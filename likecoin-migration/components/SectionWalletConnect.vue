<template>
  <div
    :class="[
      'rounded-md',
      'bg-likecoin-lightergrey',
      'p-4',
      'flex',
      'flex-col',
      'sm:flex-row',
      'gap-4',
      'sm:gap-6',
    ]"
  >
    <div :class="['flex-1']">
      <div v-if="cosmosAddress == null">
        <h3
          :class="['text-base', 'font-semibold', 'text-likecoin-black', 'mb-3']"
        >
          {{ $t('section.wallet-connect.likecoin-wallet') }}
        </h3>
        <AppButton
          :class="['w-full']"
          @click="handleConnectLikeCoinWalletClick"
        >
          {{ $t('section.wallet-connect.connect-likecoin-wallet') }}
        </AppButton>
      </div>
      <div v-else>
        <h3
          :class="[
            'text-base',
            'font-semibold',
            'text-likecoin-darkgreen',
            'mb-4',
          ]"
        >
          {{ $t('section.wallet-connect.likecoin-wallet') }}
        </h3>
        <div :class="['grid', 'grid-cols-2', 'gap-x-2.5', 'gap-y-1']">
          <div :class="['text-base', 'text-likecoin-darkgrey']">
            {{ $t('section.wallet-connect.liker-id') }}
          </div>
          <div :class="['flex', 'flex-row', 'items-center', 'gap-1']">
            <img
              v-if="avatar != null"
              :src="avatar"
              :class="['w-[18px]', 'h-[18px]', 'rounded-full']"
            />
            <div
              :class="[
                'text-base',
                'text-likecoin-votecolor-yes',
                'overflow-hidden',
                'text-ellipsis',
              ]"
            >
              {{ likerId }}
            </div>
          </div>
          <div :class="['text-base', 'text-likecoin-darkgrey']">
            {{ $t('section.wallet-connect.cosmos-address') }}
          </div>
          <div :class="['flex', 'flex-row', 'items-center', 'gap-1']">
            <div
              :class="[
                'text-base',
                'text-likecoin-votecolor-yes',
                'overflow-hidden',
                'text-ellipsis',
              ]"
            >
              {{ cosmosAddress }}
            </div>
            <button
              v-if="cosmosAddress != null"
              type="button"
              @click="handleCosmosAddressCopyClick(cosmosAddress)"
            >
              <FontAwesomeIcon :class="['text-base']" icon="copy" />
            </button>
          </div>
          <div :class="['text-base', 'text-likecoin-darkgrey']">
            {{ $t('section.wallet-connect.current-balance') }}
          </div>
          <div
            :class="[
              'text-base',
              'text-likecoin-votecolor-yes',
              'overflow-hidden',
              'text-ellipsis',
            ]"
          >
            <CoinDisplay v-if="currentBalance != null" :coin="currentBalance" />
          </div>
        </div>
      </div>
    </div>
    <div
      :class="[
        'max-sm:hidden',
        'border-l',
        'border-l-likecoin-grey',
        'self-stretch',
      ]"
    ></div>
    <div :class="['flex-1']">
      <div v-if="ethAddress == null">
        <h3
          :class="['text-base', 'font-semibold', 'text-likecoin-black', 'mb-3']"
        >
          {{ $t('section.wallet-connect.migration-wallet') }}
        </h3>
        <AppButton
          variant="secondary"
          :class="['w-full']"
          :disabled="isV3ButtonDisabled"
          @click="handleConnectTargetWalletClick"
        >
          {{ $t('section.wallet-connect.connect-migration-wallet') }}
        </AppButton>
      </div>
      <div v-else>
        <h3
          :class="[
            'text-base',
            'font-semibold',
            'text-likecoin-darkgreen',
            'mb-4',
          ]"
        >
          {{ $t('section.wallet-connect.migration-wallet') }}
        </h3>
        <div :class="['grid', 'grid-cols-2', 'gap-x-2.5', 'gap-y-1']">
          <div :class="['text-base', 'text-likecoin-darkgrey']">
            {{ $t('section.wallet-connect.eth-address') }}
          </div>
          <div :class="['flex', 'flex-row', 'items-center', 'gap-1']">
            <div
              :class="[
                'text-base',
                'text-likecoin-votecolor-yes',
                'overflow-hidden',
                'text-ellipsis',
              ]"
            >
              {{ ethAddress }}
            </div>
            <button
              v-if="ethAddress != null"
              type="button"
              @click="handleEthAddressCopyClick(ethAddress)"
            >
              <FontAwesomeIcon :class="['text-base']" icon="copy" />
            </button>
          </div>
          <div :class="['text-base', 'text-likecoin-darkgrey']">
            {{ $t('section.wallet-connect.estimated-balance') }}
          </div>
          <div
            :class="[
              'text-base',
              'text-likecoin-votecolor-yes',
              'overflow-hidden',
              'text-ellipsis',
            ]"
          >
            <CoinDisplay
              v-if="estimatedBalance != null"
              :coin="estimatedBalance"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {
  LikeCoinWalletConnectorConnectionResult,
  LikeCoinWalletConnectorMethodType,
} from '@likecoin/wallet-connector';
import Vue, { PropType } from 'vue';

import { ViewCoin } from '~/models/cosmosNetworkConfig';

export default Vue.extend({
  props: {
    avatar: {
      type: String as PropType<string | null>,
      default: null,
    },
    likerId: {
      type: String as PropType<string | null>,
      default: null,
    },
    cosmosAddress: {
      type: String as PropType<string | null>,
      default: null,
    },
    email: {
      type: String as PropType<string | null>,
      default: null,
    },
    preferredEvmProviderId: {
      type: String as PropType<string | null>,
      default: null,
    },
    currentBalance: {
      type: Object as PropType<ViewCoin | null>,
      default: null,
    },
    ethAddress: {
      type: String as PropType<string | null>,
      default: null,
    },
    estimatedBalance: {
      type: Object as PropType<ViewCoin | null>,
      default: null,
    },
    isV2WalletDataFetching: {
      type: Boolean as PropType<boolean>,
      default: false,
    },
  },

  computed: {
    isV3ButtonDisabled(): boolean {
      return this.isV2WalletDataFetching;
    },
  },

  mounted() {
    this.$likeCoinWalletConnector.on(
      'account_change',
      this.handleLikeCoinWalletAccountChange
    );
    this.$likeCoinEVMWalletConnector.eventEmitter.on(
      'connect',
      this.handleLikeCoinEVMWalletConnect
    );
  },

  destroyed() {
    this.$likeCoinWalletConnector.off(
      'account_change',
      this.handleLikeCoinWalletAccountChange
    );

    this.$likeCoinEVMWalletConnector.eventEmitter.off(
      'connect',
      this.handleLikeCoinEVMWalletConnect
    );
  },

  methods: {
    async handleConnectLikeCoinWalletClick() {
      const connection =
        await this.$likeCoinWalletConnector.openConnectionMethodSelectionDialog(
          {}
        );
      this.handleLikeCoinWalletConnection(connection);
    },

    async handleCosmosAddressCopyClick(address: string) {
      await window.navigator.clipboard.writeText(address);
      alert('Cosmos address is copied to clipboard!');
    },

    async handleEthAddressCopyClick(address: string) {
      await window.navigator.clipboard.writeText(address);
      alert('Ethereum address is copied to clipboard!');
    },

    handleConnectTargetWalletClick() {
      this.$likeCoinEVMWalletConnector.connector.showConnectPortal({
        email: this.email || undefined,
      });
    },

    handleLikeCoinWalletConnection(
      connection: LikeCoinWalletConnectorConnectionResult | undefined
    ) {
      if (!connection) return;

      const {
        accounts: [account],
      } = connection;

      this.$emit('likeCoinWalletConnected', account.address, connection);
    },

    async handleLikeCoinWalletAccountChange(
      method: LikeCoinWalletConnectorMethodType
    ) {
      const connection = await this.$likeCoinWalletConnector.init(method);
      this.handleLikeCoinWalletConnection(connection);
    },

    handleLikeCoinEVMWalletConnect({
      walletAddress,
    }: {
      walletAddress: string;
    }) {
      this.$emit('likeCoinEVMWalletConnected', walletAddress);
    },
  },
});
</script>
