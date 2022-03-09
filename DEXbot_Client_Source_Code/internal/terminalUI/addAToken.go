package terminalUI

import (
	"dexbot/internal/chain"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

//Function to trigger series of UI displays to add a new token
func addANewToken(walletName string, walletAddress string) {
	selectAChain(walletName, walletAddress)
}

//UI display to select a chain when adding a new token
func selectAChain(walletName string, walletAddress string) {

	clearTerminal()
	//Get all of the current token settings as options
	i := 1
	options := []string{}
	optionsMap := make(map[string]map[string]interface{})
	panels := pterm.Panels{}
	for _chain := range chain.Chains {
		//Initialize new sub panel list for UI display
		subPanelList := []pterm.Panel{}
		//Create new panel
		panel := pterm.Panel{}
		//Add data to the panel
		panel.Data = terminalCenterPrinter.Sprintf("%v.) %s \n", i, _chain)
		//Add panel to sub panel list
		subPanelList = append(subPanelList, panel)
		//Add sub panel list to panel list
		panels = append(panels, subPanelList)

		//Add each option to display options and options map
		options = append(options, fmt.Sprint(i))
		optionData := make(map[string]interface{})
		optionData["key"] = _chain
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

	//Update the UI to display all of the chain options
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Select a chain")),
		terminalCenterPrinter.Sprintf("\n"),
		panelDisplayString,
	)

	//Wait for the user input
	c := rawTerminalInput(options)

	//Evaluate the input
	switch c {

	//If the user selected a chain option, continue with the UI display to add token details
	default:
		addTokenDetails(walletName, walletAddress, optionsMap[c]["key"].(string))

	//If the user selected "Back", return to the UI to select a wallet to configure
	case "0":
		selectWalletToConfigure()
	}
}

//Add token details when adding a new token
func addTokenDetails(walletName string, walletAddress string, chainName string) {

	//Enter the token name
	tokenName := enterTokenDetails(
		"Enter the Token Name: ", tokenSettingDescriptions["token_name"],
		chainName, "", "", "", "", "", "", "", "")

	//Enter the contract address
	contractAddress := enterTokenDetails(
		"Enter the Contract Address: ", tokenSettingDescriptions["contract_address"],
		chainName, tokenName, "", "", "", "", "", "", "")

	//Convert address to checksum
	checksumContractAddress, err := dexbotUtils.ToChecksumAddress(contractAddress)
	handler.HandleError("Error when converting token address to checksum", err)

	//Enter the minimum token balance
	minimumTokenBalance := enterTokenDetails(
		"Enter the Minimum Token Balance: ", tokenSettingDescriptions["minimum_token_balance"],
		chainName, tokenName, checksumContractAddress, "", "", "", "", "", "")

	//Enter the minimum 24 hour price change
	minimum24HourPriceChange := enterTokenDetails(
		"Enter the Minimum 24 Hour Price Change: ", tokenSettingDescriptions["minimum_24_hour_price_change"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, "", "", "", "", "")

	//Enter the minimum USD sell price
	minimumUSDSellPrice := enterTokenDetails(
		"Enter the Minimum USD sell price: ", tokenSettingDescriptions["minimum_usd_sell_price"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, "", "", "", "")

	//Enter the minimum USD sell amount
	minimumUSDSellAmount := enterTokenDetails("Enter the Minimum USD Sell Amount: ", tokenSettingDescriptions["minimum_usd_sell_amount"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, "", "", "")

	//Enter the maximum USD sell amount
	maximumUSDSellAmount := enterTokenDetails("Enter the Maximum USD Sell Amount: ", tokenSettingDescriptions["maximum_usd_sell_amount"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, "", "")

	//Enter the minimum time between sells
	minimumTimeBetweenSells := enterTokenDetails("Enter the Minimum Time Between Sells (in seconds): ", tokenSettingDescriptions["minimum_time_in_seconds_between_sells"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, "")

	//Prompt the user to validate the inputs
	validateNewTokenInputs(walletName, walletAddress, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, minimumTimeBetweenSells)
}

//Prompt the user to enter the token details for the specified token setting passed into the function
func enterTokenDetails(inputMessage string, description string, chainName string, tokenName string, checksumContractAddress string, minimumTokenBalance string,
	minimum24HourPriceChange string, minimumUSDSellPrice string, minimumUSDSellAmount string, maximumUSDSellAmount string, minimumTimeBetweenSells string) string {

	//Update the display to show the token details
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprint(description)),
		terminalPrinter.Sprintf("Token Name: %s\n\nContract Address: %s\n\nMinimum Token Balance: %s\n\nMinimum 24 Hour Price Change: %s\n\nMinimum USD Sell Price: %s\n\nMinimum USD Sell Amount: %s\n\nMaximum USD Sell Amount: %s\n\nMinimum time between sells (in seconds): %s\n",
			tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, minimumTimeBetweenSells),
	)

	//Wait for the user input to assign a value to the specified setting
	userInput := handler.Input(inputMessage)

	//If the user inputs the chain wrapped native token, display an error and reprompt the user
	//Since DEXbot swaps tokens for the wrapped native token, the user can not swap the wrapped native token for the wrapped native token
	if inputMessage == "Enter the Contract Address: " && userInput == chain.Chains[chainName].WrappedNativeTokenAddress {
		errorMessageAndInputMessage := "Since DEXbot swaps tokens for the wrapped native token, you can not input the wrapped native token address. Please try again." + inputMessage
		return enterTokenDetails(errorMessageAndInputMessage, description, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, minimumTimeBetweenSells)
	}

	//Return the user input for the specified setting
	return userInput
}

//Prompt the user to review and confirm the token settings that they input
func validateNewTokenInputs(walletName string, walletAddress string, chainName string, tokenName string, checksumContractAddress string, minimumTokenBalance string, minimum24HourPriceChange string, minimumUSDSellPrice string, minimumUSDSellAmount string, maximumUSDSellAmount string, minimumTimeBetweenSells string) {

	//Update the display to show the user input token settings and prompt the user to confirm the settings
	clearTerminal()
	validation := enterTokenDetails("y/n: ", "Please confirm the token configuration. y/n", chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, minimumTimeBetweenSells)
	validation = strings.ToLower(validation)

	//If the user does not confirm the settings
	if validation == "n" {
		//Return to configureWallet display
		configureWallet(walletName, walletAddress)

	} else if validation == "y" {
		//^^ If the user confirms the settings

		//Add the token to the user configurations
		addTokenToUserConfig(walletAddress, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, minimumTimeBetweenSells)

		//Display the configureWallet UI
		configureWallet(walletName, walletAddress)
	} else {
		//^^ If the user input is not "y" or "n", prompt the user for input again.
		validateNewTokenInputs(walletName, walletAddress, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumUSDSellAmount, maximumUSDSellAmount, minimumTimeBetweenSells)
	}
}

//Add a token to the user configuration
func addTokenToUserConfig(walletAddress string, chainName string, tokenName string, checksumContractAddress string, minimumTokenBalance string, minimum24HourPriceChange string, minimumUSDSellPrice string, minimumUSDSellAmount string, maximumUSDSellAmount string, minimumTimeInSecondsBetweenSells string) {
	//Get the current user configuration values
	_userWallet := userConfig.UserConfig.Wallets[walletAddress]

	//If there are no tokens set up, create the tokens map
	if len(_userWallet.Tokens) == 0 {
		_userWallet.Tokens = make(map[string]*userConfig.Token)
	}

	_userWallet.Tokens[checksumContractAddress] = &userConfig.Token{
		TokenChecksumAddress: checksumContractAddress,
		TokenConfig: userConfig.CreateTokenConfig(
			chainName,
			checksumContractAddress,
			tokenName,
			minimumUSDSellAmount,
			maximumUSDSellAmount,
			minimum24HourPriceChange,
			minimumTimeInSecondsBetweenSells,
			minimumTokenBalance,
			minimumUSDSellPrice),
	}

	//Update the remote user config
	go userConfig.UpdateRemoteTokenConfig(walletAddress, checksumContractAddress)

}
