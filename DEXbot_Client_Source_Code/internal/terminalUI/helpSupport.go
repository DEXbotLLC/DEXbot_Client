package terminalUI

import (
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"fmt"
	"os/exec"

	"github.com/pterm/pterm"
)

//Display the help/support UI
func helpSupport() {
	//Update the UI to display the help options
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Official Discord Link:\nhttps://bit.ly/dexbot-discord")),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("1.) Launch Discord\n"),
		terminalCenterPrinter.Sprintf("2.) Toggle Error Reporting\n"),
		terminalCenterPrinter.Sprintf("0.) Back\n"),
	)

	//Wait for the user input
	c := rawTerminalInput([]string{"1", "2", "0"})

	//Evaluate the input
	switch c {

	//If the input is "1", launch the DEXbot Community discord server
	case "1":
		goToDiscord()
		clearTerminal()
		mainMenu()

	//If the input is "2", display the UI to toggle error reporting
	case "2":
		clearTerminal()
		toggleErrorReporting()

	//If user input is "0", go back to the main menu
	case "0":
		mainMenu()
	}
}

//Launch the DEXbot Community discord server
func goToDiscord() {
	var err error

	//DEXbot community discord url
	discordURL := "https://bit.ly/dexbot-discord"

	//Check the operating system and launch discord in the browser
	switch dexbotUtils.OperatingSystem {
	case "linux":
		err = exec.Command("xdg-open", discordURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", discordURL).Start()
	case "darwin":
		err = exec.Command("open", discordURL).Start()
	default:
		fmt.Println("This functionality does not currently support", dexbotUtils.OperatingSystem)
	}
	handler.HandleError("Error when sending user to discord channel", err)
}
