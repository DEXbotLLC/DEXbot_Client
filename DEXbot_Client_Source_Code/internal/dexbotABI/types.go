package dexbotABI

import (
	"github.com/shopspring/decimal"
)

type AbacusSwapTx struct {
	LP              string
	AmountIn        decimal.Decimal
	AmountOutMin    decimal.Decimal
	TokenIn         string
	CustomAbacusFee bool
}
