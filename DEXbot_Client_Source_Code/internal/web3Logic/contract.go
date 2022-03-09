package web3Logic

import (
	"dexbot/internal/dexbotABI"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"github.com/umbracle/go-web3/contract"
)

func CheckIfSmartContract(checksumAddress string, rpcNode *RPCNode) bool {
	code, err := rpcNode.RPCClient.Eth().GetCode(web3.HexToAddress(checksumAddress), web3.Latest)
	handler.HandleError("contract: CheckIfSmartContract: GetCode", err)
	if code == "0x" {
		return false
	} else {
		return true
	}

}

func GetContractInstance(tokenChecksumAddress string, abi *abi.ABI, rpcNode *RPCNode) *contract.Contract {
	web3Address := web3.HexToAddress(tokenChecksumAddress)

	_contractInstance := contract.NewContract(web3Address, abi, rpcNode.RPCClient)
	if _contractInstance != nil {
		return _contractInstance
	} else {
		GetContractInstance(tokenChecksumAddress, abi, rpcNode)
	}

	dexbotUtils.YellowPrinter.Println("There was an issue with the node, please restart and try again.")
	handler.Exit("")

	return nil
}

func GetLPContractAddress(tokenChecksumAddress string, wnatoChecksumAddress string, factoryContract *contract.Contract) string {
	//init lp contract, passes in addresses for tokenA(contractAddress) and tokenB(wnato) as args
	lpContractMap, err := factoryContract.Call("getPair", web3.Latest, web3.HexToAddress(tokenChecksumAddress), web3.HexToAddress(wnatoChecksumAddress))
	handler.HandleError("web3Logic: GetLPContractAddress: factoryContract.Call", err)
	lpContractAddress := lpContractMap["0"].(web3.Address)
	return lpContractAddress.String()
}

func GetLPContractInstance(tokenChecksumAddress string, wnatoChecksumAddress string, factoryContract *contract.Contract, rpcNode *RPCNode) *contract.Contract {
	//init lp contract, passes in addresses for tokenA(contractAddress) and tokenB(wnato) as args
	lpContractMap, err := factoryContract.Call("getPair", web3.Latest, web3.HexToAddress(tokenChecksumAddress), web3.HexToAddress(wnatoChecksumAddress))
	handler.HandleError("web3Logic GetLPContractInstance: factoryContract.Call", err)
	lpContractAddress := lpContractMap["0"].(web3.Address)
	contractInstance := GetContractInstance(lpContractAddress.String(), dexbotABI.LPABI, rpcNode)
	return contractInstance
}

func GetLPToken0Address(lpContract *contract.Contract) string {
	//init lp contract, passes in addresses for tokenA(contractAddress) and tokenB(wnato) as args
	lpToken0Map, err := lpContract.Call("token0", web3.Latest)
	handler.HandleError("web3Logic GetLPContractInstance: factoryContract.Call", err)
	lpToken0Address := lpToken0Map["0"].(web3.Address)
	return lpToken0Address.String()
}

func GetUSDPeggedLPContractAddress(tokenChecksumAddress string, wnatoChecksumAddress string, factoryContract *contract.Contract) string {
	//init lp contract, passes in addresses for tokenA(contractAddress) and tokenB(wnato) as args\
	lpContractMap, err := factoryContract.Call("getPair", web3.Latest, web3.HexToAddress(tokenChecksumAddress), web3.HexToAddress(wnatoChecksumAddress))
	handler.HandleError("web3Logic: GetLPContractAddress:factoryContract.Call", err)
	lpContractAddress := lpContractMap["0"].(web3.Address)
	return lpContractAddress.String()
}

func GetContractDecimals(contractInstance *contract.Contract) int {
	decimalsMap, err := contractInstance.Call("decimals", web3.Latest)
	handler.HandleError("web3Logic: GetContractDecimals: contractInstance.Call", err)
	var decimals int = int(decimalsMap["0"].(uint8))
	return decimals
}
