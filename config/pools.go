package config

import (
	_ "embed"
	"encoding/json"
)

const (
	RaydiumLiquidityPoolProgramIDV4 = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

//go:embed pools.json
var poolsBytes []byte

func GetPools() map[string]PoolConfig {
	var pools map[string]PoolConfig
	json.Unmarshal(poolsBytes, &pools)
	return pools
}

type PoolConfig struct {
	Service           string
	FromToken         string
	ToToken           string
	CoinGeckoID       string
	RaydiumPoolConfig RaydiumPoolConfig
}

type RaydiumPoolConfig struct {
	AmmId                 string
	AmmAuthority          string
	AmmOpenOrders         string
	AmmTargetOrders       string
	AmmQuantities         string
	PoolCoinTokenAccount  string
	PoolPcTokenAccount    string
	SerumProgramId        string
	SerumMarket           string
	SerumBids             string
	SerumAsks             string
	SerumEventQueue       string
	SerumCoinVaultAccount string
	SerumPcVaultAccount   string
	SerumVaultSigner      string
}
