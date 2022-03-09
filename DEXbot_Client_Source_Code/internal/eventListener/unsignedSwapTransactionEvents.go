package eventListener

import (
	"dexbot/internal/database"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"fmt"
	"strings"

	"github.com/r3labs/sse/v2"
	"github.com/umbracle/go-web3"
)

//Continuously handle incoming payloads from unsigned_swap_tx/{userId}
func handleUnsignedSwapTransactionEvents(eventChannel chan *sse.Event, finished chan bool) {

	for {
		select {

		//If there is a new payload that comes through the eventChannel
		case eventPayload := <-eventChannel:

			//If the event is a put or patch event
			if strings.Contains("put patch", string(eventPayload.Event)) {

				//Parse the reference path from the payload
				referenceTree := getReferenceTree(eventPayload.Data)

				//Unmarshal the eventData from the payload
				eventData := unmarshalEventData(eventPayload.Data)

				//If there is data in the payload
				if eventData["data"] != nil {

					//Initialize the data variable with the eventData
					data := eventData["data"].(map[string]interface{})

					//If the data payload is an error message from dexbot
					if data["error"] != nil {

						//Handle the error message from dexbot
						handleServerSideMessage(data["error"].(string))

					} else {

						//If the payload reference path is not empty, meaning that data was added to the reference during runtime
						if len(referenceTree) > 0 {
							//Get the chain name for the transaction
							chainName := referenceTree[0]

							//Get the token checksum address
							tokenChecksumAddress, err := dexbotUtils.ToChecksumAddress(referenceTree[1])
							handler.HandleError("Error when converting token address to checksum from unsigned tx payload", err)

							//Get the wallet checksum address
							walletChecksumAddress, err := dexbotUtils.ToChecksumAddress(referenceTree[2])
							handler.HandleError("Error when converting wallet address to checksum from unsigned tx payload", err)

							//Get the target wallet from the user wallets
							userWallet := userConfig.UserConfig.Wallets[walletChecksumAddress]

							//For each unsigned transaction in the payload
							for randomSeedKey, transactionData := range data {

								//Unpack the transaction data into a web3.Transaction
								unsignedTransaction := dexbotUtils.UnpackTransactionPayload(transactionData.(string))

								//Get the method signature of the transaction
								methodSignature := dexbotUtils.GetMethodSignature(unsignedTransaction.Input)

								//Handle the transaction depending on the method signature
								handleMethodSignature(unsignedTransaction, methodSignature, chainName, tokenChecksumAddress, randomSeedKey, *userWallet)

							}

						} else {

							//^^ Else if the reference path is empty, meaning that data was added to the reference on initialization
							//For each blockchain that the user has unsigned transactions for
							for chainName, unsignedTranscations := range data {
								unsignedTranscations := unsignedTranscations.(map[string]interface{})

								//For each token that the user has an unsigned transaction for
								for tokenAddress, walletData := range unsignedTranscations {
									walletData := walletData.(map[string]interface{})

									//Convert token address to checksum
									tokenChecksumAddress, err := dexbotUtils.ToChecksumAddress(tokenAddress)
									handler.HandleError("Error when converting token address to checksum from unsigned tx payload", err)

									//For each wallet that the user has an unsigned transaction for
									for walletAddress, transactionData := range walletData {
										transactionDataMap := transactionData.(map[string]interface{})

										//Convert wallet address to checksum
										walletChecksumAddress, err := dexbotUtils.ToChecksumAddress(walletAddress)
										handler.HandleError("Error when converting wallet address to checksum from unsigned tx payload", err)

										//Get the user wallet
										userWallet := userConfig.UserConfig.Wallets[walletChecksumAddress]

										//For each unsigned transaction
										for randomSeedKey, transactionData := range transactionDataMap {

											//Unpack the transaction data into a web3.Transaction
											unsignedTransaction := dexbotUtils.UnpackTransactionPayload(transactionData.(string))

											//Get the method signature of the transaction
											methodSignature := dexbotUtils.GetMethodSignature(unsignedTransaction.Input)

											//Handle the transaction depending on the method signature
											handleMethodSignature(unsignedTransaction, methodSignature, chainName, tokenChecksumAddress, randomSeedKey, *userWallet)
										}
									}
								}
							}
						}
					}
				}
			}

		//If there is a payload that notifies the listener is finished, stop the handler
		case <-finished:
			return
		}
	}

}

func handleMethodSignature(unsignedTransaction *web3.Transaction, methodSignature string, chainName string, tokenChecksumAddress string, randomSeedKey string, userWallet userConfig.UserWallet) {

	//if the transaction is swapAndTransferUnwrappedNatoWithV2
	if methodSignature == "29dccf6d" {
		handleSwapTransaction(unsignedTransaction, "swapAndTransferUnwrappedNatoWithV2", chainName, tokenChecksumAddress, randomSeedKey, userWallet)

		//Tell DEXbot to add the wallet back into the queue
		database.SendAddToQueueCommand(chainName, tokenChecksumAddress, userWallet.WalletAddress)

	} else if methodSignature == "d51a721f" {
		//^^ else if the transaction is swapAndTransferUnwrappedNatoSupportingFeeOnTransferTokensWithV2
		handleSwapTransaction(unsignedTransaction, "swapAndTransferUnwrappedNatoSupportingFeeOnTransferTokensWithV2", chainName, tokenChecksumAddress, randomSeedKey, userWallet)

		//Tell DEXbot to add the wallet back into the queue
		database.SendAddToQueueCommand(chainName, tokenChecksumAddress, userWallet.WalletAddress)

	} else if methodSignature == "095ea7b3" {
		//^^ else if the transaction is an approve tx
		handleApproveTransaction(unsignedTransaction, chainName, tokenChecksumAddress, randomSeedKey, userWallet)

	} else {
		//^^ else if the method signature is not recognized
		handler.HandleError("Unrecognized method signature when processing unsigned transaction", fmt.Errorf(fmt.Sprintf("MethodSig:%s", methodSignature)))
	}
}
