package config

import (
	"github.com/kelseyhightower/envconfig"
)

type contextKey struct{}

var ContextKey = &contextKey{}

type EnvConfig struct {
	RedisDsn string `envconfig:"REDIS_DSN" default:"redis://127.0.0.1:6379"`

	SentryDsn   string `envconfig:"SENTRY_DSN" default:""`
	SentryDebug bool   `envconfig:"SENTRY_DEBUG" default:"false"`

	// Worker config
	Concurrency int `envconfig:"WORKER_CONCURRENCY" default:"1"`

	// Logic config
	EthNetworkEventRPCURL          string `envconfig:"ETH_NETWORK_EVENT_RPC_URL"`
	EthNetworkPublicRPCURL         string `envconfig:"ETH_NETWORK_PUBLIC_RPC_URL"`
	EthLikeProtocolContractAddress string `envconfig:"ETH_LIKE_PROTOCOL_CONTRACT_ADDRESS"`

	EvmEventLikeProtocolInitialBlockHeight uint64 `envconfig:"EVM_EVENT_LIKE_PROTOCOL_INITIAL_BLOCK_HEIGHT" default:"1"`
	EvmEventQueryNumberOfBlocksLimit       uint64 `envconfig:"EVM_EVENT_QUERY_NUMBER_OF_BLOCKS_LIMIT"`
	BookNFTEventsQueryBatchSize            int    `envconfig:"BOOK_NFT_EVENTS_QUERY_BATCH_SIZE" default:"10"`
}

func LoadEnvConfigFromEnv() (*EnvConfig, error) {
	var cfg EnvConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
