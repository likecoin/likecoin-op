import { LikeCoinWalletConnector } from '@likecoin/wallet-connector';
import { Inject, NuxtApp } from '@nuxt/types/app';

import { LIKECOIN_WALLET_CONNECTOR_CONFIG } from '~/constant/network';

export default function (app: NuxtApp, inject: Inject) {
  const likeCoinWalletConnector = new LikeCoinWalletConnector(
    LIKECOIN_WALLET_CONNECTOR_CONFIG(
      app.$appConfig.isTestnet,
      app.$appConfig.authcoreRedirectUrl
    )
  );
  inject('likeCoinWalletConnector', likeCoinWalletConnector);
}

declare module 'vue/types/vue' {
  interface Vue {
    $likeCoinWalletConnector: LikeCoinWalletConnector;
  }
}

declare module '@likecoin/wallet-connector' {
  interface LikeCoinWalletConnectorConnectionResult {
    user?: {
      primary_email: string;
    };
  }
}
