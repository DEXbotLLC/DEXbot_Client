package userConfig

import (
	"dexbot/internal/chain"
	"dexbot/internal/database"
	"dexbot/internal/dexbotABI"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/web3Logic"
)

//Map to hold user config
var UserConfig *UserConfiguration

//Send a continous alive signal to the database with the user config data
func InitializePulse() {
	//send an alive signal to the database
	go database.SendPulse()
}

func SendVersionChecksum() {
	go database.SendVersionChecksum()
}

//Initialize the user config from the database
func initializeUserConfig() {
	_userConfig := UserConfiguration{}

	//Initialize user wallets
	_userConfig.Wallets = make(map[string]*UserWallet)

	//Get the user config from the database
	remoteUserConfig := database.GetUserConfig()

	if remoteUserConfig != nil {

		if remoteUserConfig["error_reporting"] != nil {

			//Set the local error reporting variable
			_userConfig.ErrorReporting = remoteUserConfig["error_reporting"].(bool)

			if remoteUserConfig["wallets"] != nil {

				//Unpack the remoteUserConfig into the local _userConfig
				for walletAddress, walletConfig := range remoteUserConfig["wallets"].(map[string]interface{}) {

					walletConfig := walletConfig.(map[string]interface{})

					//Initialize a user wallet
					_userWallet := UserWallet{}
					//Initialize wallet tokens
					_userWallet.Tokens = make(map[string]*Token)

					//Set wallet address to the wallet checksum address
					walletChecksumAddress, err := dexbotUtils.ToChecksumAddress(walletAddress)
					handler.HandleError("Error when trying to checksum wallet address", err)
					_userWallet.WalletAddress = walletChecksumAddress

					//Set the wallet name
					_userWallet.WalletName = walletConfig["wallet_name"].(string)

					//Add the wallet to the user config
					_userConfig.Wallets[walletChecksumAddress] = &_userWallet

					if walletConfig["tokens"] != nil {

						for tokenAddress, tokenConfigJSON := range walletConfig["tokens"].(map[string]interface{}) {
							tokenConfigJSON := tokenConfigJSON.(string)
							//Initialize a new token
							_token := Token{}

							//Initialize a new tokenConfiguration
							_tokenConfig := JSONToTokenConfig(tokenConfigJSON)

							//Get the chain name from the token config
							_chainName := _tokenConfig.ChainName

							//Get the token checksum address
							tokenChecksumAddress, err := dexbotUtils.ToChecksumAddress(tokenAddress)
							handler.HandleError("Error when trying to checksum token address", err)
							_token.TokenChecksumAddress = tokenChecksumAddress

							//Add the token config to the user wallet
							_token.TokenConfig = &_tokenConfig

							//Get the chain information to initialize the token parameters
							_chain := chain.Chains[_chainName]

							//Get the LP address
							lpAddress := web3Logic.GetLPContractAddress(tokenChecksumAddress, _chain.WrappedNativeTokenAddress, _chain.Factory.FactoryInstance)
							lpChecksumAddress, err := dexbotUtils.ToChecksumAddress(lpAddress)
							handler.HandleError("Error when checksumming lp address", err)
							_token.LPChecksumAddress = lpChecksumAddress

							//Initialize a contract instance for the LP
							_token.LPInstance = web3Logic.GetContractInstance(lpAddress, dexbotABI.LPABI, &_chain.RPCNode)

							//Add the token to the user wallet
							_userWallet.Tokens[tokenChecksumAddress] = &_token

						}

					}
				}
			}
		}
	}

	//Set the UserConfig variable
	UserConfig = &_userConfig

}

//Toggle error reporting on/off
func ToggleErrorReporting(toggled bool) {
	//Toggle error reporting on/off locally
	handler.ToggleErrorReporting(toggled)

	//Toggle error reporting on/off in the database
	go database.ToggleErrorReporting(toggled)
}
