# Go Parrot Twap

Go Parrot Twap will execute buy or sell orders over a specific time interval

## Getting started

Firstly copy the `.env.example` to `.env` and change set your wallet Private Key to the `WALLETPK` env variable.

You can also switch to your own Solana RCP cluster by changing the `` env variable.


Finally for **Buying** you can run the Twap using this command:

```sh
go run cmd/cli.go --side buy --pair PRT:SOL --amount 0.1 --interval 10m
```
> That translate to: for the pair `PRT:SOL`, `Buy` PRT with `0.1` SOL every `10m` (minutes)


For **Selling** will be:

```sh
go run cmd/cli.go --side sell --pair PRT:SOL --amount 100 --interval 10m
```
> That translate to: for the pair `PRT:SOL`, `Sell` `100` PRT for SOL every `10m` (minutes)


Optional you can specify a `TargetAmount` params with the `--target` argument to tell the Twap to stop buying or selling when the target amount is reached, for example:

```sh
go run cmd/cli.go --side buy --pair PRT:SOL --amount 0.001 --target 100 --interval 30s
```

It will buy 0.001 SOL worth of PRT every 30 seconds and it will **stop** buying PRT when the balance in the current wallet reach 100 PRT

## Production

For production you can run `make` and run `build/twap` or use docker container:

```sh
// First build
docker-compose build
// Then run
docker run twap --side buy --pair PRT:SOL --amount 0.001 --target 100 --interval 30s
```