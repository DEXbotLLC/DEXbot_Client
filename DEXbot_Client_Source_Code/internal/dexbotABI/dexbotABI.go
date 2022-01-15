package dexbotABI

import (
	"dexbot/internal/handler"

	"github.com/umbracle/go-web3/abi"
)

//Application Binary Interface for the DEX Router
var RouterABI *abi.ABI

//Application Binary Interface for tokens
var ERC20ABI *abi.ABI

//Initialize the router ABI
func initializeRouterABI() {
	web3ABI, err := abi.NewABI(RouterABIString)
	handler.HandleError("dexbotABI: initializeRouterABI: abi.JSON", err)
	RouterABI = web3ABI
}

//Decode transaction input data with the router ABI
func DecodeInputDataWithDEXRouterABI(methodName string, inputData []byte) map[string]interface{} {
	abiMethod := RouterABI.Methods[methodName]
	decodedInputData, err := abi.Decode(abiMethod.Inputs, inputData[4:])
	handler.HandleError("dexbotABI: DecodeInputDataWithDEXRouterABI: abi.Decode", err)
	return decodedInputData.(map[string]interface{})

}

//Initialize the ERC20 ABI to interact with token interfaces
func initializeERC20ABI() {
	web3ABI, err := abi.NewABI(ERC20ABIString)
	handler.HandleError("dexbotABI: initializeERC20ABI: abi.JSON", err)
	ERC20ABI = web3ABI
}

//Decode transaction input data with the ERC20 ABI
func DecodeInputDataWithERC20ABI(methodName string, inputData []byte) map[string]interface{} {
	abiMethod := ERC20ABI.Methods[methodName]
	decodedInputData, err := abi.Decode(abiMethod.Inputs, inputData[4:])
	handler.HandleError("dexbotABI: DecodeInputDataWithERC20ABI: abi.Decode", err)
	return decodedInputData.(map[string]interface{})
}
