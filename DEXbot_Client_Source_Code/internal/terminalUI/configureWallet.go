package terminalUI

import (
	"dexbot/internal/userConfig"
	"fmt"

	"github.com/pterm/pterm"
)

//Display UI to select a wallet to configure
func selectWalletToConfigure() {
	clearTerminal()
	_userConfig := *userConfig.UserConfig

	//Display all of the user wallets as options
	i := 1
	options := []string{}
	optionsMap := make(map[string]map[string]string)
	panels := pterm.Panels{}
	for walletAddress, walletData := range _userConfig {
		walletData := walletData.(map[string]interface{})
		walletName := walletData["wallet_name"].(string)
		//Initialize new sub panel list for UI display
		subPanelList := []pterm.Panel{}
		//Create new panel
		panel := pterm.Panel{}
		//Add data to the panel
		panel.Data = terminalCenterPrinter.Sprintf("%v.) %s (%s)\n", i, walletName, walletAddress)
		//Add panel to sub panel list
		subPanelList = append(subPanelList, panel)
		//Add sub panel list to panel list
		panels = append(panels, subPanelList)

		//Add each option to display options and options map
		options = append(options, fmt.Sprint(i))
		optionData := make(map[string]string)
		optionData["walletName"] = walletName
		optionData["walletAddress"] = walletAddress
		optionsMap[fmt.Sprint(i)] = optionData

		//Increment i
		i += 1
	}

	//Add back as option "0"
	subPanelList := []pterm.Panel{}
	panel := pterm.Panel{}
	panel.Data = terminalCenterPrinter.Sprintf("0.) Back")
	subPanelList = append(subPanelList, panel)
	panels = append(panels, subPanelList)
	options = append(options, "0")

	//Create panel UI display string with all of the options
	panelDisplayString, _ := pterm.DefaultPanel.WithPanels(panels).Srender()

	//Prompt the user to select a wallet to configure
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint("Select a wallet to configure")),
		terminalCenterPrinter.Sprintf("\n"),
		panelDisplayString,
	)

	//Wait for the user input
	c := rawTerminalInput(options)

	//If the user input is "0", go back to the main menu
	if c == "0" {
		mainMenu()

	} else {
		//^^ If the user selects a wallet to configure, display the configure wallet UI
		configureWallet(optionsMap[c]["walletName"], optionsMap[c]["walletAddress"])
	}
}

//Display UI to add a token, edit a token or delete a token from a wallet
func configureWallet(walletName string, walletAddress string) {
	clearTerminal()
	_userConfig := *userConfig.UserConfig
	walletData := _userConfig[walletAddress].(map[string]interface{})

	//Display options in the terminal UI
	if walletData["tokens"] == nil {
		configureWalletWithNoTokens(walletName, walletAddress)
	} else {

		terminalArea.Update(
			terminalCenterPrinter.Sprintf("\n"),
			terminalCenterPrinter.Sprintf("\n"),
			terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("%s (%s)", walletName, walletAddress)),
			terminalCenterPrinter.Sprintf("\n"),
			terminalCenterPrinter.Sprintf("1.) Add a new token\n"),
			terminalCenterPrinter.Sprintf("2.) Edit tokens\n"),
			terminalCenterPrinter.Sprintf("3.) Delete wallet\n"),
			terminalCenterPrinter.Sprintf("0.) Back\n"),
		)

		//Wait for the user input
		c := rawTerminalInput([]string{"1", "2", "3", "0"})

		//Evaluate the input
		switch c {

		//If user input is "1", add a new token
		case "1":
			addANewToken(walletName, walletAddress)

		//If user input is "1", add a new token
		case "2":
			selectATokenToEdit(walletName, walletAddress)

		//If user input is "2", remove the wallet
		case "3":
			confirmRemoveWallet(walletName, walletAddress)

		//If user input is "0", go back to select a wallet to configure
		case "0":
			selectWalletToConfigure()
		}

	}
}

//Display UI to add a token, edit a token or delete a token from a wallet with no tokens
func configureWalletWithNoTokens(walletName string, walletAddress string) {
	clearTerminal()
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("%s (%s)", walletName, walletAddress)),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("1.) Add a new token\n"),
		terminalCenterPrinter.Sprintf("2.) Delete wallet\n"),
		terminalCenterPrinter.Sprintf("0.) Back\n"),
	)

	//Wait for the user input
	c := rawTerminalInput([]string{"1", "2", "0"})

	//Evaluate the input
	switch c {

	//If user input is "1", add a new token
	case "1":
		addANewToken(walletName, walletAddress)

	//If user input is "2", remove the wallet
	case "2":
		confirmRemoveWallet(walletName, walletAddress)

	//If user input is "0", go back to select a wallet to configure
	case "0":
		selectWalletToConfigure()
	}

}
