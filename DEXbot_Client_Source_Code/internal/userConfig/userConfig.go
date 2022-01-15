package userConfig

import (
	"dexbot/internal/database"
	"dexbot/internal/handler"
	"time"
)

//Map to hold user config
var UserConfig *map[string]interface{}

//Send a continous alive signal to the database with the user config data
func InitializePulse() {
	//send an alive signal to the database
	go database.SendPulse(UserConfig)
}

//Initialize the user config from the database
func initializeUserConfig() {

	//Create a map to store the user config
	_userConfig := make(map[string]interface{})

	//Get the user config from the database
	remoteUserConfig := database.GetUserConfig()

	//Unpack the remoteUserConfig into the local _userConfig
	for walletAddress, tokenConfigs := range remoteUserConfig {
		tokenConfigs := tokenConfigs.(map[string]interface{})
		_userConfig[walletAddress] = tokenConfigs
	}

	//Set the UserConfig variable
	UserConfig = &_userConfig
}

//Update the remote user config with the current user config settings
func UpdateUserWalletsConfig() {
	//Update database user config
	database.UpdateUserWalletsConfig(*UserConfig)
}

//Toggle error reporting on/off
func ToggleErrorReporting(toggled bool) {
	//Toggle error reporting on/off locally
	handler.ToggleErrorReporting(toggled)

	//Toggle error reporting on/off in the database
	go database.ToggleErrorReporting(toggled)
}

//Update the last sell timestamp for a given token
func UpdateTokenLastSellTimestamp(walletAddress string, tokenAddress string) {
	//Get the current user config setting
	_userConfig := *UserConfig

	//Get the wallet with the specified wallet address
	wallet := _userConfig[walletAddress].(map[string]interface{})

	//Get the wallet tokens
	tokens := wallet["tokens"].(map[string]interface{})

	//Get the token config
	tokenConfigMap := tokens[tokenAddress].(map[string]interface{})

	//Update the token last sell timestamp
	tokenConfigMap["last_sell_timestamp"] = time.Now().Unix()

	//Update the remote user config
	UpdateUserWalletsConfig()

}
