import i18n from './i18n.config';

export default {
  // Disable server-side rendering: https://go.nuxtjs.dev/ssr-mode
  ssr: false,

  // Target: https://go.nuxtjs.dev/config-target
  target: 'static',

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    title: 'LikeCoin v3 Migration - Upgrade Your Tokens',
    htmlAttrs: {
      lang: 'en',
    },
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      {
        hid: 'description',
        name: 'description',
        content:
          'Upgrade your LikeCoin v2 tokens to v3. Connect your wallet to migrate your Liker ID and access the latest features. Join the 10M LIKE rewards pool by staking on 3ook.com.',
      },
      {
        hid: 'og:title',
        property: 'og:title',
        content: 'LikeCoin v3 Migration - Upgrade Your Tokens',
      },
      {
        hid: 'og:description',
        property: 'og:description',
        content:
          'Upgrade your LikeCoin v2 tokens to v3. Connect your wallet to migrate your Liker ID and access the latest features. Join the 10M LIKE rewards pool by staking on 3ook.com.',
      },
      {
        hid: 'og:image',
        property: 'og:image',
        content: 'https://v3-token-migration.like.co/cover.png',
      },
      {
        hid: 'og:url',
        property: 'og:url',
        content: 'https://v3-token-migration.like.co',
      },
      {
        hid: 'og:type',
        property: 'og:type',
        content: 'website',
      },
      { name: 'format-detection', content: 'telephone=no' },
    ],
    link: [
      { rel: 'icon', type: 'image/png', href: '/favicon.png' },
      { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
      {
        rel: 'preconnect',
        href: 'https://fonts.gstatic.com',
        crossorigin: true,
      },
      {
        rel: 'stylesheet',
        href: 'https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Open+Sans:ital,wght@0,300..800;1,300..800&display=swap',
      },
    ],
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [
    './assets/css/style.css',
    './assets/fontawesome/fontawesome.css',
    '@likecoin/wallet-connector/dist/style.css',
    '@likecoin/evm-wallet-connector/style.css',
  ],

  i18n,

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [
    { src: '~/plugins/textEncoder' },
    '~/plugins/config.ts',
    '~/plugins/apiClient.ts',
    '~/plugins/cosmosNetworkConfig.ts',
    '~/plugins/evmChainConfig.ts',
    '~/plugins/likeCoinWalletConnector.ts',
    '~/plugins/likeCoinEVMWalletConnector.ts',
    '~/plugins/intercom.ts',
    '~/plugins/gtag.ts',
  ],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    // https://go.nuxtjs.dev/typescript
    '@nuxt/typescript-build',
    // https://go.nuxtjs.dev/tailwindcss
    '@nuxtjs/tailwindcss',
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: ['@nuxtjs/i18n', '@nuxtjs/axios'],

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {
    babel: {
      presets: [
        [
          '@nuxt/babel-preset-app',
          {
            // Disable bable transform-exponentiation-operator
            // https://github.com/starknet-io/starknet.js/issues/37#issuecomment-955797303
            targets: {
              browsers: [
                'chrome >= 67',
                'edge >= 79',
                'firefox >= 68',
                'opera >= 54',
                'safari >= 14',
              ],
            },
          },
        ],
      ],
    },
    transpile: [
      '@cosmjs',
      'cosmjs-types',
      '@walletconnect',
      '@web3modal',
      'unstorage',
      '@likecoin/wallet-connector',
      'libsodium',
      'web3',
      'ethers',
    ],
  },
};
