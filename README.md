# Likecoin 3.0

This repo contains the code for Likecoin 3.0, which is hosted on the
[Base](https://base.org) chain, migrated from Optimism.

## Live contracts

Base mainnet (chain id `8453`) proxy addresses:

| Contract | Address |
| --- | --- |
| `Likecoin` (ERC20) | `0x1EE5DD1794C28F559f94d2cc642BaE62dC3be5cf` |
| `LikeProtocol` | `0xfb5cbb1973a092E6C77af02EA1E74B14870AbeC5` |
| `LikeCollective` | `0x4506Ac2dD1e9A470d92a3D1656E1a99C676E1c8E` |
| `LikeStakePosition` | `0x508610D3009cda82Ac1a40D2b322Ed31932D16b1` |
| `veLike` | `0xE55C2b91E688BE70e5BbcEdE3792d723b4766e2B` |

`BookNFT` is not deployed directly, each book is a beacon proxy created through
`LikeProtocol`.

veLike reward pools are rotated and new ones are added over time, so they are
deliberately not listed here. See [likecoin3/velike.md](likecoin3/velike.md) for
the current pool and the rotation procedure.

For everything else, read
`likecoin3/ignition/deployments/chain-{chainId}/deployed_addresses.json`, where
the `#ERC1967Proxy` entry of each `*Module` is that module's proxy address.
Only Base (`8453`) and Base Sepolia (`84532`) carry the full stack. The other
chain folders, Optimism (`10`), Optimism Sepolia (`11155420`), Ethereum (`1`),
Sepolia (`11155111`), Unichain (`130`) and Unichain Sepolia (`1301`), only carry
the `Likecoin` ERC20. The legacy Optimism `LikeProtocol` addresses live in
`likenft/.openzeppelin/` instead.

## Folder overview

### Smart contracts

- [`likecoin3`](likecoin3/README.md): the current contracts, Hardhat + Ignition,
  deployed with CREATE2 so that addresses stay consistent across superchain
  deployments. Includes `Likecoin` (ERC20), `LikeProtocol` and `BookNFT`,
  `LikeCollective` and `LikeStakePosition`, plus `veLike` and its reward pools.
- [`likenft`](likenft/README.md): the legacy Optimism-era LikeNFT / LikeProtocol
  contracts. Superseded by `likecoin3`, still used to run the local `eth-node`
  in docker compose and kept for the migration path.
- `abi`: ABIs and bytecode generated from `likecoin3`, consumed by the Go
  services. Regenerate with `make abigen`. `BeaconProxy.*` is the exception, it
  is copied by hand from the OpenZeppelin artifacts.

### Indexers and services

- [`likenft-indexer`](likenft-indexer/README.md): indexes NFT data around
  `LikeProtocol` and `BookNFT` for downstream applications to query.
- [`likecollective-indexer`](likecollective-indexer/README.md): indexes
  `LikeCollective` and `LikeStakePosition` staking data, serves
  `indexer-base.v3.like.co`.
- [`likecollective-staking-position-image`](likecollective-staking-position-image/README.md):
  renders the staking position NFT image, serves `asset.v3.like.co`.
- `signer-backend`: holds the migration hot wallet and signs the EVM
  transactions, segregated from the backend that holds the logic.

### Migration

- [`likecoin-migration`](likecoin-migration/README.md): the migration frontend,
  from cosmos LIKE to EVM LIKE.
- [`likenft-migration`](likenft-migration/README.md): the migration frontend for
  likerID and likeNFT, from cosmos x/nft to EVM likenft.
- [`migration-backend`](migration-backend/README.md): the backend program
  migrating the data from cosmos LIKE and cosmos NFT.
- [`migration-admin`](migration-admin/README.md): simple CRUD interface for CX to
  support end user migration.
- [`op-2-base`](op-2-base/README.md): tooling and runbook for the one-off
  Optimism to Base asset migration.

### Operation

- `operation`: scripts for operation support, local state init, funding the
  operator wallet and batch pre-mint.
- [`deploy`](deploy/README.md): the helm chart and scripts for deploying the
  services.
- [`cosmos-delegation`](cosmos-delegation/README.md): scripts for computing
  cosmos delegation commands during and after the migration period.
- `nft-tools`: cosmos-era NFT scripts, mint and send, kept for reference.

## Quick start up on backend

To setup the env files and kick start basic backend for the first time
```
make setup
make start
```

The infra will be setup with the following services

- migration-backend
- migration-backend-worker
- migration-backend-scheduler
- signer-backend
- eth-node
    - a hardhat node, running with the `likenft` hardhat config
- db-migration-backend
- db-signer-backend
- redis

`make start` also runs the db migrations and deploys the contracts to the local
node, then streams the backend logs.

To redeploy the smart contracts to the local node afterwards, run the following
in another console
```
make local-contracts
```

To stop services,

```
$ make stop
```

> Please note that the `eth-node` have no mechanism to persist state.
> To prevent stucking at nonce inconsistent, after restarting the eth-node for any reason,
> the dbs related to transactions should also be reset.
>
> ```
> $ make clean-transaction-volumes
> ```
>
> The stop command above has included this command

To wipe every volume, including the databases,

```
$ make clean-docker-volumes
```

> Both clean targets match volumes by the `likecoin-30` compose project prefix,
> so they only work when the repo is checked out into a `likecoin-3.0`
> directory. Under any other directory name they silently match nothing.

## Setup on frontends

For respective frontend, `likenft-migration`, `likecoin-migration` and
`migration-admin`, please navigate to respective folder and follow instruction
there.

## Working on the contracts

Contract development, local superchain setup, deployment and simulation are
documented in [likecoin3/README.md](likecoin3/README.md). See also
[likecoin3/DEVCMD.md](likecoin3/DEVCMD.md) for the local hardhat task reference
and [likecoin3/velike.md](likecoin3/velike.md) for veLike deployment and reward
rotation.

After changing the contracts, regenerate the shared ABIs so that the Go services
pick them up:

```
make abigen
```

This rebuilds `likecoin3`, refreshes `abi/`, and reruns `abigen` in
`likenft-indexer`, `migration-backend` and `likecollective-indexer`.

## Other make targets

| Target | Description |
| --- | --- |
| `make docker-images` | Build and push the docker images of every service |
| `make deploy` | Decrypt the secrets and deploy through `deploy/` |
| `make operator-key-link` | Symlink `deploy/env.operator` into `likenft/.env`, legacy, `likecoin3` reads its own `.env` |
| `make remove-operator-key-link` | Remove the symlink above |

## CI

`.github/workflows/ci.yaml` runs `make -C {folder} ci` for each service folder.
`likenft-indexer` is the exception, it has its own workflow in
`.github/workflows/likenft-indexer.yaml`.
