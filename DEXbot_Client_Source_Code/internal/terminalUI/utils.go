package terminalUI

import (
	"dexbot/internal/database"
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

//Human readable read names for configuration settings
var hrConfigurationNames = make(map[string]interface{})

func initializeChainNativeTokens() {
	chainNativeTokens["Ethereum"] = "ETH"
	chainNativeTokens["BSC"] = "BNB"
	chainNativeTokens["Optimism"] = "ETH"
	chainNativeTokens["Arbitrum"] = "ETH"
	chainNativeTokens["Polygon"] = "MATIC"
}

func initializeHumanReadableDescriptions() {

	humanReadableDescriptions := database.GetHumanReadableDescriptions()

	tokenSettingDescriptions["token_name"] = humanReadableDescriptions["token_name_description"].(string)
	tokenSettingDescriptions["contract_address"] = humanReadableDescriptions["contract_address_description"].(string)
	tokenSettingDescriptions["minimum_token_balance"] = humanReadableDescriptions["minimum_token_balance_description"].(string)
	tokenSettingDescriptions["minimum_24_hour_price_change"] = humanReadableDescriptions["minimum_24_hour_price_change_description"].(string)
	tokenSettingDescriptions["minimum_usd_sell_price"] = humanReadableDescriptions["minimum_usd_sell_price_description"].(string)
	tokenSettingDescriptions["minimum_usd_sell_amount"] = humanReadableDescriptions["minimum_usd_sell_amount_description"].(string)
	tokenSettingDescriptions["maximum_usd_sell_amount"] = humanReadableDescriptions["maximum_usd_sell_amount_description"].(string)
	tokenSettingDescriptions["minimum_time_in_seconds_between_sells"] = humanReadableDescriptions["minimum_time_in_seconds_between_sells_description"].(string)

	hrConfigurationNames["token_name"] = humanReadableDescriptions["token_name_title"].(string)
	hrConfigurationNames["contract_address"] = humanReadableDescriptions["contract_address_title"].(string)
	hrConfigurationNames["minimum_token_balance"] = humanReadableDescriptions["minimum_token_balance_title"].(string)
	hrConfigurationNames["minimum_24_hour_price_change"] = humanReadableDescriptions["minimum_24_hour_price_change_title"].(string)
	hrConfigurationNames["minimum_usd_sell_price"] = humanReadableDescriptions["minimum_usd_sell_price_title"].(string)
	hrConfigurationNames["minimum_usd_sell_amount"] = humanReadableDescriptions["minimum_usd_sell_amount_title"].(string)
	hrConfigurationNames["maximum_usd_sell_amount"] = humanReadableDescriptions["maximum_usd_sell_amount_title"].(string)
	hrConfigurationNames["minimum_time_in_seconds_between_sells"] = humanReadableDescriptions["minimum_time_in_seconds_between_sells_title"].(string)

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
