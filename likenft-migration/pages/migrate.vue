<template>
  <div
    class="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]"
  >
    <main class="flex flex-col gap-8 row-start-2 items-center sm:items-start">
      <primary-button @click="handleConnectCosmosWalletClick">
        {{ $t('migrate.connect-cosmos-wallet') }}
      </primary-button>
      <p>{{ $t('migrate.cosmos-wallet-address', { cosmosWalletAddress }) }}</p>
      <p>{{ $t('migrate.liker-id', { likerID }) }}</p>
      <primary-button @click="handleConnectEVMWalletClick">
        {{ $t('migrate.connect-evm-wallet') }}
      </primary-button>
      <p>{{ $t('migrate.evm-wallet-address', { evmWalletAddress }) }}</p>
      <primary-button
        v-if="cosmosWalletAddress != null && evmWalletAddress != null"
        @click="handleMigrateClick"
      >
        {{ $t('migrate.migrate') }}
      </primary-button>
    </main>
    <div
      v-if="isLoading"
      class="fixed top-0 left-0 w-full h-full bg-white/90 flex items-center justify-center"
    >
      Loading
    </div>
  </div>
</template>

<script lang="ts">
import {
  LikeCoinWalletConnector,
  LikeCoinWalletConnectorMethodType,
  LikeCoinWalletConnectorSession,
} from '@likecoin/wallet-connector';
import { Eip1193Provider } from 'ethers';
import Vue from 'vue';
import Web3 from 'web3';
import { z } from 'zod';

import { getSignMessage } from '~/apis/getSignMessage';
import { getUser } from '~/apis/getUser';
import { makeMigrateLikerIDAPI } from '~/apis/migrateLikerID';
import { LIKECOIN_WALLET_CONNECTOR_CONFIG } from '~/constant/network';

async function getEthereumAccount(
  ethereum: Eip1193Provider
): Promise<string | null> {
  const web3 = new Web3(ethereum);
  await ethereum.request({ method: 'eth_requestAccounts' });
  const accounts = await web3.eth.getAccounts();
  if (accounts.length > 0) {
    return accounts[0];
  }
  return null;
}

async function signEthereumMessage(
  message: string,
  ethereum: Eip1193Provider,
  ethereumAddress: string
) {
  const web3 = new Web3(ethereum);
  const sign = await web3.eth.personal.sign(
    message,
    ethereumAddress,
    'Password!'
  );

  return sign;
}

interface Data {
  cosmosWalletAddress: string | null;
  likerID: string | null;
  evmWalletAddress: string | null;
  isLoading: boolean;
}

export default Vue.extend({
  data(): Data {
    return {
      cosmosWalletAddress: null,
      likerID: null,
      evmWalletAddress: null,
      isLoading: false,
    };
  },

  computed: {
    connector() {
      return new LikeCoinWalletConnector(
        LIKECOIN_WALLET_CONNECTOR_CONFIG(this.$appConfig.isTestnet)
      );
    },
    getSignMessage() {
      return getSignMessage(this.$apiClient);
    },
    migrateLikerID() {
      return makeMigrateLikerIDAPI(this.$apiClient);
    },
  },

  methods: {
    async handleConnectCosmosWalletClick() {
      const connection =
        await this.connector?.openConnectionMethodSelectionDialog({});
      this.handleConnection(connection);
    },

    async handleConnection(
      connection: LikeCoinWalletConnectorSession | undefined
    ) {
      if (!connection) return;
      const {
        accounts: [account],
      } = connection;
      this.cosmosWalletAddress = account.address;

      this.isLoading = true;
      try {
        const user = await getUser(this.cosmosWalletAddress);
        this.likerID = user.liker_id;
      } finally {
        this.isLoading = false;
      }

      this.connector?.once('account_change', this.handleAccountChange);
    },

    async handleAccountChange(method: LikeCoinWalletConnectorMethodType) {
      const connection = await this.connector?.init(method);
      this.handleConnection(connection);
    },

    async handleConnectEVMWalletClick() {
      this.isLoading = true;
      if (window.ethereum == null) {
        alert('Please install metamask extension');
        return;
      }
      try {
        this.evmWalletAddress = await getEthereumAccount(window.ethereum);
      } catch (e) {
        console.error(e);
      } finally {
        this.isLoading = false;
      }
    },

    async handleMigrateClick() {
      const S = z.object({
        cosmosWalletAddress: z.string(),
        evmWalletAddress: z.string(),
      });
      const s = S.safeParse({
        cosmosWalletAddress: this.cosmosWalletAddress,
        evmWalletAddress: this.evmWalletAddress,
      });
      if (s.data == null) {
        return;
      }
      if (window.ethereum == null) {
        alert('Please install metamask extension');
        return;
      }
      const u = await getUser(s.data.cosmosWalletAddress);
      if (u.evm_address == null) {
        const signMessage = await this.getSignMessage({
          cosmos_address: s.data.cosmosWalletAddress,
          eth_address: s.data.evmWalletAddress,
          liker_id: this.likerID,
        });
        const connection = await this.connector.initIfNecessary();
        if (connection == null) {
          alert('cannot get wallet connector connection');
          return;
        }
        const {
          accounts: [account],
          offlineSigner,
        } = connection;

        if (!offlineSigner.signArbitrary) {
          alert('signArbitrary not supported');
          return;
        }
        const result = await offlineSigner.signArbitrary(
          this.connector.options.chainId,
          account.address,
          signMessage.message
        );
        const cosmosSignature = result.signature;
        const signedMessage = await signEthereumMessage(
          signMessage.message,
          window.ethereum,
          s.data.evmWalletAddress
        );

        await this.migrateLikerID({
          cosmos_pub_key: result.pub_key.value,
          cosmos_signature: cosmosSignature,
          eth_address: s.data.evmWalletAddress,
          eth_signature: signedMessage,
          like_id: this.likerID,
          signing_message: signMessage.message,
        });
      }

      this.$router.push(`/migration-preview/${s.data.cosmosWalletAddress}`);
    },
  },
});
</script>
