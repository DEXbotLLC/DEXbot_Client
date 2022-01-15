package userWallets

import (
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"dexbot/internal/web3Fork"
	"fmt"
	"strings"

	"github.com/umbracle/go-web3"
)

//Map to store user wallets
var UserWallets = make(map[string]UserWallet)

//Initialize UserWallet variables with user config data
func initalizeUserWallets() {
	_userConfig := *userConfig.UserConfig
	//for each wallet in the config file, initialize a wallet
	for walletAddress, walletData := range _userConfig {
		walletData := walletData.(map[string]interface{})
		initializeUserWallet(walletAddress, walletData)
	}
}

//Initialize a UserWallet
func initializeUserWallet(walletAddress string, walletData map[string]interface{}) {

	//Create a new UserWallet variable
	userWallet := UserWallet{}

	//Convert the wallet address to checksum
	checksumWalletAddress, err := dexbotUtils.ToChecksumAddress(walletAddress)
	handler.HandleError(fmt.Sprintf("Error when converting %s to checksum address\n", walletAddress), err)

	//Set the UserWallet address to the checksum address
	userWallet.WalletAddress = checksumWalletAddress

	//initialize user wallet private key with byteslice zeroing to wipe the input after usage
	userWallet.PrivateKey = InitializePrivateKeyWithByteSliceZeroing(checksumWalletAddress)

	//If the public address generated from the private key does not match the walletAddress specified, log an error message and exit the program
	if !strings.EqualFold(walletAddress, userWallet.PrivateKey.Address().String()) {
		handler.Exit(fmt.Sprintf("Incorrect or invalid private key for %s. Please check your wallet address/private key and try again.\n", walletAddress))
	} else {
		//^^ If the public address generated from the private key matches the walletAddress specified
		UserWallets[userWallet.WalletAddress] = userWallet
	}

}

//Initialize a new private key and wipe the input after usage
func InitializePrivateKeyWithByteSliceZeroing(walletChecksumAddress string) *web3Fork.Key {
	//Initialize a key variable
	var walletKey *web3Fork.Key

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
			walletKey = web3Fork.NewKey(ecdsaPrivateKey)
		}
	}
	//Return the wallet key
	return walletKey
}

//Sign a transaction with a specified wallet key
func SignTransaction(chainID uint64, txn *web3.Transaction, privateKey *web3Fork.Key) *web3.Transaction {
	//create a signer object with the bsc mainnet chainID
	signer := web3Fork.NewEIP155Signer(chainID)

	//sign the transaction
	signedTx, err := signer.SignTx(txn, privateKey)
	handler.HandleError("Error: SignTransaction, SignTx", err)

	//return the signed transaction
	return signedTx
}
