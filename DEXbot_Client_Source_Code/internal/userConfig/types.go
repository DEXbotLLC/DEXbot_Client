package userConfig

import (
	"dexbot/internal/web3Logic"

	"github.com/shopspring/decimal"
	"github.com/umbracle/go-web3/contract"
)

type UserConfiguration struct {
	ErrorReporting bool `mapstructure:"error_reporting"`
	Wallets        map[string]*UserWallet
}

type UserWallet struct {
	WalletAddress string
	WalletName    string
	Tokens        map[string]*Token
	PrivateKey    *web3Logic.Key
}

type Token struct {
	TokenChecksumAddress string
	TokenInstance        *contract.Contract
	Decimals             int
	LPChecksumAddress    string
	LPInstance           *contract.Contract
	LPToken0Address      string
	TokenConfig          *TokenConfiguration
}

type TokenConfiguration struct {
	ChainName                        string
	TokenName                        string
	ContractAddress                  string
	LastSellTimestamp                int64
	MinimumUSDSellAmount             decimal.Decimal
	MaximumUSDSellAmount             decimal.Decimal
	Minimum24HourPriceChange         float64
	MinimumTimeInSecondsBetweenSells int64
	MinimumTokenBalance              decimal.Decimal
	MinimumUSDSellPrice              float64
}
