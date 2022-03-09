package chain

import (
	"dexbot/internal/web3Logic"

	"github.com/umbracle/go-web3/contract"
)

type ChainOptions struct {
	Polygon  Chain
	Ethereum Chain
	BSC      Chain
	Optimism Chain
	Arbitrum Chain
}

type Chain struct {
	RPCNode                   web3Logic.RPCNode
	AbacusAddress             string
	Factory                   *_Factory
	WrappedNativeTokenAddress string
	ChainID                   uint64
}

type DEX struct {
	Router  _Router
	Factory _Factory
}

type NodeURLs struct {
	Polygon  string `json:"Polygon"`
	Arbitrum string `json:"Arbitrum"`
	Optimism string `json:"Optimism"`
	BSC      string `json:"BSC"`
}

type _Router struct {
	//address

	//contract instance
}

type _Factory struct {
	//address
	FactoryAddress  string
	FactoryInstance *contract.Contract
}
