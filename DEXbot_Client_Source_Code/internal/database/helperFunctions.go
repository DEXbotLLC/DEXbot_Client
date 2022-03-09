package database

import (
	"dexbot/internal/authentication"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/web3Logic"
	"fmt"
	"time"

	"github.com/umbracle/go-web3"
)

//Check if the user has enabled error reporting or not
func GetErrorReportingSetting() bool {
	//Get the user config
	userConfig := Get(fmt.Sprintf("client_data/%s", authentication.FirebaseAuthToken.LocalId))

	//If error reporting is not set return false
	if userConfig["error_reporting"] == nil {
		return false
	} else {
		//If error reporting is set, return the setting
		return userConfig["error_reporting"].(bool)
	}
}

//Toggle error reporting on and off
func ToggleErrorReporting(toggled bool) {

	//Create a data payload to send to the database
	errorReportingSettingPayload := make(map[string]interface{})

	//Set the value to the toggle setting
	errorReportingSettingPayload["error_reporting"] = toggled

	//Send the payload to the database
	Update(fmt.Sprintf("client_data/%s", authentication.FirebaseAuthToken.LocalId), errorReportingSettingPayload)
}

//Get the user wallet/token configurations from the database
func GetUserConfig() map[string]interface{} {

	//Get the user configurations settings
	return Get(fmt.Sprintf("client_data/%s", authentication.FirebaseAuthToken.LocalId))
}

//Update the user wallet/token configurations
func UpdateTokenConfig(walletChecksumAddress string, tokenChecksumAddress string, tokenConfigJSON []byte) {
	payload := make(map[string]interface{})
	payload[tokenChecksumAddress] = tokenConfigJSON

	//Send the payload to the database
	Update(fmt.Sprintf("client_data/%s/wallets/%s/tokens", authentication.FirebaseAuthToken.LocalId, walletChecksumAddress), payload)
}

func AddWalletToUserConfig(walletAddress string, walletName string) {
	//Initialize a new payload
	payload := make(map[string]interface{})

	//Create a new wallet
	newWallet := make(map[string]interface{})

	//Add the wallet name
	newWallet["wallet_name"] = walletName

	//Add the new wallet to the payload
	payload[walletAddress] = newWallet

	//Send the payload
	Update(fmt.Sprintf("client_data/%s/wallets", authentication.FirebaseAuthToken.LocalId), payload)

}

func RemoveUserWallet(walletChecksumAddress string) {
	Delete(fmt.Sprintf("client_data/%s/wallets/%s", authentication.FirebaseAuthToken.LocalId, walletChecksumAddress))
}

func RemoveTokenConfig(walletChecksumAddress string, tokenChecksumAddress string) {
	Delete(fmt.Sprintf("client_data/%s/wallets/%s/tokens/%s", authentication.FirebaseAuthToken.LocalId, walletChecksumAddress, tokenChecksumAddress))
}

//Get human readable descriptions for the token settings
func GetHumanReadableDescriptions() map[string]interface{} {
	//Get the user configurations settings
	return Get("client_descriptions")
}

//Send an alive signal to the database with the wallet/token configurations to be added to the token queues
func SendPulse() {

	//Create a data payload to be sent to the database
	pulsePayload := make(map[string]interface{})

	//Add a time connected timestamp
	pulsePayload["time_connected"] = time.Now().Unix()

	//As an infinite loop
	for {

		//Add the current timestamp to the payload
		pulsePayload["alive_timestamp"] = time.Now().Unix()

		//Send the pulse payload to the database
		go Update(fmt.Sprintf("client_data/%s", authentication.FirebaseAuthToken.LocalId), pulsePayload)

		//Sleep for 30 seconds
		time.Sleep(30 * time.Second)
	}
}

//Send a signed transaction to the database
func SendSignedTransactionToDexbot(signedTransaction *web3.Transaction, chainName string, tokenAddress string, walletAddress string, randomSeedKey string) {
	//Create a data payload to be sent to the database
	signedTransactionPayload := make(map[string]interface{})

	//Marshal the signed transaction to a JSON format
	signedTransactionJSON, err := web3Logic.MarshalJSON(signedTransaction)
	handler.HandleError("database, SendSignedTransactionToDexbot, signedTransaction.MarshalJSON", err)

	//Add the signed transaction to the payload
	signedTransactionPayload[randomSeedKey] = signedTransactionJSON

	//Send the payload to the database
	go Update(fmt.Sprintf("signed_swap_tx/%s/%s/%s/%s", chainName, authentication.FirebaseAuthToken.LocalId, tokenAddress, walletAddress), signedTransactionPayload)
}

//Send a command to DEXbot to connect the client to the token queues
func SendConnectToDEXbotCommand() {

	//Initialize the payload
	payload := make(map[string]interface{})

	//Add the command to the payload
	payload["ctd"] = "_"

	//Send the payload to active client commands
	go Update(fmt.Sprintf("active_client_commands/%s", authentication.FirebaseAuthToken.LocalId), payload)
}

//Send a command to DEXbot to connect the client to a specific token queue
func SendAddToQueueCommand(chainName string, tokenChecksumAddress string, walletChecksumAddress string) {

	//Initialize the payload
	payload := make(map[string]interface{})

	//Add the command to the payload
	payload[fmt.Sprintf("atq(%s, %s, %s)", chainName, tokenChecksumAddress, walletChecksumAddress)] = "_"

	//Send the payload to active client commands
	go Update(fmt.Sprintf("active_client_commands/%s", authentication.FirebaseAuthToken.LocalId), payload)

}

//Function to send the verison checksum to DEXbot
func SendVersionChecksum() {

	//Initialize the payload
	payload := make(map[string]interface{})

	//Add version checksum to payload
	payload["version_checksum"] = dexbotUtils.DEXbotClientChecksum

	//Send the payload
	go Update(fmt.Sprintf("client_data/%s", authentication.FirebaseAuthToken.LocalId), payload)

}
