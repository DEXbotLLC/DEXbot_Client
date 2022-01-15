package chainOptions

//Initialize the chainOptions package
func Initialize() {

	//Initialize the available chain options DEXbot can connect to
	initializeChainIDMap()

	//Initialize the dex router addresses
	initializeDEXRouters()

	//Initialize the wrapped native token addresses
	initializeWrappedNativeTokens()

	//Initialize the Abacus address for each chain
	initalizeAbacusAddresses()
}
