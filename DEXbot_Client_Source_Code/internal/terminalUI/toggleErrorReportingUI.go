package terminalUI

import (
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"

	"github.com/pterm/pterm"
)

//Initialize the error reporting configruation for the user
func initializeErrorReporting() {
	clearTerminal()
	//If the user has not set up a wallet yet
	if len(userConfig.UserConfig.Wallets) == 0 {
		//Prompt the user to toggle error reporting on or off
		toggleErrorReporting()
	}
}

//Toggle error reporting on or off
func toggleErrorReporting() {

	//Prompt the user to allow or deny error reporting to be sent to DEXbot
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint("Allow error messages to be sent to DEXbot?\nYou can change this setting anytime in the 'Help' menu.\ny/n")),
	)

	//Wait for the user input
	userInput := handler.Input("y/n: ")

	//Evaluate the input
	//If the user inputs "y", toggle error reporting on and display the main menu
	if userInput == "y" {
		userConfig.ToggleErrorReporting(true)
		mainMenu()

	} else if userInput == "n" {
		//^^If the user inputs "n", toggle error reporting off and display the main menu
		userConfig.ToggleErrorReporting(false)
		mainMenu()
	} else {
		//^^If the user does not input "y" or "n", prompt the user again
		toggleErrorReporting()
	}
}
