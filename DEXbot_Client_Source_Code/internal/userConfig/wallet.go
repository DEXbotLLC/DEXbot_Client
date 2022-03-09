package userConfig

import (
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/web3Logic"
	"fmt"
)

//Initialize a new private key and wipe the input after usage
func InitializePrivateKeyWithByteSliceZeroing(walletChecksumAddress string) *web3Logic.Key {
	//Initialize a key variable
	var walletKey *web3Logic.Key

	//Get private key input in isolated scope to destroy the input through garbage collector
	{
		//Initialize private key
		privateKeyByteSlice := handler.InputPrivateKey(fmt.Sprintf("Enter the private key for %s: ", walletChecksumAddress))
		ecdsaPrivateKey, err := dexbotUtils.BytesToECDSA(privateKeyByteSlice)
		//Print an empty line in the terminal for readability
		fmt.Println()
		//If there is an error with the private key initialization, stop the program
		if err != nil {
			handler.Exit(fmt.Sprintf("Incorrect or invalid private key for %s. Please check your wallet address/private key and try again.\n", walletChecksumAddress))
		} else {
			//Zero the byte slice containing the private key
			for i := range privateKeyByteSlice {
				privateKeyByteSlice[i] = 0
			}
			//Set user wallet private key
			walletKey = web3Logic.NewKey(ecdsaPrivateKey)
		}
	}
	//Return the wallet key
	return walletKey
}

func AddWalletToUserConfig(walletName string, walletChecksumAddress string) {
	UserConfig.Wallets[walletChecksumAddress] = &UserWallet{
		WalletName:    walletName,
		WalletAddress: walletChecksumAddress,
	}

}

func InitializeWalletKeys() {
	for walletChecksumAddress := range UserConfig.Wallets {
		userWallet := UserConfig.Wallets[walletChecksumAddress]
		userWallet.PrivateKey = InitializePrivateKeyWithByteSliceZeroing(walletChecksumAddress)
	}
}
