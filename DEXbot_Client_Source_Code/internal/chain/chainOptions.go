package chain

import (
	"dexbot/internal/dexbotABI"
	"dexbot/internal/handler"
	"dexbot/internal/web3Logic"
	"encoding/json"
	"io/ioutil"
)

var Chains = make(map[string]*Chain)

func initializeChains() {
	nodeURLs := initializeNodeURLs()
	//Initialize BSC
	bscRPCNode := web3Logic.RPCNode{
		NodeURL:   nodeURLs.BSC,
		RPCClient: web3Logic.InitRPCClient(nodeURLs.BSC)}
	Chains["BSC"] = &Chain{
		RPCNode: bscRPCNode,
		//Pancakeswap Factory
		Factory: &_Factory{
			FactoryAddress:  "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73",
			FactoryInstance: web3Logic.GetContractInstance("0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73", dexbotABI.FactoryABI, &bscRPCNode)},
		//Wrapped BNB
		WrappedNativeTokenAddress: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",

		//Address of the on chain Abacus
		AbacusAddress: "0xc685148D517680aB6a5242B72d67D6dE8079F054",

		ChainID: uint64(56),
	}
}

func initializeNodeURLs() *NodeURLs {
	nodeURLsBytes, err := ioutil.ReadFile("internal/config/nodeURLs.json")
	handler.HandleError("Error when trying to read nodeURL.json", err)
	nodeURLs := &NodeURLs{}
	json.Unmarshal(nodeURLsBytes, nodeURLs)
	return nodeURLs
}
