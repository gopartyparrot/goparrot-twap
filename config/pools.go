package config

const (
	SerumProgramIDV3                = "9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin"
	RaydiumLiquidityPoolProgramIDV4 = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

type PoolConfig struct {
	Service           string
	FromToken         string
	ToToken           string
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

var PoolConfigs = map[string]PoolConfig{
	"from_sol_to_prt_raydium": {
		Service:   "raydium_swap",
		FromToken: NativeSOL,                                     // SOL
		ToToken:   "PRT88RkA4Kg5z7pKnezeNH4mafTvtQdfFgpQTGRjz44", // PRT
		RaydiumPoolConfig: RaydiumPoolConfig{
			AmmId:                 "7rVAbPFzqaBmydukTDFAuBiuyBrTVhpa5LpfDRrjX9mr",
			AmmAuthority:          "5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1",
			AmmOpenOrders:         "7nsGyAGAawvpVF2JQRKLJ9PVwE64Xc2CzhbTukJdZ4TY",
			AmmTargetOrders:       "DqR8zK676oafdCMAtRm6Jc5d8ADQtoiUKnQb6DkTnisE",
			AmmQuantities:         NativeSOL,
			PoolCoinTokenAccount:  "Bh8KFmkkXZQzNgQ9qpjegfWQjNupLugtoNDZSacawGbb",
			PoolPcTokenAccount:    "ArBXA3NvfSmSDq4hhR17qyKpwkKvGvgnBiZC4K36eMvz",
			SerumProgramId:        SerumProgramIDV3,
			SerumMarket:           "H7ZmXKqEx1T8CTM4EMyqR5zyz4e4vUpWTTbCmYmzxmeW",
			SerumBids:             "5Yfr8HHzV8FHWBiCDCh5U7bUNbnaUL4UKMGasaveAXQo",
			SerumAsks:             "A2gckowJzAv3P2fuYtMTQbEvVCpKZa6EbjwRsBzzeLQj",
			SerumEventQueue:       "2hYscTLaWWWELYNsHmYqK9XK8TnbGF2fn2cSqAvVrwrd",
			SerumCoinVaultAccount: "4Zm3aQqQHJFb7Q4oQotfxUFBcf9FVP6qvt2pkJA35Ymn",
			SerumPcVaultAccount:   "B34rGhNUNxnSfxodkUoqYC3kGMdF4BjFHV2rQZAzQPMF",
			SerumVaultSigner:      "9ZGDGCN9BHiqEy44JAd1ExaAiRoh9HWou8nw44MbhnNX",
		},
	},
}
