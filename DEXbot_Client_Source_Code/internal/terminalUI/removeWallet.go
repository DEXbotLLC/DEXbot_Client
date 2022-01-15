package terminalUI

import (
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"strings"

	"github.com/pterm/pterm"
)

//Display UI to confirm action to remove the wallet
func confirmRemoveWallet(walletName string, walletAddress string) {
	//Prompt user to confirm action to remove the wallet
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Are you sure you want to delete\n%s (%s)?", walletName, walletAddress)),
	)

	//Wait for the user input
	validation := handler.Input("Delete wallet? y/n: ")
	validation = strings.ToLower(validation)

	//If the user denies the action
	if validation == "n" {

		//Return to configure wallet UI
		configureWallet(walletName, walletAddress)

	} else if validation == "y" {
		//If the user confirms the action

		//Get the current user config
		_userConfig := *userConfig.UserConfig

		//Delete the wallet from the user config
		delete(_userConfig, walletAddress)

		//Update the remote user config
		go userConfig.UpdateUserWalletsConfig()

		//Return to the main menu
		mainMenu()
	} else {
		//^^ If the user input is not "y" or "n", prompt the user for input again
		confirmRemoveWallet(walletName, walletAddress)
	}

}
