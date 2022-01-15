package database

import (
	"dexbot/internal/authentication"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/web3Fork"
	"fmt"
	"time"

	"github.com/umbracle/go-web3"
)

//Check if the user has enabled error reporting or not
func GetErrorReportingSetting() bool {
	//Get the user config
	userConfig := Get(fmt.Sprintf("user_configs/%s", authentication.FirebaseAuthToken.LocalId))

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
	Update(fmt.Sprintf("user_configs/%s", authentication.FirebaseAuthToken.LocalId), errorReportingSettingPayload)
}

//Get the user wallet/token configurations from the database
func GetUserConfig() map[string]interface{} {

	//Get the user configurations settings
	return Get(fmt.Sprintf("user_configs/%s/wallets", authentication.FirebaseAuthToken.LocalId))
}

//Update the user wallet/token configurations
func UpdateUserWalletsConfig(userConfig map[string]interface{}) {
	//Create a data payload to send to the database
	updatedConfigPayload := make(map[string]interface{})

	//Set the value as the updated userConfig
	updatedConfigPayload["wallets"] = userConfig

	//Send the payload to the database
	Update(fmt.Sprintf("user_configs/%s", authentication.FirebaseAuthToken.LocalId), updatedConfigPayload)
}

//Send an alive signal to the database with the wallet/token configurations to be added to the token queues
func SendPulse(data *map[string]interface{}) {

	//Create a data payload to be sent to the database
	pulsePayload := make(map[string]interface{})

	//Get the version checksum so DEXbot can validate the integrity of the application
	versionChecksum := dexbotUtils.RunChecksumOnDEXbotBinary()

	//Create a map to store the pulse data
	pulseData := make(map[string]interface{})

	//Add a time connected timestamp
	pulseData["time_connected"] = time.Now().Unix()

	//Add the current userConfig data to the payload
	pulseData["wallets"] = *data

	//Put the pulse data into the pulse payload
	pulsePayload[versionChecksum] = pulseData

	//As an infinite loop
	for {

		//Add the current timestamp to the payload
		timestamp := time.Now().Unix()
		pulseData["alive_timestamp"] = timestamp

		//Send the pulse payload to the database
		go Update(fmt.Sprintf("active_clients/%s", authentication.FirebaseAuthToken.LocalId), pulsePayload)

		//Sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}
}

//Send a signed transaction to the database
func SendSignedTransactionToDexbot(signedTransaction *web3.Transaction, chainName string, tokenAddress string, walletAddress string, randomSeedKey string) {
	//Create a data payload to be sent to the database
	signedTransactionPayload := make(map[string]interface{})

	//Marshal the signed transaction to a JSON format
	signedTransactionJSON, err := web3Fork.MarshalJSON(signedTransaction)
	handler.HandleError("database, SendSignedTransactionToDexbot, signedTransaction.MarshalJSON", err)

	//Add the signed transaction to the payload
	signedTransactionPayload[randomSeedKey] = signedTransactionJSON

	//Send the payload to the database
	go Update(fmt.Sprintf("signed_swap_tx/%s/%s/%s/%s", chainName, authentication.FirebaseAuthToken.LocalId, tokenAddress, walletAddress), signedTransactionPayload)
}
