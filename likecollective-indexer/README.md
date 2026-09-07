# Likecollective indexer

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Development

Init env file

```bash
make secret
```

Install tools

```bash
make vendor
```

Start api server

```bash
make start-api
```

Linting

```bash
make lint
```

Stop api server

```bash
make stop-api
```

Run the test suite:

```bash
make test
```

## MakeFile

Build the application

```bash
make build
```

Clean up binary from the last build:

```bash
make clean
```

### Alchemy config

Quick reference for setting up webhook in the alchemy webhook

URL For LikeCollective contract should be
`https://indexer-base.v3.like.co/collective/alchemy/like-collective/ethlog`

Webhooktype: Custom
Query as below: (all event)
```
{
  block {
    hash,
    number,
    timestamp,
    logs(filter: {addresses: ["0x4506Ac2dD1e9A470d92a3D1656E1a99C676E1c8E"]}) {
      # Account is the account which generated this log - this will always be a contract account.
      account {
        address,
      },
      # Topics is a list of 0-4 indexed topics for the log.
      topics,
      # Data is unindexed data for this log.
      data,
      # Transaction is the transaction that generated this log entry.
      transaction {
        hash,
        index,
      },
      index,
    }
  }
}
```

URL for LikeStakePosition should be
`https://indexer-base.v3.like.co/collective/alchemy/like-stake-position/ethlog`
Webhooktype: Custom
Query as below: (all event)
```
{
  block {
    hash,
    number,
    timestamp,
    logs(filter: {addresses: ["0x508610D3009cda82Ac1a40D2b322Ed31932D16b1"]}) {
      # Account is the account which generated this log - this will always be a contract account.
      account {
        address,
      },
      # Topics is a list of 0-4 indexed topics for the log.
      topics,
      # Data is unindexed data for this log.
      data,
      # Transaction is the transaction that generated this log entry.
      transaction {
        hash,
        index,
      },
      index,
    }
  }
}
```

## Simulation

> To prepare the logs, please take a look at [likecoin3](../likecoin3/README.md)

### Simulate with local node

1. Retrieve the like collective address and like stake position address after local contract is deployed.

    ```bash
    go run ./cmd/cli simulate {likecoin3_simulation_output.json} --rpc http://localhost:8545 --like-collective-address {0xlike_collective_address} --like-stake-position-address {0xlike_stake_position_address} | jq
    go run ./cmd/cli simulate {likecoin3_simulation_output.json} --rpc http://localhost:8545 --like-collective-address {0xlike_collective_address} --like-stake-position-address {0xlike_stake_position_address} --verify
    ```

    e.g.

    ```bash
    go run ./cmd/cli simulate ../likecoin3/simulate/likecollective/simulations/case1.output.json --rpc http://localhost:8545 --like-collective-address 0x227eFaea699FDe0B11e52d7AF81ebF215c5532E1 --like-stake-position-address 0x75C986519A6F4a144520aAd0AFF7846d17e66175 | jq
    go run ./cmd/cli simulate ../likecoin3/simulate/likecollective/simulations/case1.output.json --rpc http://localhost:8545 --like-collective-address 0x227eFaea699FDe0B11e52d7AF81ebF215c5532E1 --like-stake-position-address 0x75C986519A6F4a144520aAd0AFF7846d17e66175 --verify
    ```

2. Compare the result of `{likecoin3_simulation_output.json}` and the console output and see if there are inconsistency.

## Keeping the indexed state true to the chain

`accounts`, `nft_classes` and `stakings` hold running totals accumulated from
events, not values derived from the chain. A missing or wrong event therefore
does not delay a number, it offsets it, and every later event builds on the
wrong base. Two things guard that.

### Noticing

- **`check-evm-event-gaps`** (scheduler, `*/5`) reads the logs the contracts
  actually emitted in a recent window and compares them with `evm_events`. The
  webhook is the only ingest path, so this is the only check that can see a
  delivery that never arrived -- `check-received-evm-events` drains rows the
  webhook already wrote and cannot miss what was never written. A non-empty
  result is returned as an error so the sentry middleware raises it. The window
  is `EVM_EVENT_GAP_QUERY_NUMBER_OF_BLOCKS_LIMIT` blocks (500 by default,
  matching what providers accept for one `eth_getLogs`); keep the cron well
  inside the time the chain takes to produce that many blocks.
- **`retry-failed-evm-events`** (scheduler, `*/15`) re-drives events that
  exhausted their asynq retries and were parked at `failed`. Without it a
  single failed event stayed failed forever.

### Repairing

`resync` rebuilds all three tables from one pinned-block on-chain snapshot:

```bash
go run ./cmd/cli resync                  # report the diff, write nothing
go run ./cmd/cli resync --apply          # write it
go run ./cmd/cli resync --report out.json
```

It enumerates every live position through `LikeStakePosition`, unions those
pairs with the rows already stored so a burned position is zeroed rather than
left at its last value, and refuses to write when the contract's two
independent accessors disagree with each other.

It is deliberately **manual**. Running it on a schedule would have it rewriting
financial columns unattended, and it takes no lock against the live worker, so
an event applied between its read and its commit is overwritten. Run it when
`check-evm-event-gaps` says something was lost, and prefer running it in
cluster -- the mainnet write takes ~80s there against ~35 minutes over a
port-forward.

### What none of this repairs

- **`claimed_reward_amount`** is lifetime-cumulative and the contract keeps no
  such counter, so a snapshot cannot rebuild it. Only replaying `RewardClaimed`
  from logs can.
- **`staking_events`** is history, not current state, so `resync` leaves it
  alone. Anything reading it -- including the TimescaleDB
  `book_nft_delta_time_bucket_*` aggregates behind
  `GET /collective/api/book-nfts/{7d|30d|1y}/delta` -- still reflects the gap.
- **`nft_classes.last_staked_at`** is not part of a snapshot; rows a resync
  creates get the ent schema default.

`check-evm-event-gaps` has two blind spots of its own. It only looks at a
window behind the head, so if the scheduler itself is down for longer than that
window, the interval it skipped is never examined by any later pass -- it keeps
no cursor. And it only reports logs the chain has that the database does not;
a row the webhook delivered from a block later reorged out stays put and is not
flagged. The block padding makes the second unlikely rather than impossible.

Expect the first run in a new environment to be loud: the two Alchemy webhooks
are registered out of band, so if either is missing, paused, or filtered on a
stale address, every log from that contract is correctly reported missing. Run
it somewhere quiet before pointing it at sentry.

Replaying missing logs back into `evm_events` would repair the first two, and
`check-evm-event-gaps` already identifies exactly which logs to replay. It is
not implemented yet because it is not safe yet: `RewardDeposited` fans out into
a per-staker split derived from the stake distribution held **at the time it is
applied**, so a log replayed late would distribute against today's distribution
rather than the one at its own block. That event has to read its amounts from
the chain at its own block first.
