package terminalUI

import (
	"dexbot/internal/authentication"
	"dexbot/internal/chain"
	"dexbot/internal/database"
	"dexbot/internal/dexbotABI"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
)

//Initialize the terminalUI Package
func InitalizeTerminalUI() {
	//Get the operating system for OS specific funcationality compatiblilty
	dexbotUtils.InitializeOperatingSystem()
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
	//Initialize the chain options for DEXbot
	initializeChainNativeTokens()

	//Display a message while the client is initializing
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n\n\n"),
		terminalPrinter.Sprintf(dexbotUtils.GreenPrinter.Sprintf("Initializing DEXbot...")),
	)

	//Initialize Human Readable names for token settings and descriptions
	initializeHumanReadableDescriptions()
	//Initialize the ABIs
	dexbotABI.Initialize()
	//Initialize chain options
	chain.Initialize()
	//Initialize the client checksum
	dexbotUtils.Initialize()
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
