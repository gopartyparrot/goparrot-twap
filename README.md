# Go Parrot Twap

Go Parrot Twap will execute buy or sell orders over a specific time interval.

## Getting started

Firstly copy the `.env.example` to `.env` and set your wallet Private Key in the `WALLETPK` env variable.

You can also switch to your own Solana RCP cluster by changing the `RPCURL` and `RPCWS` env variable.


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

### Stop buying

Optional you can specify a `StopAmount` params with the `--stopAmount` argument to tell the Twap to stop buying or selling when the given amount, in the balance, is reached, for example:

```sh
go run cmd/cli.go --side buy --pair PRT:SOL --amount 0.001 --stopAmount 100 --interval 30s
```

It will buy 0.001 SOL worth of PRT every 30 seconds and it will **stop** buying PRT when the balance in the wallet reach 100 PRT (aka stopAmount)

### Transfer amount

Optional you can specify a `TransferAddress` and `TransferThreshold` with the params  `--transferAddress` and `--transferThreshold` respectively. 

When the balance of the buying asset reach the `TransferThreshold`, all it's balance will be transfer to the `TransferAddress` associated token account.

Example:

```sh
go run cmd/cli.go --side buy --pair PRT:SOL --amount 0.37722  --interval 10m --transferThreshold 100000 --transferAddress FRnCC8dBCcRabRv8xNbR5WHiGPGxdphjiRhE2qJZvwpm
```

It will buy 0.37722 SOL worth of PRT every 10 minutes and it will transfer to the Parrot Protocol address all the PRT balance once greater than 100,000 PRT

## Production

For production you can run `make` and run `build/twap`.

Or you can use docker container:

```sh
# First build
docker-compose build
# Then run
docker run twap --side buy --pair PRT:SOL --amount 0.001 --target 100 --interval 30s
```
