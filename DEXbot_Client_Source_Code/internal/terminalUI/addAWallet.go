package terminalUI

import (
	"dexbot/internal/database"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

//Display UI to add a new wallet
func addANewWallet() {
	//Enter the wallet name
	walletName := enterWalletDetails("Wallet Name: ", "Enter the wallet name.", "", "")

	//Enter the wallet address
	checksumWalletAddress := enterWalletAddress("Wallet Address: ", "Enter the wallet address.", walletName, "")

	//Prompt the user to validate the new wallet details
	validateNewWalletInputs(walletName, checksumWalletAddress)
}

//Prompt the user to enter the wallet details for the specified wallet setting passed into the function
func enterWalletDetails(inputMessage string, description string, walletName string, walletAddress string) string {
	clearTerminal()
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint(description)),
		terminalPrinter.Sprintf("\nWallet name: %s\n\nWallet address: %s\n", walletName, walletAddress),
	)
	//Wait for the user input
	userInput := handler.Input(inputMessage)

	//Return the user input
	return userInput
}

func enterWalletAddress(inputMessage string, description string, walletName string, walletAddress string) string {
	clearTerminal()
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint(description)),
		terminalPrinter.Sprintf("\nWallet name: %s\n\nWallet address: %s\n", walletName, walletAddress),
	)
	//Wait for the user input
	userInput := handler.Input(inputMessage)

	//Checksum the input wallet address
	checksumWalletAddress, err := dexbotUtils.ToChecksumAddress(userInput)

	//If the address is not a valid ethereum address, prompt the user to try again
	if err != nil {
		fmt.Println(err)
		checksumWalletAddress = enterWalletAddress(inputMessage, "Address is not valid, please check your input and re-enter the wallet address.", walletName, walletAddress)
	}

	//return the checksummed wallet address
	return checksumWalletAddress

}

//Prompt the user to review and confirm the wallet settings that they input
func validateNewWalletInputs(walletName string, walletAddress string) {

	//Update the display to show the user input token settings and prompt the user to confirm the settings
	clearTerminal()
	validation := enterWalletDetails("y/n: ", "Please confirm your wallet information. y/n", walletName, walletAddress)
	validation = strings.ToLower(validation)

	//If the user denies the settings
	if validation == "n" {
		//Return to the main menu
		mainMenu()
	} else if validation == "y" {
		//^^ If the user confirms the settings

		//Add the wallet to the user configurations
		userConfig.AddWalletToUserConfig(walletName, walletAddress)

		//Add the wallet the the remote user config
		database.AddWalletToUserConfig(walletAddress, walletName)

		//Display the configureWallet UI
		configureWallet(walletName, walletAddress)

	} else {
		//^^ If the user input is not "y" or "n", prompt the user for input again
		validateNewWalletInputs(walletName, walletAddress)
	}
}
