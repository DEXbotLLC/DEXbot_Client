package terminalUI

import (
	"dexbot/internal/authentication"
	"dexbot/internal/handler"

	"github.com/pterm/pterm"
)

//Get username and password and initialize authentication package to authenticate user
func authenticateUser() {

	//Prompt the user to enter their username
	username := enterUsername()

	//Promt the user to enter their password
	password := enterPassword()

	//Initialize the authentication package with the username and password
	authentication.Initalize(username, password)
}

//Prompt the user for their username
func enterUsername() string {

	//Display the username entry UI
	clearTerminal()
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint("Enter your DEXbot email")),
	)

	//Wait for the user input
	userInput := handler.Input("Username: ")

	//Return the user input
	return userInput
}

//Prompt the user for their password
func enterPassword() string {

	//Display the password entry UI
	clearTerminal()
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint("Enter your DEXbot password")),
	)

	//Wait for the user input
	userInput := handler.InputObscured("Password: ")

	//Return the user input
	return userInput
}
