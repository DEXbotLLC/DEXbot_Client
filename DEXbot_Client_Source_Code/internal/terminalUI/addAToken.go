package terminalUI

import (
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"fmt"
	"strconv"
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
	//Update the UI to display all of the chain options
	terminalArea.Update(
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf(pterm.DefaultHeader.Sprintf("Select a chain")),
		terminalCenterPrinter.Sprintf("\n"),
		terminalCenterPrinter.Sprintf("1.) Binance Smart Chain\n"),

		terminalCenterPrinter.Sprintf("0.) Back\n"),
	)

	//Wait for the user to select a chain from the options
	c := rawTerminalInput([]string{"1", "0"})
	chains := []string{"BSC"}

	//Parse the user input
	inputAsInt, err := strconv.ParseInt(c, 0, 32)
	handler.HandleError("Issue when parsing int during chain selection", err)

	//Evaluate the input
	switch c {

	//If the user selected a chain option, continue with the UI display to add token details
	default:
		addTokenDetails(walletName, walletAddress, chains[inputAsInt-1])

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
		chainName, "", "", "", "", "", "", "", "", "")

	//Enter the contract address
	contractAddress := enterTokenDetails(
		"Enter the Contract Address: ", tokenSettingDescriptions["contract_address"],
		chainName, tokenName, "", "", "", "", "", "", "", "")

	//Convert address to checksum
	checksumContractAddress, err := dexbotUtils.ToChecksumAddress(contractAddress)
	handler.HandleError("Error when converting token address to checksum", err)

	//Enter the minimum token balance
	minimumTokenBalance := enterTokenDetails(
		"Enter the Minimum Token Balance: ", tokenSettingDescriptions["minimum_token_balance"],
		chainName, tokenName, checksumContractAddress, "", "", "", "", "", "", "")

	//Enter the minimum 24 hour price change
	minimum24HourPriceChange := enterTokenDetails(
		"Enter the Minimum 24 Hour Price Change: ", tokenSettingDescriptions["minimum_24_hour_price_change"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, "", "", "", "", "", "")

	//Enter the minimum USD sell price
	minimumUSDSellPrice := enterTokenDetails(
		"Enter the Minimum USD sell price: ", tokenSettingDescriptions["minimum_usd_sell_price"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, "", "", "", "", "")

	//Enter the minimum nato sell amount
	minimumNatoSellAmount := enterTokenDetails(fmt.Sprintf("Enter the Minimum %s Sell Amount: ", chainNativeTokens[chainName]), tokenSettingDescriptions["minimum_nato_sell_amount"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, "", "", "", "")

	//Enter the maximum nato sell amount
	maximumNatoSellAmount := enterTokenDetails(fmt.Sprintf("Enter the Maximum %s Sell Amount: ", chainNativeTokens[chainName]), tokenSettingDescriptions["maximum_nato_sell_amount"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, "", "", "")

	//Enter the maximum sell percent of buy transaction
	maximumSellPercentOfBuyTransaction := enterTokenDetails("Enter the Maximum Sell Percent Of Buy Transaction: ", tokenSettingDescriptions["maximum_sell_percent_of_buy_transaction"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, maximumNatoSellAmount, "", "")

	//Enter the minimum time between sells
	minimumTimeBetweenSells := enterTokenDetails("Enter the Minimum Time Between Sells (seconds): ", tokenSettingDescriptions["minimum_time_in_seconds_between_sells"],
		chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, maximumNatoSellAmount, maximumSellPercentOfBuyTransaction, "")

	//Prompt the user to validate the inputs
	validateNewTokenInputs(walletName, walletAddress, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, maximumNatoSellAmount, maximumSellPercentOfBuyTransaction, minimumTimeBetweenSells)
}

//Prompt the user to enter the token details for the specified token setting passed into the function
func enterTokenDetails(inputMessage string, description string, chainName string, tokenName string, checksumContractAddress string, minimumTokenBalance string,
	minimum24HourPriceChange string, minimumUSDSellPrice string, minimumNatoSellAmount string, maximumNatoSellAmount string,
	maximumSellPercentOfBuyTransaction string, minimumTimeBetweenSells string) string {

	//Update the display to show the token details
	clearTerminal()
	terminalArea.Update(
		terminalPrinter.Sprintf("\n\n"),
		terminalPrinter.Sprintf(pterm.DefaultHeader.Sprint(description)),
		terminalPrinter.Sprintf("Token Name: %s\n\nContract Address: %s\n\nMinimum Token Balance: %s\n\nMinimum 24 Hour Price Change: %s\n\nMinimum USD Sell Price: %s\n\nMinimum %s Sell Amount: %s\n\nMaximum %s Sell Amount: %s\n\nMaximum Sell Percent of Buy Transaction: %s\n\nMinimum time between sells (in seconds): %s\n",
			tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, chainNativeTokens[chainName], minimumNatoSellAmount, chainNativeTokens[chainName], maximumNatoSellAmount, maximumSellPercentOfBuyTransaction, minimumTimeBetweenSells),
	)

	//Wait for the user input to assign a value to the specified setting
	userInput := handler.Input(inputMessage)

	//Return the user input for the specified setting
	return userInput
}

//Prompt the user to review and confirm the token settings that they input
func validateNewTokenInputs(walletName string, walletAddress string, chainName string, tokenName string, checksumContractAddress string, minimumTokenBalance string, minimum24HourPriceChange string, minimumUSDSellPrice string, minimumNatoSellAmount string, maximumNatoSellAmount string, maximumSellPercentOfBuyTransaction string, minimumTimeBetweenSells string) {

	//Update the display to show the user input token settings and prompt the user to confirm the settings
	clearTerminal()
	validation := enterTokenDetails("y/n: ", "Please confirm the token configuration. y/n", chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, maximumNatoSellAmount, maximumSellPercentOfBuyTransaction, minimumTimeBetweenSells)
	validation = strings.ToLower(validation)

	//If the user does not confirm the settings
	if validation == "n" {
		//Return to configureWallet display
		configureWallet(walletName, walletAddress)

	} else if validation == "y" {
		//^^ If the user confirms the settings

		//Add the token to the user configurations
		addTokenToUserConfig(walletAddress, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, maximumNatoSellAmount, maximumSellPercentOfBuyTransaction, minimumTimeBetweenSells)

		//Display the configureWallet UI
		configureWallet(walletName, walletAddress)
	} else {
		//^^ If the user input is not "y" or "n", prompt the user for input again.
		validateNewTokenInputs(walletName, walletAddress, chainName, tokenName, checksumContractAddress, minimumTokenBalance, minimum24HourPriceChange, minimumUSDSellPrice, minimumNatoSellAmount, maximumNatoSellAmount, maximumSellPercentOfBuyTransaction, minimumTimeBetweenSells)
	}
}

//Add a token to the user configuration
func addTokenToUserConfig(walletAddress string, chainName string, tokenName string, checksumContractAddress string, minimumTokenBalance string, minimum24HourPriceChange string, minimumUSDSellPrice string, minimumNatoSellAmount string, maximumNatoSellAmount string, maximumSellPercentOfBuyTransaction string, minimumTimeBetweenSells string) {
	//Get the current user configuration values
	_userConfig := *userConfig.UserConfig

	//Get the wallet configuration for the walletAddress
	walletData := _userConfig[walletAddress].(map[string]interface{})

	//If there are no tokens set up, create the tokens map
	if walletData["tokens"] == nil {
		walletData["tokens"] = make(map[string]interface{})
	}

	//Get the token data for the wallet
	tokens := walletData["tokens"].(map[string]interface{})

	//Initalize a map for the new token config and add the token settings
	tokenConfig := make(map[string]interface{})
	tokenConfig["chain_name"] = chainName
	tokenConfig["token_name"] = tokenName
	tokenConfig["contract_address"] = checksumContractAddress
	tokenConfig["minimum_token_balance"] = minimumTokenBalance
	tokenConfig["minimum_24_hour_price_change"] = minimum24HourPriceChange
	tokenConfig["minimum_usd_sell_price"] = minimumUSDSellPrice
	tokenConfig["minimum_nato_sell_amount"] = minimumNatoSellAmount
	tokenConfig["maximum_nato_sell_amount"] = maximumNatoSellAmount
	tokenConfig["maximum_sell_percent_of_buy_transaction"] = maximumSellPercentOfBuyTransaction
	tokenConfig["minimum_time_in_seconds_between_sells"] = minimumTimeBetweenSells
	tokenConfig["last_sell_timestamp"] = "0"

	//Add the token config to the wallet tokens
	tokens[checksumContractAddress] = tokenConfig

	//Update the remote user config
	go userConfig.UpdateUserWalletsConfig()

}
