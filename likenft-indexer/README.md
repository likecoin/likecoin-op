# Project likenft-indexer

Indexes the `LikeProtocol` and `BookNFT` contracts, on Base in production, so
that downstream applications can query BookNFT classes, tokens, owners and
transfer memos without hitting the chain directly. It ingests EVM events, keeps
the NFT state in postgres and exposes it through an api. The chain and contract
address come from `ETH_NETWORK_EVENT_RPC_URL` and
`ETH_LIKE_PROTOCOL_CONTRACT_ADDRESS`, see `.env.example`.

## Pre-requisite

Install atlas db migration tool

```bash
curl -sSf https://atlasgo.sh | sh
```

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Development

Init env file

```bash
make secret
```

Start api server

```bash
make start-api
```

#### Workers

Tasks:

- `cmd/worker/task`

To start workers of all tasks

```bash
make start-worker
```

Start task scheduler

```bash
make start-scheduler
```

To stop workers of all tasks

```bash
make stop-worker
```

Stop task scheduler

```bash
make stop-scheduler
```

## MakeFile

Init env file

```bash
make secret
```

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
