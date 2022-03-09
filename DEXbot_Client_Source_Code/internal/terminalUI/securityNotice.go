package terminalUI

import (
	"dexbot/internal/dexbotUtils"
	"os"
)

//UI prompting the user to accept or decline the security notice
func securityNotice() {
	//Display the security message and options
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(dexbotUtils.YellowPrinter.Sprintf("**!!IMPORTANT!!**\n")),
		terminalCenterPrinter.Sprintf(dexbotUtils.WhitePrinter.Sprintf("For security, it is extremely important that you verify the authenticity\nof this application. To verify, quit the app and run the following command:\n")),
		terminalCenterPrinter.Sprintf("%s: %s", dexbotUtils.WhitePrinter.Sprintf("Linux"), dexbotUtils.BluePrinter.Sprintf("sha256sum dexbot-linux")),
		terminalCenterPrinter.Sprintf("%s: %s", dexbotUtils.WhitePrinter.Sprintf("MacOS"), dexbotUtils.BluePrinter.Sprintf("openssl sha256 dexbot-mac")),
		terminalCenterPrinter.Sprintf("%s: %s", dexbotUtils.WhitePrinter.Sprintf("Windows"), dexbotUtils.BluePrinter.Sprintf("certutil -hashfile dexbot-windows.exe SHA256\n")),
		terminalCenterPrinter.Sprintf(dexbotUtils.WhitePrinter.Sprintf("After running the command, compare the checksum to the checksum on\nthe official DEXbot GitHub repo:\n")),
		terminalCenterPrinter.Sprintf(dexbotUtils.BluePrinter.Sprintf("https://github.com/DEXbotLLC/DEXbot_Client\n")),
		terminalCenterPrinter.Sprintf(dexbotUtils.WhitePrinter.Sprintf("You should only proceed if the values match exactly.\n")),
		"------------------------------\n",
		"1.) I have verified the checksum and agree to the terms (https://www.dexbot.io/terms).\n\n",
		"0.) I have not verified the checksum or do not agree with the terms.\n",
	)

	//Wait for the user input
	c := rawTerminalInput([]string{"1", "0"})

	//Evaluate the input
	switch c {
	//If the user input is "1", continue to authentication

	case "1":
	//continue to authentication

	//If the user input is "0", exit the program
	case "0":
		clearTerminal()
		os.Exit(0)
	}

}
