# Deployment Notes

This folder covers deploying the backend services. For deploying the smart
contracts themselves, see [likecoin3/README.md](../likecoin3/README.md), the
current contracts are deployed with Hardhat Ignition using the CREATE2 strategy.

We rely on operation to maintains the initial contract address of the evm smart contract.

Testnet operator keys are stored in the `env.operator` file. [Blackbox](https://github.com/StackExchange/blackbox) encrypts the keys.

To decrypt the keys, run `blackbox_decrypt_all_files`, then `make operator-key-link` at the root of the repo.

## Expected procedure

Current encrypted operator wallet: `0xC71fe89e4C0e5458a793fc6548EF6B392417A7Fb` (For testnet, for production, we use ledger `0xB0318A8f049b625dA5DdD184FfFF668Aa6E96261`, see `likenft/.env.optimism` for now)

The live deployment is on Base, the current proxy addresses are listed in the
[root README](../README.md) and in `likecoin3/ignition/deployments/`.

Legacy Optimism-era deployment, from the `likenft` folder, kept for reference:

- Optimism `LikeProtocol`: `0x526237a676444A67bc79E9009756df128Ca9a619`
- Optimism-sepolia `LikeProtocol`: `0xfF79df388742f248c61A633938710559c61faEF1`

Once the initial deployment is done, we should updated the expected proxy address in the `env.operator` file for reference. On new L2 chain launch, we should operate the following steps to obtain same proxy address.

0. Checkout a specific commit
1. Someone sends ETH to the operator's wallet (!!failed operation becasue of lack of fund will make the address different from the expected one)
2. Use operator's wallet to run the deployment from `likecoin3`, with the
   CREATE2 strategy and the shared `ignition/parameters.json`, so that the
   proxy addresses match the existing chains. See
   [likecoin3/README.md](../likecoin3/README.md) for which ignition module to
   use, the plain `LikeProtocol.ts` and `BookNFT.ts` modules are dev only.
   - 2.1. Deploy first, with the operator itself as the owner.
   - 2.2. Transfer owner to community wallet (0x91093818ef6A195fb9b453F8335a25Bd930DF4bd).

> These steps previously ran against `likenft` on Optimism, where the addresses
> were CREATE-based and therefore nonce dependent.

Formula refs:
- newAddressCREATE = `keccak256(deployingAddress ++ nonce)[12:]`
- newAddressCREATE2 = `keccak256(0xff ++ deployingAddress ++ salt ++ initCodeHash)[12:]`


## Migration signer hotwallet

Because of the natural of the project, some credential are manage independently for segregation.

For the signer private key, please consult @rickmak/@williamchong  For the htpassword, we generate separately
The `helm/{env}` value files and their Makefile are gitignored rather than
committed, see `.gitignore`, they provide helper and default for illustration.

For deployment to likeco

```
APP_VERSION=8c9106e4 ENV=likeco-base K8S_NAMESPACE=likecoin-op make deploy
```

`ENV` selects the `helm/{env}` value files, the known envs are `likeco-base`,
`likeco-op`, `pandawork-base` and `pandawork-op`. `K8S_NAMESPACE` defaults to
`likecoin-op`, see `Makefile`.

## Deployed services

The chart in `helm/` covers every backend service and frontend, see
`helm/templates/`. [`likecollective-staking-position-image`](../likecollective-staking-position-image/README.md)
is the exception, it has its own chart in
`../likecollective-staking-position-image/helm`.

To render or lint the chart before deploying,

```
ENV={env} make debug-template
ENV={env} make lint
```
