package eventListener

import (
	"dexbot/internal/chainOptions"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userWallets"
	"fmt"
	"strings"

	"github.com/r3labs/sse/v2"
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
						userWallet := userWallets.UserWallets[walletChecksumAddress]

						//For each unsigned transaction in the payload
						for randomSeedKey, transactionData := range data {

							//Unpack the transaction data into a web3.Transaction
							unsignedTransaction := dexbotUtils.UnpackTransactionPayload(transactionData.(string))

							//Get the method signature of the transaction
							methodSignature := dexbotUtils.GetMethodSignature(unsignedTransaction.Input)

							//if the transaction is swapExactTokensForTokens
							if methodSignature == "38ed1739" {
								handleSwapTransaction(unsignedTransaction, "swapExactTokensForTokens", chainName, chainOptions.ChainOptions[chainName], tokenChecksumAddress, randomSeedKey, userWallet)

							} else if methodSignature == "5c11d795" {
								//^^ else if the transaction is swapExactTokensForTokensSupportingFeeOnTransferTokens
								handleSwapTransaction(unsignedTransaction, "swapExactTokensForTokensSupportingFeeOnTransferTokens", chainName, chainOptions.ChainOptions[chainName], tokenChecksumAddress, randomSeedKey, userWallet)

							} else if methodSignature == "095ea7b3" {
								//^^ else if the transaction is an approve tx
								handleApproveTransaction(unsignedTransaction, chainName, chainOptions.ChainOptions[chainName], tokenChecksumAddress, randomSeedKey, userWallet)

							} else {
								//^^ else if the method signature is not recognized
								handler.HandleError("Unrecognized method signature when processing unsigned transaction", fmt.Errorf(fmt.Sprintf("MethodSig:%s", methodSignature)))
							}
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
									userWallet := userWallets.UserWallets[walletChecksumAddress]

									//For each unsigned transaction
									for randomSeedKey, transactionData := range transactionDataMap {

										//Unpack the transaction data into a web3.Transaction
										unsignedTransaction := dexbotUtils.UnpackTransactionPayload(transactionData.(string))

										//Get the method signature of the transaction
										methodSignature := dexbotUtils.GetMethodSignature(unsignedTransaction.Input)

										//if the transaction is swapExactTokensForTokens
										if methodSignature == "38ed1739" {
											handleSwapTransaction(unsignedTransaction, "swapExactTokensForTokens", chainName, chainOptions.ChainOptions[chainName], tokenChecksumAddress, randomSeedKey, userWallet)

										} else if methodSignature == "5c11d795" {

											//^^ else if the transaction is swapExactTokensForTokensSupportingFeeOnTransferTokens
											handleSwapTransaction(unsignedTransaction, "swapExactTokensForTokensSupportingFeeOnTransferTokens", chainName, chainOptions.ChainOptions[chainName], tokenChecksumAddress, randomSeedKey, userWallet)

										} else if methodSignature == "095ea7b3" {
											//^^ else if the transaction is an approve tx
											handleApproveTransaction(unsignedTransaction, chainName, chainOptions.ChainOptions[chainName], tokenChecksumAddress, randomSeedKey, userWallet)

										} else {
											//^^ else if the method signature is not recognized
											handler.HandleError("Unrecognized method signature when processing unsigned transaction", fmt.Errorf(fmt.Sprintf("MethodSig:%s", methodSignature)))
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
