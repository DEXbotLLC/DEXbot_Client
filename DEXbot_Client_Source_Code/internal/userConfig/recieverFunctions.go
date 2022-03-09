package userConfig

import (
	"dexbot/internal/web3Logic"
	"math"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/umbracle/go-web3"
)

func (t *Token) GetContractDecimals() int {
	return web3Logic.GetContractDecimals(t.TokenInstance)
}

func (t *Token) GetNatoPricePerToken(blockNumber web3.BlockNumber) decimal.Decimal {
	var price decimal.Decimal
	reserves := t.GetLPReserves(blockNumber)
	var longTokenValue *big.Int
	var longWNatoValue *big.Int

	if reserves["_reserve0"] == nil {
		return t.GetNatoPricePerToken(blockNumber - 1)
	}

	if strings.EqualFold(t.LPToken0Address, t.TokenChecksumAddress) {
		//extract token and wnato reserves
		longTokenValue = reserves["_reserve0"].(*big.Int)
		longWNatoValue = reserves["_reserve1"].(*big.Int)
	} else {
		//extract token and wnato reserves
		longWNatoValue = reserves["_reserve0"].(*big.Int)
		longTokenValue = reserves["_reserve1"].(*big.Int)
	}

	//convert to short form values
	wnatoValue := (decimal.NewFromBigInt(longWNatoValue, 0).Div(decimal.NewFromFloat(math.Pow(10, 18))))
	tokenValue := (decimal.NewFromBigInt(longTokenValue, 0).Div(decimal.NewFromFloat(math.Pow(10, float64(t.Decimals)))))

	//get the price per token in wnato
	price = wnatoValue.Div(tokenValue)
	return price
}

func (t *Token) GetLPReserves(blockNumber web3.BlockNumber) map[string]interface{} {
	reservesMap, err := t.LPInstance.Call("getReserves", blockNumber)
	//if there are no reserves at that block, retry at blockNumber-1
	if err != nil || reservesMap["_reserve0"] == nil || reservesMap["_reserve1"] == nil {
		return t.GetLPReserves(blockNumber - 1)
	}
	return reservesMap
}
