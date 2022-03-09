package terminalUI

import (
	"dexbot/internal/database"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/shopspring/decimal"
)

//Display UI to select a token to configure
func selectATokenToEdit(walletName string, walletAddress string) {
	clearTerminal()

	//Display all of the tokens in the wallet as options
	i := 1
	options := []string{}
	optionsMap := make(map[string]map[string]interface{})
	panels := pterm.Panels{}
	tokens := userConfig.UserConfig.Wallets[walletAddress].Tokens
	for _, token := range tokens {
		tokenConfig := token.TokenConfig
		//Initialize new sub panel list for UI display
		subPanelList := []pterm.Panel{}
		//Create new panel
		panel := pterm.Panel{}
		//Add data to the panel
		panel.Data = terminalCenterPrinter.Sprintf("%v.) %s\n", i, fmt.Sprintf("%s: %s", tokenConfig.ChainName, tokenConfig.TokenName))
		//Add panel to sub panel list
		subPanelList = append(subPanelList, panel)
		//Add sub panel list to panel list
		panels = append(panels, subPanelList)

		//Add each option to display options and options map
		options = append(options, fmt.Sprint(i))
		optionData := make(map[string]interface{})
		optionData["tokenName"] = tokenConfig.TokenName
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
		_tokenConfig := optionsMap[c]["tokenConfig"].(*userConfig.TokenConfiguration)
		editToken(walletName, walletAddress, optionsMap[c]["tokenName"].(string), _tokenConfig)
	}
}

//Display UI to edit token settings, or delete the token
func editToken(walletName string, walletAddress string, tokenName string, tokenConfig *userConfig.TokenConfiguration) {
	clearTerminal()
	//Display options in the terminal UI
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Edit %s (%s)", tokenName, tokenConfig.ContractAddress)),
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
func confirmRemoveToken(walletName string, walletChecksumAddress string, tokenName string, tokenConfig *userConfig.TokenConfiguration) {

	//Display the confirmation prompt
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Are you sure you want to delete\n%s (%s)?", tokenName, tokenConfig.ContractAddress)),
	)

	//Wait for the user input
	validation := handler.Input("Delete token? y/n: ")
	validation = strings.ToLower(validation)

	//Wait for the user input
	if validation == "n" {

		//Return to deit token
		editToken(walletName, walletChecksumAddress, tokenName, tokenConfig)

	} else if validation == "y" {
		//^^ If the user confirms the action

		//Delete the token from tokens
		delete(userConfig.UserConfig.Wallets[walletChecksumAddress].Tokens, tokenConfig.ContractAddress)

		//Update the remote user config
		go database.RemoveTokenConfig(walletChecksumAddress, tokenConfig.ContractAddress)

		//Return to select a token to edit UI
		selectATokenToEdit(walletName, walletChecksumAddress)
	} else {
		//^^ If the user input is not "y" or "n", prompt the user for input again
		confirmRemoveToken(walletName, walletChecksumAddress, tokenName, tokenConfig)
	}

}

//Display UI to select a token setting to configure
func configureTokenSettings(walletName string, walletAddress string, tokenName string, tokenConfig *userConfig.TokenConfiguration) {
	clearTerminal()
	//Get all of the current token settings as options
	i := 1
	options := []string{}
	optionsMap := make(map[string]map[string]interface{})
	panels := pterm.Panels{}
	for key, setting := range map[string]interface{}{
		"maximum_usd_sell_amount":               tokenConfig.MaximumUSDSellAmount,
		"minimum_usd_sell_amount":               tokenConfig.MinimumUSDSellAmount,
		"minimum_time_in_seconds_between_sells": tokenConfig.MinimumTimeInSecondsBetweenSells,
		"minimum_24_hour_price_change":          tokenConfig.Minimum24HourPriceChange,
		"minimum_usd_sell_price":                tokenConfig.MinimumUSDSellPrice,
		"minimum_token_balance":                 tokenConfig.MinimumTokenBalance,
	} {
		//Initialize new sub panel list for UI display
		subPanelList := []pterm.Panel{}
		//Create new panel
		panel := pterm.Panel{}
		//Add data to the panel
		panel.Data = terminalCenterPrinter.Sprintf("%v.) %s: %v\n", i, hrConfigurationNames[key], setting)
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
			fmt.Sprintf("Chain Name: %s", tokenConfig.ChainName),
			fmt.Sprintf("Token Name: %s", tokenConfig.TokenName),
			fmt.Sprintf("Contract Address: %s", tokenConfig.ContractAddress),
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
		changeTokenSetting(walletAddress, tokenConfig.ContractAddress, optionsMap[c]["key"].(string), tokenConfig)
		//Return to configure token settings UI
		configureTokenSettings(walletName, walletAddress, tokenName, tokenConfig)
	}
}

//Display UI to change token setting
func changeTokenSetting(walletAddress string, contractAddress string, key string, tokenConfig *userConfig.TokenConfiguration) {

	//Update the UI to prompt the user to change the token setting
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprint(tokenSettingDescriptions[key])),
	)

	//Wait for the user input
	userInput := handler.Input(fmt.Sprintf("Update %s: ", hrConfigurationNames[key]))

	//Add the updated setting to the tokenConfig
	switch key {
	case "maximum_usd_sell_amount":
		_maximumUSDSellAmount, err := decimal.NewFromString(userInput)
		handler.HandleError("Error when trying to convert maximum USD sell amount from input", err)
		tokenConfig.MaximumUSDSellAmount = _maximumUSDSellAmount
	case "minimum_usd_sell_amount":
		_minimumUSDSellAmount, err := decimal.NewFromString(userInput)
		handler.HandleError("Error when trying to convert minimum USD sell amount from string", err)
		tokenConfig.MinimumUSDSellAmount = _minimumUSDSellAmount
	case "minimum_time_in_seconds_between_sells":
		dexbotUtils.StringToInt64(userInput)
	case "minimum_24_hour_price_change":
		_minimum24HourPriceChange, err := decimal.NewFromString(userInput)
		handler.HandleError("Error when trying to convert minimum 24 hour price change", err)
		minimum24HourPriceChangeFloat, _ := _minimum24HourPriceChange.Float64()
		tokenConfig.Minimum24HourPriceChange = minimum24HourPriceChangeFloat
	case "minimum_usd_sell_price":
		_minimumUSDSellPrice, err := decimal.NewFromString(userInput)
		handler.HandleError("Error when trying to convert minimum USD sell price", err)
		minimumUSDSellPriceFloat, _ := _minimumUSDSellPrice.Float64()
		tokenConfig.MinimumUSDSellPrice = minimumUSDSellPriceFloat
	case "minimum_token_balance":
		_minimumTokenBalance, err := decimal.NewFromString(userInput)
		handler.HandleError("Error when trying to convert minimum 24 hour price change", err)
		tokenConfig.MinimumTokenBalance = _minimumTokenBalance

	}

	//Update the remote config
	go userConfig.UpdateRemoteTokenConfig(walletAddress, contractAddress)

}
