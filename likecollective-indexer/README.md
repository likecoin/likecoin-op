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

`accounts`, `nft_classes` and `stakings` are re-read from `LikeCollective`
whenever an event touches them, rather than moved by that event's delta. An
event is the trigger to re-read an amount, not the source of it.

That is what makes the numbers recoverable. The reads are deliberately not
pinned to the event's block: nothing orders the pipeline -- events are enqueued
oldest-first but asynq retries land late, and `retry-failed-evm-events`
re-drives old events on purpose -- so a pinned read would let an older event
commit an older truth over a newer one and leave it there. Reading the head,
every writer converges on the same answer whatever order they run in, and no
archive node is needed. The head is resolved once per event and every read in
that pass is pinned to it, so a pool-wide re-read cannot mix rows from either
side of a block that lands mid-pass; and an event whose block the node has not
reached yet fails and is retried rather than read against an older head.
`staking_events` is still accumulated from deltas -- it is history, and history
has to be. The deltas are neither applied to nor checked against the loaded
totals: those may already have been re-read at a head past the event, so an
unstake larger than the loaded amount is expected, and refusing it would fail
the event on every retry.

The totals used to be accumulated too, which is why they drifted: a missing or
wrong event did not delay a number, it offset it, and every later event built
on the wrong base. `RewardDeposited` was worse than that, re-deriving the
per-staker split with one integer floor per account where the contract takes
one per position against a reward index -- so the two disagreed by dust on
every single deposit, by construction rather than by accident.

`claimed_reward_amount` is the exception. It is lifetime cumulative and the
contract keeps no such counter, so it stays accumulated and stays
unrecoverable.

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
an event applied between its read and its commit is overwritten. Stop the
workers before `--apply`. It refuses a snapshot block older than a staking
event already applied, because its default of head less `--confirmations`
would otherwise roll back events that are never replayed, and one at or after
a staking event still to be applied (received, enqueued, or failed or
processing before persisting), which a worker would apply on top of the
snapshot again. Let the workers drain those first. A `processing` row left by a
crashed worker is never retried, so it keeps refusing until it is reset. Nor
can it see a webhook delivery that arrives late for a block the snapshot
already covers, during the write or after it; `--confirmations` is the only
margin against that, and running `resync` again at a later block repairs it. Run it when
`check-evm-event-gaps` says something was lost, and prefer running it in
cluster -- the mainnet write takes ~80s there against ~35 minutes over a
port-forward.

### What none of this repairs

- **`claimed_reward_amount`** is lifetime-cumulative and the contract keeps no
  such counter, so neither the per-event re-read nor a snapshot can rebuild it.
  Only replaying `RewardClaimed` from logs can.
- **`staking_events`** is history, not current state, so `resync` leaves it
  alone. Anything reading it -- including the TimescaleDB
  `book_nft_delta_time_bucket_*` aggregates behind
  `GET /collective/api/book-nfts/{7d|30d|1y}/delta` -- still reflects the gap.
- **`nft_classes.last_staked_at`** is not part of a snapshot; rows a resync
  creates get the ent schema default.

Two drift cases the re-read cannot heal by itself. A row that is missing
outright -- its `Staked` delivery was lost, so nothing ever created it -- fails
an unstake or claim with "not found" before the re-read is reached, and the
event retries into `retry-failed-evm-events` until a `cli resync` creates the
row. That is loud rather than silent: the events pile up at `failed` and sentry
says so. And an account total already lower than a correction takes out of one
of its stakings is clamped to zero rather than wrapped, which is logged as an
error and discards what the account holds in pools the event did not load; it
also takes a `cli resync`, which rebuilds the account from every pool.

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
