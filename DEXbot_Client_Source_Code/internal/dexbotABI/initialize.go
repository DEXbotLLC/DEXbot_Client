package dexbotABI

func Initialize() {

	//Initialize the Abacus ABI
	initializeAbacusABI()

	//Initialize the ERC20 ABI to interact with token interfaces
	initializeERC20ABI()

	//Initialize the LP ABI
	initializeLPABI()

	//Initialize the factory ABI
	initializeFactoryABI()
}
