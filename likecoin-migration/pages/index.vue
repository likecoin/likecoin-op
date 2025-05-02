<template>
  <div :class="['flex-1', 'min-h-0', 'bg-likecoin-lightergrey', 'pb-4']">
    <div :class="['-mb-[65px]']">
      <HeroBanner>
        <h1
          :class="[
            'text-3xl',
            'font-inter',
            'font-semibold',
            'text-likecoin-votecolor-yes',
          ]"
        >
          {{ $t('app.title') }}
        </h1>
      </HeroBanner>
    </div>
    <div :class="['relative', 'max-w-[880px]', 'px-4', 'mx-auto']">
      <div :class="['bg-white', 'p-[30px]', 'rounded-md', 'shadow-md']">
        <StepSection :step="1" :current-step="currentStep.step">
          <template #default="{ isCurrent, isPast }">
            <h2
              :class="[
                'text-base',
                {
                  ['font-semibold']: isCurrent() || isPast(),
                },
                'leading-[30px]',
                'text-likecoin-darkgrey',
              ]"
            >
              {{ $t('section.introduction.title') }}
            </h2>
          </template>
          <template #current>
            <div :class="['mt-2', 'mb-5']">
              <SectionIntroduction
                @confirmClicked="handleIntroductionSectionConfirmClick"
              />
            </div>
          </template>
        </StepSection>
        <StepSection :step="2" :current-step="currentStep.step">
          <template #default="{ isCurrent, isPast }">
            <h2
              :class="[
                'text-base',
                {
                  ['font-semibold']: isCurrent() || isPast(),
                },
                'leading-[30px]',
                'text-likecoin-darkgrey',
              ]"
            >
              {{ $t('section.connect-wallet.title') }}
            </h2>
          </template>
          <template #current>
            <SectionWalletConnect
              :class="['mt-2.5', 'mb-4']"
              :liker-id="likerId"
              :avatar="avatar"
              :cosmos-address="cosmosAddress"
              :email="email"
              :preferred-evm-provider-id="preferredEvmProviderId"
              :current-balance="currentBalance"
              :eth-address="ethAddress"
              :estimated-balance="estimatedBalance"
              @likeCoinWalletConnected="handleLikeCoinWalletConnected"
              @likeCoinEVMWalletConnected="handleLikeCoinEVMWalletConnected"
            />
          </template>
          <template #past>
            <SectionWalletConnect
              :class="['mt-2.5', 'mb-4']"
              :liker-id="likerId"
              :avatar="avatar"
              :cosmos-address="cosmosAddress"
              :email="email"
              :preferred-evm-provider-id="preferredEvmProviderId"
              :current-balance="currentBalance"
              :eth-address="ethAddress"
              :estimated-balance="estimatedBalance"
            />
          </template>
        </StepSection>
        <StepSection :step="3" :current-step="currentStep.step">
          <template #future>
            <h2
              :class="['text-base', 'leading-[30px]', 'text-likecoin-darkgrey']"
            >
              {{ $t('section.confirm-by-signing.title') }}
            </h2>
          </template>
          <template #current>
            <SectionEvmSign
              :signing-message="
                currentStep.step === 3 ? currentStep.ethSigningMessage : ''
              "
              @signed="handleEvmSigned"
            >
              <template #title>
                <h2
                  :class="[
                    'text-base',
                    'font-semibold',
                    'leading-[30px]',
                    'text-likecoin-darkgrey',
                  ]"
                >
                  {{ $t('section.confirm-by-signing.title') }}
                </h2>
              </template>
            </SectionEvmSign>
          </template>
          <template #past>
            <h2
              :class="[
                'text-base',
                'leading-[30px]',
                'font-semibold',
                'text-likecoin-darkgrey',
              ]"
            >
              {{ $t('section.confirm-by-signing.title') }}
            </h2>
          </template>
        </StepSection>
        <StepSection :step="4" :current-step="currentStep.step">
          <template #future>
            <h2
              :class="['text-base', 'leading-[30px]', 'text-likecoin-darkgrey']"
            >
              {{ $t('section.migration-progress.title') }}
            </h2>
          </template>
          <template #current>
            <SectionMigrationProgress
              v-if="migration != null"
              :enable-retry="!isMigratedThroughStep"
              :estimated-balance="estimatedBalance"
              :migration-status="migration.status"
              :cosmos-tx-hash="migration.cosmos_tx_hash"
              :evm-tx-hash="migration.evm_tx_hash"
              :error-message="migrationErrorMessage"
              @retry="handleRetry"
              @restart="handleRestart"
            >
              <template #title>
                <h2
                  :class="[
                    'text-base',
                    'font-semibold',
                    'text-likecoin-darkgrey',
                  ]"
                >
                  {{ $t('section.migration-progress.title') }}
                </h2>
              </template>
            </SectionMigrationProgress>
          </template>
          <template #past>
            <SectionMigrationProgress
              v-if="migration != null"
              :enable-retry="!isMigratedThroughStep"
              :estimated-balance="estimatedBalance"
              :migration-status="migration.status"
              :cosmos-tx-hash="migration.cosmos_tx_hash"
              :evm-tx-hash="migration.evm_tx_hash"
              :error-message="migrationErrorMessage"
              @retry="handleRetry"
              @restart="handleRestart"
            >
              <template #title>
                <h2
                  :class="[
                    'text-base',
                    'font-semibold',
                    'text-likecoin-darkgrey',
                  ]"
                >
                  {{ $t('section.migration-progress.title') }}
                </h2>
              </template>
            </SectionMigrationProgress>
          </template>
        </StepSection>
        <StepSection
          v-slot="{ isFuture }"
          :step="5"
          :current-step="currentStep.step"
        >
          <h2
            :class="[
              'text-base',
              'leading-[30px]',
              {
                'font-semibold': !isFuture(),
              },
              'text-likecoin-darkgrey',
            ]"
          >
            {{ $t('section.migration-completed.title') }}
          </h2>
        </StepSection>
      </div>
    </div>
    <DelayedFullScreenLoading :is-loading="isTransitioning" />
  </div>
</template>

<script lang="ts">
import {
  DeliverTxResponse,
  parseCoins,
  SigningStargateClient,
} from '@cosmjs/stargate';
import {
  LikeCoinWalletConnectorConnectionResult,
  LikeCoinWalletConnectorMethodType,
} from '@likecoin/wallet-connector';
import { isAxiosError } from 'axios';
import { Decimal } from 'decimal.js';
import Vue from 'vue';

import {
  CreateCosmosMemoData,
  makeCreateCosmosMemoDataAPI,
} from '~/apis/createCosmosMemoData';
import {
  CreateEthSigningMessage,
  makeCreateEthSigningMessageAPI,
} from '~/apis/createEthSigningMessage';
import {
  CreateLikeCoinMigration,
  makeCreateLikeCoinMigrationAPI,
} from '~/apis/createLikeCoinMigration';
import {
  GetLatestLikeCoinMigration,
  makeGetLatestLikeCoinMigrationAPI,
} from '~/apis/getLatestLikeCoinMigration';
import { GetUserProfile, makeGetUserProfileAPI } from '~/apis/getUserProfile';
import {
  Completed,
  Failed,
  LikeCoinMigration,
  Pending,
  Polling,
} from '~/apis/models/likeCoinMigration';
import {
  makeUpdateLikeCoinMigrationCosmosTxHash,
  UpdateLikeCoinMigrationCosmosTxHash,
} from '~/apis/updateLikeCoinMigrationCosmosTxHash';
import {
  ChainCoin,
  convertViewCoinToChainCoin,
  isChainCoin,
  isViewCoin,
} from '~/models/cosmosNetworkConfig';
import {
  authcoreRedirected,
  authcoreRedirectionFailed,
  completedMigrationResolved,
  EitherEthConnected,
  EthConnected,
  ethSignConfirming,
  evmConnected,
  failedMigrationResolved,
  gasEstimated,
  initCosmosConnected,
  introductionConfirmed,
  isEthConnected,
  likerIdResolved,
  migrationCancelledByCosmosNotSigned,
  migrationRetryCosmosSign,
  migrationRetryFailed,
  pendingMigrationResolved,
  pollingMigrationResolved,
  restart,
  StepState,
  StepStateStep2AuthcoreRedirected,
  StepStateStep2CosmosConnected,
  StepStateStep2GasEstimated,
  StepStateStep2Init,
  StepStateStep2LikerIdResolved,
  StepStateStep3AwaitSignature,
  StepStateStep4Failed,
  StepStateStep4Pending,
  StepStateStep4PendingCosmosSignCancelled,
  StepStateStep4Polling,
  StepStateStepEnd,
} from '~/pageModels';

interface Data {
  isTransitioning: boolean;
  isMigratedThroughStep: boolean;

  currentStep: StepState;

  migrationFetchTimeout: ReturnType<typeof setTimeout> | null;
}

export default Vue.extend({
  data(): Data {
    return {
      isTransitioning: false,
      isMigratedThroughStep: false,

      currentStep: { step: 1 },

      migrationFetchTimeout: null,
    };
  },

  computed: {
    likerId(): string | null {
      if ('likerId' in this.currentStep) {
        return this.currentStep.likerId;
      }
      return null;
    },
    avatar(): string | null {
      if ('avatar' in this.currentStep) {
        return this.currentStep.avatar;
      }
      return null;
    },
    cosmosAddress(): string | null {
      if ('cosmosAddress' in this.currentStep) {
        return this.currentStep.cosmosAddress;
      }
      return null;
    },
    email(): string | null {
      let email: string | null = null;
      if ('connection' in this.currentStep) {
        const { user } = this.currentStep.connection;
        email = user?.primary_email ?? null;
      }
      return email;
    },
    preferredEvmProviderId(): string | null {
      if ('connection' in this.currentStep) {
        const { method } = this.currentStep.connection;
        switch (method) {
          case 'liker-id':
          case 'likerland-app':
            return 'email';

          case 'keplr':
          case 'keplr-mobile':
            return 'app.keplr';

          case 'cosmostation':
          case 'cosmostation-mobile':
            return 'io.cosmostation';

          case 'leap':
            return 'io.leapwallet.LeapWallet';

          case 'metamask-leap':
            return 'io.metamask';

          default:
            break;
        }
      }
      return null;
    },
    currentBalance(): ChainCoin | null {
      if ('currentBalance' in this.currentStep) {
        return this.currentStep.currentBalance;
      }
      return null;
    },
    currentBalanceOverride(): ChainCoin | null {
      if (
        this.$route.query.currentBalanceOverride != null &&
        typeof this.$route.query.currentBalanceOverride === 'string'
      ) {
        const coin = parseCoins(this.$route.query.currentBalanceOverride)[0];
        if (isChainCoin(coin)) {
          return coin;
        }
        if (isViewCoin(coin)) {
          return convertViewCoinToChainCoin(coin, this.$cosmosNetworkConfig);
        }
      }
      return null;
    },
    ethAddress(): string | null {
      if ('ethAddress' in this.currentStep) {
        return this.currentStep.ethAddress;
      }
      return null;
    },
    estimatedBalance(): ChainCoin | null {
      if ('estimatedBalance' in this.currentStep) {
        return this.currentStep.estimatedBalance;
      }
      return null;
    },
    migration(): LikeCoinMigration | null {
      if ('migration' in this.currentStep) {
        return this.currentStep.migration;
      }
      return null;
    },
    migrationErrorMessage(): string | null {
      if (
        this.currentStep.step === 4 &&
        this.currentStep.state === 'PendingCosmosSignCancelled'
      ) {
        return this.currentStep.cancelReason;
      }
      return this.migration?.failed_reason ?? null;
    },

    getUserProfile(): GetUserProfile {
      if (this.cosmosAddress == null) {
        throw new Error('cosmos address not connected');
      }
      return makeGetUserProfileAPI(this.cosmosAddress)(this.$apiClient);
    },
    createEthSigningMessage(): CreateEthSigningMessage {
      return makeCreateEthSigningMessageAPI(this.$apiClient);
    },
    createCosmosMemoData(): CreateCosmosMemoData {
      return makeCreateCosmosMemoDataAPI(this.$apiClient);
    },
    createLikeCoinMigration(): CreateLikeCoinMigration {
      return makeCreateLikeCoinMigrationAPI(this.$apiClient);
    },
    getLatestLikeCoinMigration(): GetLatestLikeCoinMigration {
      return makeGetLatestLikeCoinMigrationAPI(this.$apiClient);
    },
    updateLikeCoinMigrationCosmosTxHash(): UpdateLikeCoinMigrationCosmosTxHash {
      return makeUpdateLikeCoinMigrationCosmosTxHash(this.$apiClient);
    },
  },

  watch: {
    currentStep(currentStep: StepState) {
      if (this.migrationFetchTimeout != null) {
        clearTimeout(this.migrationFetchTimeout);
        this.migrationFetchTimeout = null;
      }

      if (currentStep.step === 4 && currentStep.state === 'Polling') {
        this.migrationFetchTimeout = setTimeout(async () => {
          // Also check after timeout
          if (
            this.currentStep.step !== 4 ||
            this.currentStep.state !== 'Polling'
          ) {
            return;
          }
          this.currentStep = await this._asyncStateTransition(
            this.currentStep,
            (s) => this._refreshMigration(s)
          );
        }, 1000);
      }
    },
  },

  async mounted() {
    await this.handleMaybeLikeCoinWalletConnectedFromRedirect();
  },

  methods: {
    handleIntroductionSectionConfirmClick() {
      if (this.currentStep.step !== 1) {
        return;
      }
      this.currentStep = introductionConfirmed(this.currentStep);
    },

    async handleMaybeLikeCoinWalletConnectedFromRedirect() {
      const { code, method, ...query } = this.$route.query;
      if (method && code) {
        this.$router.replace({ query });
        await this.handleLikeCoinWalletAuthcoreRedirected(method, code);
      }
    },

    async handleLikeCoinWalletAuthcoreRedirected(
      method: string | (string | null)[],
      code: string | (string | null)[]
    ) {
      this.currentStep = authcoreRedirected(this.currentStep, method, code);
      this.currentStep = await this._asyncStateTransition(
        this.currentStep,
        (s) => this._handleAuthcoreRedirect(s)
      );

      if (this.currentStep.state === 'CosmosConnected') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._resolveLikerId
        );
      }

      if (this.currentStep.state === 'LikerIdResolved') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._estimateBalance
        );
      }

      if (this.currentStep.state === 'GasEstimated') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._getLatestMigration
        );
      }

      if (
        this.currentStep.step === 2 &&
        this.currentStep.state === 'GasEstimated' &&
        isEthConnected(this.currentStep)
      ) {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._resolveEvmSigningMessage
        );
      }

      if (this.currentStep.step === 4 && this.currentStep.state === 'Pending') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._sendCosmosToken
        );
      }
    },

    async handleLikeCoinWalletConnected(
      cosmosAddress: string,
      connection: LikeCoinWalletConnectorConnectionResult
    ) {
      this.currentStep = initCosmosConnected(
        this.currentStep,
        cosmosAddress,
        connection
      );

      this.currentStep = await this._asyncStateTransition(
        this.currentStep,
        this._resolveLikerId
      );

      if (this.currentStep.state === 'LikerIdResolved') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._estimateBalance
        );
      }

      if (this.currentStep.state === 'GasEstimated') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._getLatestMigration
        );
      }

      if (
        this.currentStep.step === 2 &&
        this.currentStep.state === 'GasEstimated' &&
        isEthConnected(this.currentStep)
      ) {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._resolveEvmSigningMessage
        );
      }

      if (this.currentStep.step === 4 && this.currentStep.state === 'Pending') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._sendCosmosToken
        );
      }
    },

    async handleLikeCoinEVMWalletConnected(ethAddress: string) {
      if (this.currentStep.step === 2) {
        this.currentStep = evmConnected(this.currentStep, ethAddress);
      }

      if (
        this.currentStep.step === 2 &&
        this.currentStep.state === 'GasEstimated' &&
        isEthConnected(this.currentStep)
      ) {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._resolveEvmSigningMessage
        );
      }
    },

    async handleEvmSigned(signature: string) {
      if (this.currentStep.step !== 3) {
        return;
      }
      this.currentStep = await this._asyncStateTransition(
        this.currentStep,
        (s) => this._createMigration(s, signature, s.ethSigningMessage)
      );

      if (this.currentStep.step === 4 && this.currentStep.state === 'Pending') {
        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          this._sendCosmosToken
        );
      }
    },

    async handleRetry() {
      if (this.currentStep.step === 4) {
        if (this.currentStep.state === 'PendingCosmosSignCancelled') {
          this.currentStep = migrationRetryCosmosSign(this.currentStep);
          this.currentStep = await this._asyncStateTransition(
            this.currentStep,
            this._sendCosmosToken
          );
        } else if (this.currentStep.state === 'Failed') {
          this.currentStep = migrationRetryFailed(this.currentStep);
        }
      }
    },

    async handleRestart() {
      if (this.currentStep.step === 99999) {
        this.currentStep = restart(this.currentStep);

        this.currentStep = await this._asyncStateTransition(
          this.currentStep,
          (s) => this._estimateBalance(s)
        );
      }
    },

    async _sendCosmosToken(
      s: StepStateStep4Pending
    ): Promise<
      | StepStateStep4Pending
      | StepStateStep4PendingCosmosSignCancelled
      | StepStateStep4Polling
      | StepStateStep4Failed
      | StepStateStepEnd
    > {
      const cosmosMemoData = await this.createCosmosMemoData({
        signature: s.evmSignature,
        amount: s.estimatedBalance,
        ethAddress: s.ethAddress,
      });

      let tx: DeliverTxResponse;
      try {
        tx = await s.signingStargateClient.sendTokens(
          s.cosmosAddress,
          this.$appConfig.cosmosDepositAddress,
          [s.estimatedBalance],
          {
            amount: [s.estimatedBalance],
            gas: `${s.gasEstimation}`,
          },
          cosmosMemoData.memo_data
        );
      } catch (e) {
        if (e instanceof Error) {
          return migrationCancelledByCosmosNotSigned(s, e.message);
        }
        throw e;
      }

      const migration = await this.updateLikeCoinMigrationCosmosTxHash(
        s.cosmosAddress,
        {
          cosmos_tx_hash: tx.transactionHash,
        }
      );
      return this._resolveMigration(s, migration.migration);
    },

    async _asyncStateTransition<S1 extends StepState, S2 extends StepState>(
      state1: S1,
      asyncT: (state1: S1) => Promise<S2>
    ): Promise<S1 | S2> {
      if (this.isTransitioning) {
        return state1;
      }
      this.isTransitioning = true;
      try {
        return await asyncT(state1);
      } finally {
        this.isTransitioning = false;
      }
    },

    async _handleAuthcoreRedirect(
      currentStep: EitherEthConnected<StepStateStep2AuthcoreRedirected>
    ): Promise<
      | EitherEthConnected<StepStateStep2CosmosConnected>
      | EitherEthConnected<StepStateStep2Init>
    > {
      const { method, code } = currentStep;
      const connection = await this.$likeCoinWalletConnector.handleRedirect(
        method as LikeCoinWalletConnectorMethodType,
        { code }
      );
      if (connection != null) {
        if ('method' in connection) {
          const {
            accounts: [account],
          } = connection;
          return initCosmosConnected(currentStep, account.address, connection);
        }
      }
      return authcoreRedirectionFailed(currentStep);
    },

    async _resolveLikerId(
      s: EitherEthConnected<StepStateStep2CosmosConnected>
    ): Promise<EitherEthConnected<StepStateStep2LikerIdResolved>> {
      const userProfile = await this.getUserProfile();
      return likerIdResolved(
        s,
        userProfile.user_profile.avatar,
        userProfile.user_profile.liker_id
      );
    },

    async _estimateBalance(
      s: EitherEthConnected<StepStateStep2LikerIdResolved>
    ): Promise<EitherEthConnected<StepStateStep2GasEstimated>> {
      const { offlineSigner } = s.connection;
      const client = await SigningStargateClient.connectWithSigner(
        this.$likeCoinWalletConnector.options.rpcURL,
        offlineSigner
      );

      const balance =
        this.currentBalanceOverride != null
          ? this.currentBalanceOverride
          : ((await client.getBalance(
              s.cosmosAddress,
              this.$cosmosNetworkConfig.coinLookup[0].chainDenom
            )) as unknown as ChainCoin);

      const cosmosMemoData = await this.createCosmosMemoData({
        ethAddress: isEthConnected(s) ? s.ethAddress : `0x${'0'.repeat(40)}`,
        amount: balance,
        signature: `0x${'0'.repeat(130)}`,
      });

      const gasEstimation = await this._estimateGas(
        client,
        s.cosmosAddress,
        balance,
        cosmosMemoData.memo_data
      );

      // this is the tier selection from keplr sign dialog
      const tierMultiplier = {
        low: 1000,
        average: 10000,
        high: 1000000,
      };

      const gasFee = Decimal(
        gasEstimation *
          // Assume worst case user select high without insufficient fund
          // in the signing ui
          tierMultiplier.average
      );

      const estimatedBalance: ChainCoin = {
        denom: balance.denom,
        amount: Decimal.max(
          Decimal(balance.amount).minus(gasFee),
          Decimal(0)
        ).toString(),
      };

      return gasEstimated(s, client, balance, gasEstimation, estimatedBalance);
    },

    async _getLatestMigration(
      prev: EitherEthConnected<StepStateStep2GasEstimated>
    ): Promise<
      | EitherEthConnected<StepStateStep2GasEstimated>
      | StepStateStep4Pending
      | StepStateStep4Polling
      | StepStateStep4Failed
      | StepStateStepEnd
    > {
      try {
        const { migration } = await this.getLatestLikeCoinMigration(
          prev.cosmosAddress
        );
        return this._resolveMigration(prev, migration);
      } catch (e) {
        if (isAxiosError(e)) {
          if (e.status === 404) {
            return prev;
          }
        }
        throw e;
      }
    },

    async _resolveEvmSigningMessage(
      prev: EthConnected<StepStateStep2GasEstimated>
    ): Promise<StepStateStep3AwaitSignature> {
      const m = await this.createEthSigningMessage({
        amount: prev.estimatedBalance,
      });
      return ethSignConfirming(prev, m.signing_message);
    },

    async _createMigration(
      prev: StepStateStep3AwaitSignature,
      signature: string,
      signatureMessage: string
    ): Promise<
      | StepStateStep4Pending
      | StepStateStep4Polling
      | StepStateStep4Failed
      | StepStateStepEnd
    > {
      const { migration } = await this.createLikeCoinMigration({
        amount: prev.estimatedBalance,
        cosmos_address: prev.cosmosAddress,
        eth_address: prev.ethAddress,
        evm_signature: signature,
        evm_signature_message: signatureMessage,
      });
      this.isMigratedThroughStep = true;
      return this._resolveMigration(prev, migration);
    },

    async _refreshMigration(
      s: StepStateStep4Polling
    ): Promise<
      | StepStateStep4Pending
      | StepStateStep4Polling
      | StepStateStepEnd
      | StepStateStep4Failed
    > {
      const { migration } = await this.getLatestLikeCoinMigration(
        s.cosmosAddress
      );
      // expect throw on error
      return this._resolveMigration(s, migration);
    },

    async _updateLikeCoinMigrationCosmosTxHash(
      prev: StepStateStep4Pending,
      txHash: string
    ) {
      const migration = await this.updateLikeCoinMigrationCosmosTxHash(
        prev.cosmosAddress,
        {
          cosmos_tx_hash: txHash,
        }
      );

      return this._resolveMigration(prev, migration.migration);
    },

    async _estimateGas(
      signingStargateClient: SigningStargateClient,
      cosmosAddress: string,
      amount: ChainCoin,
      memo: string
    ) {
      const msg = {
        typeUrl: '/cosmos.bank.v1beta1.MsgSend',
        value: {
          fromAddress: cosmosAddress,
          toAddress: this.$appConfig.cosmosDepositAddress,
          amount: [amount],
        },
      };
      const fluctuationMultiplier = 1.3;
      const gasEstimation = await signingStargateClient.simulate(
        cosmosAddress,
        [msg],
        memo
      );
      return Math.floor(gasEstimation * fluctuationMultiplier);
    },

    _resolveMigration(
      prev:
        | EitherEthConnected<StepStateStep2GasEstimated>
        | StepStateStep3AwaitSignature
        | StepStateStep4Pending
        | StepStateStep4Polling,
      migration: LikeCoinMigration
    ):
      | StepStateStep4Pending
      | StepStateStep4Polling
      | StepStateStep4Failed
      | StepStateStepEnd {
      switch (migration.status) {
        case 'completed': {
          return completedMigrationResolved(
            prev,
            migration as Completed<LikeCoinMigration>
          );
        }
        case 'failed': {
          return failedMigrationResolved(
            prev,
            migration as Failed<LikeCoinMigration>
          );
        }
        case 'evm_minting':
        case 'evm_verifying':
        case 'verifying_cosmos_tx':
          return pollingMigrationResolved(
            prev,
            migration as Polling<LikeCoinMigration>
          );
        case 'pending_cosmos_tx_hash':
          return pendingMigrationResolved(
            prev,
            migration as Pending<LikeCoinMigration>
          );
      }
    },
  },
});
</script>
