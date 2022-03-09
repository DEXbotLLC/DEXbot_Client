package dexbotABI

import (
	"dexbot/internal/handler"
	"math/big"

	"github.com/shopspring/decimal"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
)

//Application Binary Interfaces
var AbacusABI *abi.ABI
var ERC20ABI *abi.ABI
var LPABI *abi.ABI
var FactoryABI *abi.ABI

//Initialize the Abacus ABI
func initializeAbacusABI() {
	web3ABI, err := abi.NewABI(AbacusABIString)
	handler.HandleError("There was an issue initializing the Abacus ABI", err)
	AbacusABI = web3ABI
}

//Decode Abacus Swap transaction input data with the Abacus ABI
func DecodeAbacusSwapInputData(methodName string, inputData []byte) *AbacusSwapTx {
	abiMethod := AbacusABI.Methods[methodName]
	decodedInputData, err := abi.Decode(abiMethod.Inputs, inputData[4:])
	handler.HandleError("There was an issue decoding the input data with the Abacus ABI", err)
	abacusSwapTx := MapToAbacusSwapTx(decodedInputData.(map[string]interface{}))
	return abacusSwapTx
}

func MapToAbacusSwapTx(inputData map[string]interface{}) *AbacusSwapTx {

	amountIn := inputData["_amountIn"].(*big.Int)
	amountOutMin := inputData["_amountOutMin"].(*big.Int)
	customAbacusFee := inputData["_customAbacusFee"].(bool)
	lp := inputData["_lp"].(web3.Address)
	tokenIn := inputData["_tokenIn"].(web3.Address)

	abacusSwapTx := &AbacusSwapTx{
		AmountIn:        decimal.NewFromBigInt(amountIn, 0),
		AmountOutMin:    decimal.NewFromBigInt(amountOutMin, 0),
		LP:              lp.String(),
		TokenIn:         tokenIn.String(),
		CustomAbacusFee: customAbacusFee,
	}

	return abacusSwapTx
}

//Initialize the ERC20 ABI to interact with token interfaces
func initializeERC20ABI() {
	web3ABI, err := abi.NewABI(ERC20ABIString)
	handler.HandleError("There was an issue initializing the ERC20 ABI", err)
	ERC20ABI = web3ABI
}

//Decode transaction input data with the ERC20 ABI
func DecodeInputDataWithERC20ABI(methodName string, inputData []byte) map[string]interface{} {
	abiMethod := ERC20ABI.Methods[methodName]
	decodedInputData, err := abi.Decode(abiMethod.Inputs, inputData[4:])
	handler.HandleError("There was an issue decoding the input data with the ERC20 ABI", err)
	return decodedInputData.(map[string]interface{})
}

//Initialize the Factory ABI
func initializeFactoryABI() {
	web3ABI, err := abi.NewABI(FactoryABIString)
	handler.HandleError("There was an issue initializing the Factory ABI", err)
	FactoryABI = web3ABI
}

//Initialize the LP ABI
func initializeLPABI() {
	web3ABI, err := abi.NewABI(LPABIString)
	handler.HandleError("There was an issue initializing the LP ABI", err)
	LPABI = web3ABI
}
