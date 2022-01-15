package eventListener

import (
	"dexbot/internal/chainOptions"
	"dexbot/internal/database"
	"dexbot/internal/dexbotABI"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/userConfig"
	"dexbot/internal/userWallets"
	"strings"

	"github.com/umbracle/go-web3"
)

func handleSwapTransaction(swapTransaction *web3.Transaction, methodName string, chainName string, chainID uint64, tokenChecksumAddress string, randomSeedKey string, userWallet userWallets.UserWallet) {
	//Get the swapTransaction "to" address
	transactionToAddress := swapTransaction.To

	//Decode the transaction input data
	transactionData := dexbotABI.DecodeInputDataWithDEXRouterABI(methodName, swapTransaction.Input)
	swapToAddress := transactionData["to"].(web3.Address)
	swapPath := transactionData["path"].([]web3.Address)

	//Check that the transaction "to" address is the dex router
	if strings.EqualFold(transactionToAddress.String(), chainOptions.DEXRouters[chainName]) {

		//Check that the swap "to" address is to the Dexbot Abacus
		if strings.EqualFold(swapToAddress.String(), chainOptions.AbacusAddresses[chainName]) {

			//Check that the swap path has a length of 2 and that the token out is the wrapped native token
			if len(swapPath) == 2 && strings.EqualFold(swapPath[1].String(), chainOptions.WrappedNativeTokens[chainName]) {

				//Sign the transaction
				dexbotUtils.BluePrinter.Println("Transaction recieved, signing transaction")

				signedTransaction := userWallets.SignTransaction(chainID, swapTransaction, userWallet.PrivateKey)

				//Update the last sell timestamp for the token
				userConfig.UpdateTokenLastSellTimestamp(userWallet.WalletAddress, tokenChecksumAddress)

				//Send the signed transaction back to DEXbot
				database.SendSignedTransactionToDexbot(signedTransaction, chainName, tokenChecksumAddress, userWallet.WalletAddress, randomSeedKey)
				dexbotUtils.GreenPrinter.Println("Sent signed transaction to DEXbot")
			}
		}
	}
}

func handleApproveTransaction(approveTransaction *web3.Transaction, chainName string, chainID uint64, tokenChecksumAddress string, randomSeedKey string, userWallet userWallets.UserWallet) {

	//Decode the transaction input data
	transactionData := dexbotABI.DecodeInputDataWithERC20ABI("approve", approveTransaction.Input)
	spenderAddress := transactionData["_spender"].(web3.Address)

	//Check that the approve transaction is approving the dex router
	if strings.EqualFold(spenderAddress.String(), chainOptions.DEXRouters[chainName]) {
		//Sign the transaction
		dexbotUtils.BluePrinter.Println("Transaction recieved, signing transaction")
		signedTransaction := userWallets.SignTransaction(chainID, approveTransaction, userWallet.PrivateKey)

		//Send the signed transaction back to DEXbot
		database.SendSignedTransactionToDexbot(signedTransaction, chainName, tokenChecksumAddress, userWallet.WalletAddress, randomSeedKey)
		dexbotUtils.GreenPrinter.Println("Sent signed transaction to DEXbot")
	}

}
