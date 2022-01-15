package chainOptions

//Map for DEXbot to check the current available chain options
var ChainOptions = make(map[string]uint64)

//Map to identify the swap router address for each chain
var DEXRouters = make(map[string]string)

//Map to identify the wrapped native token address for each chain
var WrappedNativeTokens = make(map[string]string)

//Map to identify the address of the DEXbot Abacus
var AbacusAddresses = make(map[string]string)

//Initialize the available chain options DEXbot can connect to
func initializeChainIDMap() {

	//Binance Smart Chain
	ChainOptions["BSC"] = uint64(56)

	//Polygon
	ChainOptions["Polygon"] = uint64(137)
}

//Initialize the dex router addresses for each chain
func initializeDEXRouters() {

	//Binance Smart Chain
	DEXRouters["BSC"] = "0x10ED43C718714eb63d5aA57B78B54704E256024E"

	//Polygon
	DEXRouters["Polygon"] = ""

}

//Initialize the wrapped native token addresses for each chain
func initializeWrappedNativeTokens() {

	//Binance Smart Chain
	WrappedNativeTokens["BSC"] = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"

	//Polygon
	WrappedNativeTokens["Polygon"] = "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"

}

//Initialize the Abacus address for each chain
func initalizeAbacusAddresses() {

	//Binance Smart Chain
	AbacusAddresses["BSC"] = "0x0000000000009f100078627C52A022382b5f561D"

	//Polygon
	AbacusAddresses["Polygon"] = ""

}
