package terminalUI

import (
	"dexbot/internal/authentication"
	"dexbot/internal/chainOptions"
	"dexbot/internal/database"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
)

//Initialize the terminalUI Package
func InitalizeTerminalUI() {
	//Get the operating system for OS specific funcationality compatiblilty
	dexbotUtils.InitializeOperatingSystem()
	//Initialize the chain options for DEXbot
	initializeChainNativeTokens()
	//Initialize Human Readable names for token settings
	initalizeHRConfigurationNames()
	//Initialize the descriptions for token settings
	initializeTokenSettingDescriptions()
	//Initialize the terminal Display
	initializeTerminalArea()
	initializeTerminalCenterPrinter()
	//Clear the terminal
	clearTerminal()
	//Prompt the user to read the security notice
	securityNotice()
	//Clear the terminal
	clearTerminal()
	//Display the authentication UI and log the user in
	authenticateUser()
	//Initialize the database connection
	database.Initialize()
	//Initialize chainOptions package
	chainOptions.Initialize()
	//Initialize the user configs
	userConfig.Initialize()
	//Initialize error reporting setting
	handler.InitializeErrorReporting(database.GetErrorReportingSetting(), authentication.FirebaseAuthToken.LocalId)
	//Clear the terminal
	clearTerminal()
	//Initialize error reporting setting
	initializeErrorReporting()
	//Display the main menu
	mainMenu()
}
