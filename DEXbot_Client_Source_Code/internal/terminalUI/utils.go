package terminalUI

import (
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
	"golang.org/x/crypto/ssh/terminal"
)

//Map to hold the native token name for each chain
var chainNativeTokens = make(map[string]interface{})

//Map to hold descriptions for each token setting that will display when configuring token settings
var tokenSettingDescriptions = make(map[string]string)

//Easy to read names for configuration settings
var hrConfigurationNames = make(map[string]interface{})

func initializeChainNativeTokens() {
	chainNativeTokens["Ethereum"] = "ETH"
	chainNativeTokens["BSC"] = "BNB"
	chainNativeTokens["Optimism"] = "ETH"
	chainNativeTokens["Arbitrum"] = "ETH"
	chainNativeTokens["Polygon"] = "MATIC"
}

func initializeTokenSettingDescriptions() {
	tokenSettingDescriptions["token_name"] = "Enter the name of the token.\nGive the token a nick name for easy navigation in the terminal."
	tokenSettingDescriptions["contract_address"] = "Enter the contract address for the token."
	tokenSettingDescriptions["minimum_token_balance"] = "Enter the minimum token balance.\nDEXbot will never sell below this threshold."
	tokenSettingDescriptions["minimum_24_hour_price_change"] = "Enter the minimum 24 hour price change as a percent (ex: 1.00). \nSwaps will never happen when the percent change\nis below this threshold."
	tokenSettingDescriptions["minimum_usd_sell_price"] = "Enter the minimum sell price in USD. DEXbot will never make a sell below this USD value."
	tokenSettingDescriptions["minimum_nato_sell_amount"] = "Enter the minimum Native Token sell amount.\nDEXbot will never make a sell below this threshold."
	tokenSettingDescriptions["maximum_nato_sell_amount"] = "Enter the maximum Native Token sell amount.\nDEXbot will never make a sell above this threshold."
	tokenSettingDescriptions["maximum_sell_percent_of_buy_transaction"] = "Enter the maximum sell percent of buy transaction (ex. 90.00)\nDEXbot creates a sell transaction when tokens are bought.\nWhat percent of the buy transaction do you want to countersell?"
	tokenSettingDescriptions["minimum_time_in_seconds_between_sells"] = "Enter the minimum amount of time that DEXbot will wait\nbetween each sell transaction for the token."
}

func initalizeHRConfigurationNames() {
	hrConfigurationNames["chain_name"] = "Chain Name"
	hrConfigurationNames["token_name"] = "Token Name"
	hrConfigurationNames["contract_address"] = "Contract Address"
	hrConfigurationNames["minimum_token_balance"] = "Minimum Token Balance"
	hrConfigurationNames["minimum_24_hour_price_change"] = "Minimum 24 Hour Price Change"
	hrConfigurationNames["minimum_usd_sell_price"] = "Minimum USD Sell Price"
	hrConfigurationNames["minimum_nato_sell_amount"] = "Minimum Native Token Sell Amount"
	hrConfigurationNames["maximum_nato_sell_amount"] = "Maximum Native Token Sell Amount"
	hrConfigurationNames["maximum_sell_percent_of_buy_transaction"] = "Maximum Sell Percent Of Buy Transaction"
	hrConfigurationNames["minimum_time_in_seconds_between_sells"] = "Minimum Time Between Sells"
}

//Print the DEXbot main menu with big letters
func bigLetters(text string) string {
	str, err := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(text)).Srender()
	handler.HandleError("error when printing big letters", err)
	return str
}

//Get input from the user without requiring the user to press enter
func rawTerminalInput(options []string) string {
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	handler.HandleError("setting stdin to raw:", err)
	for {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)

		char := string(b)
		for _, v := range options {
			if v == char {
				err = terminal.Restore(fd, state)
				handler.HandleError("failed to restore terminal:", err)
				return char
			}
		}
	}

}

//Clear the terminal UI
func clearTerminal() {
	if dexbotUtils.OperatingSystem == "linux" || dexbotUtils.OperatingSystem == "darwin" {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	} else if dexbotUtils.OperatingSystem == "windows" {
		c := exec.Command("cmd", "/c", "cls")
		c.Stdout = os.Stdout
		c.Run()
	}
}
