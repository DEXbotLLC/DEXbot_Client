package terminalUI

import (
	"dexbot/internal/userConfig"
	"os"

	"github.com/pterm/pterm"
)

//Display the main menu
func mainMenu() {
	clearTerminal()
	_userConfig := *userConfig.UserConfig

	//If the user has already set up a wallet, display the main menu options
	if len(_userConfig) > 0 {
		terminalArea.Update(
			terminalPrinter.Sprintf("\n"),
			terminalPrinter.Sprintf("\n"),
			terminalPrinter.Sprintf(bigLetters("DEXbot")),
			terminalPrinter.Sprintf("\n"),
			terminalPrinter.Sprintf("1.) Start DEXbot\n"),
			terminalPrinter.Sprintf("2.) Add a wallet\n"),
			terminalPrinter.Sprintf("3.) Configure a wallet\n"),
			terminalPrinter.Sprintf("4.) Help/Support\n"),
			terminalPrinter.Sprintf("0.) Quit\n"),
		)

		//Wait for the user input
		c := rawTerminalInput([]string{"1", "2", "3", "4", "0"})

		//Evaluate the input
		switch c {

		//If user input is "1", start DEXbot
		case "1":
			startDEXbot()

		//If user input is "2", add a new wallet
		case "2":
			addANewWallet()

		//If user input is "3", select a wallet to configure
		case "3":
			selectWalletToConfigure()

		//If user input is "4", display the help menu
		case "4":
			helpSupport()

		//If user input is "0", exit the program
		case "0":
			confirmExitDEXbot()
		}

	} else {
		//^^ If the user has not set up a wallet yet, display the limited main menu
		newUserMainMenu()
	}
}

//Main menu UI that is displayed when the user has not set up a wallet yet
func newUserMainMenu() {
	clearTerminal()
	//Display the limited main menu options
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(bigLetters("DEXbot")),
		terminalCenterPrinter.Sprintf("Welcome to DEXbot, the sidecar to your crytpocurrency wallet!\nAdd a wallet to get started.\n\n"),
		terminalCenterPrinter.Sprintf("1.) Add a wallet\n"),
		terminalCenterPrinter.Sprintf("2.) Help/Support\n"),
		terminalCenterPrinter.Sprintf("0.) Quit\n"),
	)

	//Wait for the user input
	c := rawTerminalInput([]string{"1", "2", "0"})

	//Evaluate the input
	switch c {

	//If user input is "1", start DEXbot
	case "1":
		addANewWallet()

	//If user input is "2", add a new wallet
	case "2":
		helpSupport()

	//If user input is "0", exit the program
	case "0":
		confirmExitDEXbot()
	}
}

//Display UI to confirm action to exit DEXbot
func confirmExitDEXbot() {
	//Prompt user to confirm action to exit DEXbot
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Exit DEXbot? y/n")),
	)

	//Wait for the user input
	c := rawTerminalInput([]string{"y", "n"})

	//If the user denies the action
	if c == "n" {
		clearTerminal()
		_userConfig := *userConfig.UserConfig

		//If the user has a wallet set up already
		if len(_userConfig) > 0 {
			//Return to main menu
			mainMenu()
		} else {
			//^^ If the user does not have a wallet set up yet

			//Return to new user main menu
			newUserMainMenu()
		}

	} else if c == "y" {
		//Clear the terminal
		clearTerminal()
		//Exit DEXbot
		os.Exit(0)
	}

}
