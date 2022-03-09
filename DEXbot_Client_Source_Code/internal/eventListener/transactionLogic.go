package eventListener

import (
	"dexbot/internal/chain"
	"dexbot/internal/database"
	"dexbot/internal/dexbotABI"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"dexbot/internal/web3Logic"
	"time"

	"github.com/umbracle/go-web3"
)

func handleSwapTransaction(swapTransaction *web3.Transaction, methodName string, chainName string, tokenChecksumAddress string, randomSeedKey string, userWallet userConfig.UserWallet) {
	_token := userWallet.Tokens[tokenChecksumAddress]
	_tokenConfig := _token.TokenConfig
	_chain := chain.Chains[chainName]

	//Get the swapTransaction "to" address
	transactionToAddress := swapTransaction.To
	transactionToChecksumAddress, err := dexbotUtils.ToChecksumAddress(transactionToAddress.String())
	handler.HandleError("Error when trying to checksum transactionToAddress", err)

	//Decode the transaction input data
	abacusSwapTx := dexbotABI.DecodeAbacusSwapInputData(methodName, swapTransaction.Input)
	handler.HandleError("Error when trying to checksum swapToAddress", err)

	//Check that the transaction "to" address is the Abacus address
	if transactionToChecksumAddress == _chain.AbacusAddress {

		// Check that the tokenIn is a token in the user wallet
		if abacusSwapTx.TokenIn == userWallet.Tokens[abacusSwapTx.TokenIn].TokenChecksumAddress {

			//Check that minimum time between sells has passed
			if time.Now().Unix()-_tokenConfig.LastSellTimestamp > _tokenConfig.MinimumTimeInSecondsBetweenSells {

				//Check that the LP passed into the abacus swap tx is the token LP
				if abacusSwapTx.LP == _token.LPChecksumAddress {

					//Sign the transaction
					dexbotUtils.BluePrinter.Println("Transaction recieved, signing transaction")

					signedTransaction := web3Logic.SignTransaction(chain.Chains[chainName].ChainID, swapTransaction, userWallet.PrivateKey)

					//Update the last sell timestamp for the token
					userConfig.UpdateLastSellTimestamp(userWallet.WalletAddress, tokenChecksumAddress)

					//Send the signed transaction back to DEXbot
					database.SendSignedTransactionToDexbot(signedTransaction, chainName, tokenChecksumAddress, userWallet.WalletAddress, randomSeedKey)
					dexbotUtils.GreenPrinter.Println("Sent signed transaction to DEXbot")

					//Add the wallet back into the queue
					database.SendAddToQueueCommand(chainName, tokenChecksumAddress, userWallet.WalletAddress)
				} else {
					dexbotUtils.YellowPrinter.Println("Debug ---- incoming transaction skipped, code:3")
				}
			} else {
				dexbotUtils.YellowPrinter.Println("Debug ---- incoming transaction skipped, code:2")
			}
		} else {
			dexbotUtils.YellowPrinter.Println("Debug ---- incoming transaction skipped, code:1")
		}
	} else {
		dexbotUtils.YellowPrinter.Println("Debug ---- incoming transaction skipped, code:0")
	}
}

func handleApproveTransaction(approveTransaction *web3.Transaction, chainName string, tokenChecksumAddress string, randomSeedKey string, userWallet userConfig.UserWallet) {

	//Decode the transaction input data
	transactionData := dexbotABI.DecodeInputDataWithERC20ABI("approve", approveTransaction.Input)
	spenderAddress := transactionData["_spender"].(web3.Address)

	//Check that the approve transaction is approving the Abacus
	if spenderAddress.String() == chain.Chains[chainName].AbacusAddress {
		//Sign the transaction
		dexbotUtils.BluePrinter.Println("Transaction recieved, signing transaction")
		signedTransaction := web3Logic.SignTransaction(chain.Chains[chainName].ChainID, approveTransaction, userWallet.PrivateKey)

		//Send the signed transaction back to DEXbot
		database.SendSignedTransactionToDexbot(signedTransaction, chainName, tokenChecksumAddress, userWallet.WalletAddress, randomSeedKey)
		dexbotUtils.GreenPrinter.Println("Sent signed transaction to DEXbot")

		//Add the wallet back into the queue
		database.SendAddToQueueCommand(chainName, tokenChecksumAddress, userWallet.WalletAddress)
	}

}
