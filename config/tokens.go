package config

import "math"

const (
	TokenAccountSize = 165
	NativeSOL        = "11111111111111111111111111111111"
	WrappedSOL       = "So11111111111111111111111111111111111111112"
)

type TokenInfo struct {
	Symbol   string
	Decimals uint8
}

func (s *TokenInfo) Pow() float64 {
	return math.Pow10(int(s.Decimals))
}

func (s *TokenInfo) ToFloat(v uint64) float64 {
	return float64(v) / math.Pow10(int(s.Decimals))
}

func (s *TokenInfo) FromFloat(v float64) uint64 {
	return uint64(v * math.Pow10(int(s.Decimals)))
}

var TokenInfos = map[string]TokenInfo{
	NativeSOL: {
		Symbol:   "SOL",
		Decimals: 9,
	},
	// Parrot tokens
	"Ea5SjE2Y6yvCeW5dYTn7PYMuW5ikXkvbGdcmSnXeaLjS": {
		Symbol:   "PAI",
		Decimals: 6,
	},
	"PRT88RkA4Kg5z7pKnezeNH4mafTvtQdfFgpQTGRjz44": {
		Symbol:   "PRT",
		Decimals: 6,
	},
	"E2Ub8wPfxxEvdrtumbfeL2HaQHgpd3gUGkDxDmmgN3p9": {
		Symbol:   "PTT",
		Decimals: 9,
	},
	"DYDWu4hE4MN3aH897xQ3sRTs5EAjJDmQsKLNhbpUiKun": {
		Symbol:   "pBTC",
		Decimals: 8,
	},
	"9EaLkQrbjmbbuZG9Wdpo8qfNUEjHATJFSycEmw6f1rGX": {
		Symbol:   "pSOL",
		Decimals: 9,
	},
	"BdZPG9xWrG3uFrx2KrUW1jT4tZ9VKPDWknYihzoPRJS3": {
		Symbol:   "prtSOL",
		Decimals: 9,
	},
	// Main Tokens
	"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {
		Symbol:   "USDC",
		Decimals: 6,
	},
	"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB": {
		Symbol:   "USDT",
		Decimals: 6,
	},
	WrappedSOL: {
		Symbol:   "WSOL",
		Decimals: 9,
	},
	"SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt": {
		Symbol:   "SRM",
		Decimals: 6,
	},
	"MERt85fc5boKw3BW1eYdxonEuJNvXbiMbs6hvheau5K": {
		Symbol:   "MER",
		Decimals: 6,
	},
	"9n4nbM75f5Ui33ZbPYXn59EwSgE8CGsHtAeTH5YFeJ9E": {
		Symbol:   "BTC(Sollet)",
		Decimals: 6,
	},
	"CDJWUqTcYTVAKXAVXoQZFes5JUFc7owSeq7eMQcDSbo5": {
		Symbol:   "renBTC",
		Decimals: 8,
	},
	"mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So": {
		Symbol:   "mSOL",
		Decimals: 9,
	},
	// MER LPs
	"57h4LEnBooHrKbacYWGCFghmrTzYPVn8PwZkzTzRLvHa": {
		Symbol:   "MER LP USDC-USDT-UST",
		Decimals: 9,
	},
	"9s6dXtMgV5E6v3rHqBF2LejHcA2GWoZb7xNUkgXgsBqt": {
		Symbol:   "MER LP USDC-USDT-PAI",
		Decimals: 6,
	},
	"GHhDU9Y7HM37v6cQyaie1A3aZdfpCDp6ScJ5zZn2c3uk": {
		Symbol:   "MER LP SOL-pSOL",
		Decimals: 9,
	},
	// Saber LPs
	"PrsVdKtXDDf6kJQu5Ff6YqmjfE4TZXtBgHM4bjuvRnR": {
		Symbol:   "SBR LP prtSOL-SOL",
		Decimals: 9,
	},
	"SLPbsNrLHv8xG4cTc4R5Ci8kB9wUPs6yn6f7cKosoxs": {
		Symbol:   "SBR LP BTC-renBTC",
		Decimals: 8,
	},
	"SoLEao8wTzSfqhuou8rcYsVoLjthVmiXuEjzdNPMnCz": {
		Symbol:   "SBR LP mSOL-SOL",
		Decimals: 9,
	},
	"2poo1w1DL6yd2WNTCnNTzDqkC6MBXq7axo77P16yrBuf": {
		Symbol:   "SBR LP USDC-USDT",
		Decimals: 6,
	},
	"UST32f2JtPGocLzsL41B3VBBoJzTm1mK1j3rwyM3Wgc": {
		Symbol:   "SBR LP UST-USDC",
		Decimals: 9,
	},
	// Raydium LPs
	"8HoQnePLqPj4M7PUDzfw8e3Ymdwgc7NLGnaTUapubyvu": {
		Symbol:   "RAY LP SOL-USDC",
		Decimals: 9,
	},
	"3H9NxvaZoxMZZDZcbBDdWMKbrfNj7PCF5sbRwDr7SdDW": {
		Symbol:   "RAY LP MER-USDC",
		Decimals: 6,
	},
}
