<template>
  <div
    :class="[
      'flex',
      'flex-col',
      'flex-1',
      'min-h-0',
      'bg-likecoin-lightergrey',
    ]"
  >
    <div :class="['-mb-[120px]']">
      <HeroBanner>
        <h1
          :class="[
            '-mt-[75px]',
            'text-3xl',
            'font-inter',
            'font-semibold',
            'text-likecoin-votecolor-yes',
          ]"
        >
          {{ $t("app.likecoin.transaction_details") }}
        </h1>
      </HeroBanner>
    </div>
    <div
      :class="[
        'mx-16',
        'mb-16',
        'px-4',
        'flex',
        'flex-1',
        'flex-col',
        'gap-6',
        'z-10',
      ]"
    >
      <!-- Loading state -->
      <div
        v-if="loading"
        :class="[
          'flex-1',
          'text-center',
          'flex',
          'justify-center',
          'items-center',
          'flex-col',
        ]"
      >
        <LoadingIcon />
        <p>{{ $t("common.loading") }}</p>
      </div>

      <!-- Error state -->
      <div
        v-else-if="error"
        :class="[
          'flex-1',
          'text-center',
          'flex',
          'justify-center',
          'items-center',
          'flex-col',
        ]"
      >
        <p :class="['text-red-500']">{{ error }}</p>
      </div>

      <!-- Content -->
      <template v-else-if="migration">
        <!-- Transaction Header -->
        <UCard :class="['w-full']">
          <div :class="['flex', 'flex-col', 'gap-4']">
            <!-- Transaction Status -->
            <div :class="['flex', 'items-center', 'justify-between']">
              <div :class="['flex', 'items-center', 'gap-4']">
                <div :class="['flex', 'items-center', 'gap-2']">
                  <div
                    :class="[
                      'w-8',
                      'h-8',
                      'rounded-full',
                      'flex',
                      'items-center',
                      'justify-center',
                      {
                        'bg-red-100': migration.status === 'failed',
                        'bg-green-100': migration.status === 'completed',
                        'bg-yellow-100':
                          migration.status !== 'failed' &&
                          migration.status !== 'completed',
                      },
                    ]"
                  >
                    <span
                      :class="[
                        {
                          'text-red-500': migration.status === 'failed',
                          'text-green-500': migration.status === 'completed',
                          'text-yellow-500':
                            migration.status !== 'failed' &&
                            migration.status !== 'completed',
                        },
                      ]"
                    >
                      <template v-if="migration.status === 'failed'"
                        >✕</template
                      >
                      <template v-else-if="migration.status === 'completed'"
                        >✓</template
                      >
                      <template v-else>⟳</template>
                    </span>
                  </div>
                  <div>
                    <h2 :class="['text-xl', 'font-semibold']">
                      {{ $t("section.likecoin-migration.title") }}
                    </h2>
                    <p :class="['text-sm', 'text-gray-500']">
                      {{ formatDate(migration.created_at) }}
                    </p>
                  </div>
                </div>
              </div>
              <div
                :class="[
                  'px-4',
                  'py-2',
                  'rounded-full',
                  'text-sm',
                  'font-medium',
                  {
                    'bg-red-100 text-red-800': migration.status === 'failed',
                    'bg-green-100 text-green-800':
                      migration.status === 'completed',
                    'bg-yellow-100 text-yellow-800':
                      migration.status === 'pending_cosmos_tx_hash',
                    'bg-blue-100 text-blue-800':
                      migration.status === 'verifying_cosmos_tx',
                    'bg-purple-100 text-purple-800':
                      migration.status === 'evm_minting',
                    'bg-indigo-100 text-indigo-800':
                      migration.status === 'evm_verifying',
                  },
                ]"
              >
                {{ $t(getStatusKey(migration.status)) }}
              </div>
            </div>

            <!-- Transaction Amount -->
            <div :class="['flex', 'justify-between', 'items-center', 'mt-2']">
              <span :class="['text-gray-600']">{{
                $t("section.likecoin-migration.table.header.amount")
              }}</span>
              <span :class="['text-2xl', 'font-bold']"
                >{{ amount }} {{ currency }}</span
              >
            </div>
          </div>
        </UCard>

        <!-- Transaction Information -->
        <UCard :class="['w-full']">
          <h3 :class="['text-lg', 'font-semibold', 'mb-4']">
            {{ $t("migration.information") }}
          </h3>
          <div :class="['grid', 'grid-cols-1', 'md:grid-cols-2', 'gap-4']">
            <!-- Transaction ID -->
            <div :class="['flex', 'flex-col']">
              <span :class="['text-sm', 'text-gray-500']">{{
                $t("section.likecoin-migration.table.header.id")
              }}</span>
              <div :class="['flex', 'items-center', 'gap-2']">
                <span :class="['font-medium']">{{ migration.id }}</span>
              </div>
            </div>

            <!-- User Cosmos Address -->
            <div :class="['flex', 'flex-col']">
              <span :class="['text-sm', 'text-gray-500']">{{
                $t("section.likecoin-migration.table.header.cosmos-address")
              }}</span>
              <div :class="['flex', 'items-center', 'gap-2']">
                <span :class="['font-medium', 'truncate']">{{
                  migration.user_cosmos_address
                }}</span>
                <button
                  :class="[
                    'p-1.5',
                    'rounded-md',
                    'bg-gray-100',
                    'hover:bg-gray-200',
                    'text-gray-600',
                    'hover:text-gray-800',
                    'transition-colors',
                  ]"
                  @click="copyToClipboard(migration.user_cosmos_address)"
                >
                  <FontAwesomeIcon icon="copy" />
                </button>
              </div>
            </div>

            <!-- User ETH Address -->
            <div :class="['flex', 'flex-col']">
              <span :class="['text-sm', 'text-gray-500']">{{
                $t("section.likecoin-migration.table.header.eth-address")
              }}</span>
              <div :class="['flex', 'items-center', 'gap-2']">
                <span :class="['font-medium', 'truncate']">{{
                  migration.user_eth_address
                }}</span>
                <button
                  :class="[
                    'p-1.5',
                    'rounded-md',
                    'bg-gray-100',
                    'hover:bg-gray-200',
                    'text-gray-600',
                    'hover:text-gray-800',
                    'transition-colors',
                  ]"
                  @click="copyToClipboard(migration.user_eth_address)"
                >
                  <FontAwesomeIcon icon="copy" />
                </button>
              </div>
            </div>

            <!-- Network -->
            <div :class="['flex', 'flex-col']">
              <span :class="['text-sm', 'text-gray-500']">{{
                $t("migration.network")
              }}</span>
              <span :class="['font-medium']">{{ chain }}</span>
            </div>
          </div>
        </UCard>

        <!-- Transaction Hashes -->
        <UCard :class="['w-full']">
          <h3 :class="['text-lg', 'font-semibold', 'mb-4']">
            {{ $t("migration.hashes") }}
          </h3>
          <div :class="['flex', 'flex-col', 'gap-4']">
            <!-- Cosmos TX Hash -->
            <div :class="['flex', 'flex-col']">
              <span :class="['text-sm', 'text-gray-500']">{{
                $t("section.likecoin-migration.table.header.cosmos-tx-hash")
              }}</span>
              <div
                v-if="migration.cosmos_tx_hash"
                :class="['flex', 'items-center', 'gap-2']"
              >
                <span :class="['font-medium', 'truncate']">{{
                  migration.cosmos_tx_hash
                }}</span>
                <button
                  :class="[
                    'p-1.5',
                    'rounded-md',
                    'bg-gray-100',
                    'hover:bg-gray-200',
                    'text-gray-600',
                    'hover:text-gray-800',
                    'transition-colors',
                  ]"
                  @click="copyToClipboard(migration.cosmos_tx_hash)"
                >
                  <FontAwesomeIcon icon="copy" />
                </button>
              </div>
              <span v-else :class="['text-gray-400']">-</span>
            </div>

            <!-- EVM TX Hash -->
            <div :class="['flex', 'flex-col']">
              <span :class="['text-sm', 'text-gray-500']">{{
                $t("section.likecoin-migration.table.header.evm-tx-hash")
              }}</span>
              <div
                v-if="migration.evm_tx_hash"
                :class="['flex', 'items-center', 'gap-2']"
              >
                <span :class="['font-medium', 'truncate']">{{
                  migration.evm_tx_hash
                }}</span>
                <button
                  :class="[
                    'p-1.5',
                    'rounded-md',
                    'bg-gray-100',
                    'hover:bg-gray-200',
                    'text-gray-600',
                    'hover:text-gray-800',
                    'transition-colors',
                  ]"
                  @click="copyToClipboard(migration.evm_tx_hash)"
                >
                  <FontAwesomeIcon icon="copy" />
                </button>
              </div>
              <span v-else :class="['text-gray-400']">-</span>
            </div>
          </div>
        </UCard>

        <!-- Error Details (if failed) -->
        <UCard
          v-if="migration.status === 'failed' && migration.failed_reason"
          :class="['w-full', 'bg-red-50']"
        >
          <h3 :class="['text-lg', 'font-semibold', 'mb-4', 'text-red-700']">
            {{ $t("migration.error_details") }}
          </h3>
          <span :class="['font-medium']">{{ migration.failed_reason }}</span>
        </UCard>

        <!-- Action Buttons -->
        <div :class="['flex', 'justify-end', 'gap-4', 'mt-4']">
          <AppButton
            v-if="pendingMigration"
            variant="warning"
            :loading="deleting"
            @click="promptFailLikeCoinMigration(pendingMigration)"
          >
            {{ $t("migration.fail_likecoin_migration") }}
          </AppButton>
          <AppButton
            v-if="failedMigration"
            variant="warning"
            :loading="deleting"
            @click="promptRetryLikeCoinMigration(failedMigration)"
          >
            {{ $t("migration.retry_likecoin_migration") }}
          </AppButton>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import { format } from "date-fns";
import numeral from "numeral";
import Vue from "vue";

import { makeFailLikeCoinMigrationAPI } from "~/apis/FailLikeCoinMigration";
import { makeGetLikeCoinMigrationsAPI } from "~/apis/GetLikeCoinMigration";
import {
  LikeCoinMigration,
  LikeCoinMigrationStatus,
} from "~/apis/models/likecoinMigration";
import { makeRetryLikeCoinMigrationAPI } from "~/apis/RetryLikeCoinMigration";
import AppButton from "~/components/AppButton.vue";
import HeroBanner from "~/components/HeroBanner.vue";
import { LIKECOIN_CHAIN_DENOM, LIKECOIN_CHAIN_NAME } from "~/constant";
import UCard from "~/nuxtui/components/UCard.vue";

interface Data {
  migration: LikeCoinMigration | null;
  loading: boolean;
  deleting: boolean;
  error: string | null;
}

export default Vue.extend({
  components: {
    HeroBanner,
    AppButton,
    UCard,
  },

  data(): Data {
    return {
      migration: null,
      loading: true,
      deleting: false,
      error: null,
    };
  },

  computed: {
    migrationId(): number {
      return parseInt(this.$route.params.id, 10);
    },
    amount(): string {
      const am = parseInt(this.migration?.amount || "0") * Math.pow(10, -9);
      return numeral(am).format("0,0.[0000]");
    },
    currency(): string {
      return LIKECOIN_CHAIN_DENOM(this.$appConfig.isTestnet);
    },
    chain(): string {
      return LIKECOIN_CHAIN_NAME(this.$appConfig.isTestnet);
    },
    failedMigration(): LikeCoinMigration | null {
      if (this.migration != null && this.migration.status === "failed") {
        return this.migration;
      }
      return null;
    },
    pendingMigration(): LikeCoinMigration | null {
      if (
        this.migration != null &&
        this.migration.status === "pending_cosmos_tx_hash"
      ) {
        return this.migration;
      }
      return null;
    },
  },

  mounted() {
    this.fetchMigration();
  },

  methods: {
    async fetchMigration() {
      this.loading = true;
      this.error = null;

      try {
        const resp = await makeGetLikeCoinMigrationsAPI(this.migrationId)(
          this.$apiClient
        )();
        this.migration = resp.migration;
      } catch (error) {
        this.error = this.$t("migration.error_fetching") as string;
      } finally {
        this.loading = false;
      }
    },

    formatDate(date: Date): string {
      return format(date, "yyyy-MM-dd HH:mm:ss");
    },

    getStatusKey(status: LikeCoinMigrationStatus): string {
      switch (status) {
        case LikeCoinMigrationStatus.PendingCosmosTxHash:
          return "section.likecoin-migration.table.data.status.pending_cosmos_tx_hash";
        case LikeCoinMigrationStatus.VerifyingCosmosTx:
          return "section.likecoin-migration.table.data.status.verifying_cosmos_tx";
        case LikeCoinMigrationStatus.EvmMinting:
          return "section.likecoin-migration.table.data.status.evm_minting";
        case LikeCoinMigrationStatus.EvmVerifying:
          return "section.likecoin-migration.table.data.status.evm_verifying";
        case LikeCoinMigrationStatus.Completed:
          return "section.likecoin-migration.table.data.status.completed";
        case LikeCoinMigrationStatus.Failed:
          return "section.likecoin-migration.table.data.status.failed";
        default:
          throw new Error("Invalid status");
      }
    },

    async retryLikeCoinMigration(failedMigration: LikeCoinMigration) {
      this.deleting = true;
      try {
        const resp = await makeRetryLikeCoinMigrationAPI(failedMigration.id)(
          this.$apiClient
        )();
        this.migration = resp.migration;
      } finally {
        this.deleting = false;
      }
    },

    async promptRetryLikeCoinMigration(failedMigration: LikeCoinMigration) {
      if (
        confirm(this.$t("migration.confirm_retry_likecoin_migration") as string)
      ) {
        await this.retryLikeCoinMigration(failedMigration);
      }
    },

    async failLikeCoinMigration(pendingMigration: LikeCoinMigration) {
      this.deleting = true;
      try {
        const resp = await makeFailLikeCoinMigrationAPI(pendingMigration.id)(
          this.$apiClient
        )();
        this.migration = resp.migration;
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : String(error);
        alert(`${this.$t("migration.error_failing")}: ${errorMessage}`);
      } finally {
        this.deleting = false;
      }
    },

    async promptFailLikeCoinMigration(pendingMigration: LikeCoinMigration) {
      if (
        confirm(this.$t("migration.confirm_fail_likecoin_migration") as string)
      ) {
        await this.failLikeCoinMigration(pendingMigration);
      }
    },

    async copyToClipboard(text: string) {
      await navigator.clipboard.writeText(text);
      alert(this.$t("common.copied_to_clipboard"));
    },
  },
});
</script>
