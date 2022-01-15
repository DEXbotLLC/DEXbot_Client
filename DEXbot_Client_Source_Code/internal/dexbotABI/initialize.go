package dexbotABI

func Initialize() {

	//Initialize the router ABI
	initializeRouterABI()

	//Initialize the ERC20 ABI to interact with token interfaces
	initializeERC20ABI()
}
