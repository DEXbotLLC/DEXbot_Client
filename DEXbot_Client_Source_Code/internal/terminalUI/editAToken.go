package terminalUI

import (
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

//Display UI to select a token to configure
func selectATokenToEdit(walletName string, walletAddress string) {
	clearTerminal()
	_userConfig := *userConfig.UserConfig

	//Display all of the tokens in the wallet as options
	i := 1
	options := []string{}
	optionsMap := make(map[string]map[string]interface{})
	panels := pterm.Panels{}
	walletData := _userConfig[walletAddress].(map[string]interface{})
	tokens := walletData["tokens"].(map[string]interface{})
	for _, tokenConfig := range tokens {
		tokenConfig := tokenConfig.(map[string]interface{})
		//Initialize new sub panel list for UI display
		subPanelList := []pterm.Panel{}
		//Create new panel
		panel := pterm.Panel{}
		//Add data to the panel
		panel.Data = terminalCenterPrinter.Sprintf("%v.) %s\n", i, fmt.Sprintf("%s: %s", tokenConfig["chain_name"], tokenConfig["token_name"]))
		//Add panel to sub panel list
		subPanelList = append(subPanelList, panel)
		//Add sub panel list to panel list
		panels = append(panels, subPanelList)

		//Add each option to display options and options map
		options = append(options, fmt.Sprint(i))
		optionData := make(map[string]interface{})
		optionData["tokenName"] = tokenConfig["token_name"]
		optionData["tokenConfig"] = tokenConfig
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

	//Prompt the user to select a token to configure
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprint("Select a token to configure")),
		terminalCenterPrinter.Sprintf("\n"),
		panelDisplayString,
	)

	//Wait for the user input
	c := rawTerminalInput(options)

	//If the user input is "0", go back to configure wallet UI
	if c == "0" {
		configureWallet(walletName, walletAddress)

	} else {
		//^^ If the user selects a token to configure, display the edit token UI
		editToken(walletName, walletAddress, optionsMap[c]["tokenName"].(string), optionsMap[c]["tokenConfig"].(map[string]interface{}))
	}
}

//Display UI to edit token settings, or delete the token
func editToken(walletName string, walletAddress string, tokenName string, tokenConfig map[string]interface{}) {
	clearTerminal()
	//Display options in the terminal UI
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Edit %s (%s)", tokenName, tokenConfig["contract_address"].(string))),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("1.) Configure token settings \n"),
		terminalCenterPrinter.Sprintf("2.) Delete token\n"),
		terminalCenterPrinter.Sprintf("0.) Back\n"),
	)

	//Wait for the user input
	c := rawTerminalInput([]string{"1", "2", "0"})

	//Evaluate the input
	switch c {

	//If user input is "1", configure the token settings
	case "1":
		configureTokenSettings(walletName, walletAddress, tokenName, tokenConfig)

	//If user input is "2", remove the token
	case "2":
		confirmRemoveToken(walletName, walletAddress, tokenName, tokenConfig)

	//If user input is "0", go back to select a token to configure
	case "0":
		selectATokenToEdit(walletName, walletAddress)
	}
}

//Display UI to confirm action when removing token
func confirmRemoveToken(walletName string, walletAddress string, tokenName string, tokenConfig map[string]interface{}) {
	contractAddress := tokenConfig["contract_address"].(string)

	//Display the confirmation prompt
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Are you sure you want to delete\n%s (%s)?", tokenName, contractAddress)),
	)

	//Wait for the user input
	validation := handler.Input("Delete token? y/n: ")
	validation = strings.ToLower(validation)

	//Wait for the user input
	if validation == "n" {

		//Return to deit token
		editToken(walletName, walletAddress, tokenName, tokenConfig)

	} else if validation == "y" {
		//^^ If the user confirms the action

		//Get the current user config
		_userConfig := *userConfig.UserConfig

		//Get the wallet data for the wallet address
		walletData := _userConfig[walletAddress].(map[string]interface{})

		//Get the wallet tokens
		tokens := walletData["tokens"].(map[string]interface{})

		//Delete the token from tokens
		delete(tokens, contractAddress)

		//Update the remote user config
		go userConfig.UpdateUserWalletsConfig()

		//Return to select a token to edit UI
		selectATokenToEdit(walletName, walletAddress)
	} else {
		//^^ If the user input is not "y" or "n", prompt the user for input again
		confirmRemoveToken(walletName, walletAddress, tokenName, tokenConfig)
	}

}

//Display UI to select a token setting to configure
func configureTokenSettings(walletName string, walletAddress string, tokenName string, tokenConfig map[string]interface{}) {
	clearTerminal()
	//Get all of the current token settings as options
	i := 1
	options := []string{}
	optionsMap := make(map[string]map[string]interface{})
	panels := pterm.Panels{}
	for key, setting := range tokenConfig {
		if !strings.Contains("chain_name token_name contract_address, last_sell_timestamp", key) {
			//Initialize new sub panel list for UI display
			subPanelList := []pterm.Panel{}
			//Create new panel
			panel := pterm.Panel{}
			//Add data to the panel
			panel.Data = terminalCenterPrinter.Sprintf("%v.) %s: %s\n", i, hrConfigurationNames[key], setting)
			//Add panel to sub panel list
			subPanelList = append(subPanelList, panel)
			//Add sub panel list to panel list
			panels = append(panels, subPanelList)

			//Add each option to display options and options map
			options = append(options, fmt.Sprint(i))
			optionData := make(map[string]interface{})
			optionData["key"] = key
			optionData["setting"] = setting
			optionsMap[fmt.Sprint(i)] = optionData

			//Increment i
			i += 1
		}
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

	//Prompt the user to select a token setting to configure
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultBox.Sprintf(
			"%s\n%s\n%s",
			fmt.Sprintf("Chain Name: %s", tokenConfig["chain_name"]),
			fmt.Sprintf("Token Name: %s", tokenConfig["token_name"]),
			fmt.Sprintf("Contract Address: %s", tokenConfig["contract_address"]),
		)),
		terminalCenterPrinter.Sprintf("\n"),
		panelDisplayString,
	)

	//Wait for the user input
	c := rawTerminalInput(options)

	//If the user input is "0", go back to edit token
	if c == "0" {
		editToken(walletName, walletAddress, tokenName, tokenConfig)

	} else {
		//^^ If the user selects a token setting to configure, display the change token setting UI
		changeTokenSetting(walletAddress, tokenConfig["contract_address"].(string), optionsMap[c]["key"].(string))
		//Return to configure token settings UI
		configureTokenSettings(walletName, walletAddress, tokenName, tokenConfig)
	}
}

//Display UI to change token setting
func changeTokenSetting(walletAddress string, contractAddress string, key string) {
	_userConfig := *userConfig.UserConfig

	//Update the UI to prompt the user to change the token setting
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprint(tokenSettingDescriptions[key])),
	)

	//Wait for the user input
	userInput := handler.Input(fmt.Sprintf("Update %s: ", hrConfigurationNames[key]))

	//Get the wallet data for the wallet address
	walletData := _userConfig[walletAddress].(map[string]interface{})

	//Get the wallet tokens
	tokens := walletData["tokens"].(map[string]interface{})

	//Get the target token
	targetToken := tokens[contractAddress].(map[string]interface{})

	//Update the token setting with the user input
	targetToken[key] = userInput

	//Update the remote config
	go userConfig.UpdateUserWalletsConfig()

}
