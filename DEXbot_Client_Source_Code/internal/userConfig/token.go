package userConfig

import (
	"dexbot/internal/database"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/handler"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

func TokenConfigToJSON(_tokenConfig TokenConfiguration) []byte {
	tokenConfigJSONBytes, err := json.Marshal(_tokenConfig)
	handler.HandleError("Error when converting token config to bytes", err)
	return tokenConfigJSONBytes
}

func JSONToTokenConfig(tokenConfigJSON string) TokenConfiguration {
	//Initialize a new tokenConfiguration
	_tokenConfig := TokenConfiguration{}

	tokenConfigJSONBytes, err := base64.StdEncoding.DecodeString(tokenConfigJSON)
	handler.HandleError("Error when converting token config json to bytes", err)
	unmarshalErr := json.Unmarshal(tokenConfigJSONBytes, &_tokenConfig)
	handler.HandleError("Error when unmarshing from JSON to TokenConfiguration", unmarshalErr)

	return _tokenConfig
}

func CreateTokenConfig(chainName string, contractAddress string, tokenName string, minimumUSDSellAmount string, maximumUSDSellAmount string,
	minimum24HourPriceChange string, minimumTimeInSecondsBetweenSells string,
	minimumTokenBalance string, minimumUSDSellPrice string) *TokenConfiguration {

	_minimumUSDSellAmount, err := decimal.NewFromString(minimumUSDSellAmount)
	handler.HandleError("Error when trying to convert minimum USD sell amount from string", err)

	_maximumUSDSellAmount, err := decimal.NewFromString(maximumUSDSellAmount)
	handler.HandleError("Error when trying to convert maximum USD sell amount from string", err)

	_minimum24HourPriceChange, err := decimal.NewFromString(minimum24HourPriceChange)
	handler.HandleError("Error when trying to convert minimum 24 hour price change", err)
	minimum24HourPriceChangeFloat, _ := _minimum24HourPriceChange.Float64()

	_minimumTokenBalance, err := decimal.NewFromString(minimumTokenBalance)
	handler.HandleError("Error when trying to convert minimum 24 hour price change", err)

	_minimumUSDSellPrice, err := decimal.NewFromString(minimumUSDSellPrice)
	handler.HandleError("Error when trying to convert minimum USD sell price", err)
	minimumUSDSellPriceFloat, _ := _minimumUSDSellPrice.Float64()

	return &TokenConfiguration{
		ChainName:                        chainName,
		TokenName:                        tokenName,
		ContractAddress:                  contractAddress,
		LastSellTimestamp:                int64(0),
		MinimumUSDSellAmount:             _minimumUSDSellAmount,
		MaximumUSDSellAmount:             _maximumUSDSellAmount,
		Minimum24HourPriceChange:         minimum24HourPriceChangeFloat,
		MinimumTimeInSecondsBetweenSells: dexbotUtils.StringToInt64(minimumTimeInSecondsBetweenSells),
		MinimumTokenBalance:              _minimumTokenBalance,
		MinimumUSDSellPrice:              minimumUSDSellPriceFloat,
	}
}

//Update the remote user config with the current user config settings
func UpdateRemoteTokenConfig(walletChecksumAddress string, tokenChecksumAddress string) {
	_tokenConfig := UserConfig.Wallets[walletChecksumAddress].Tokens[tokenChecksumAddress].TokenConfig

	_tokenConfigJSON := TokenConfigToJSON(*_tokenConfig)

	//Update database user config
	database.UpdateTokenConfig(walletChecksumAddress, tokenChecksumAddress, _tokenConfigJSON)
}

//
func UpdateLastSellTimestamp(walletChecksumAddress string, tokenChecksumAddress string) {
	lastSellTimestamp := time.Now().Unix()

	//Update the local last sell timestamp
	userWallet := UserConfig.Wallets[walletChecksumAddress]
	token := userWallet.Tokens[tokenChecksumAddress]
	_tokenConfig := token.TokenConfig
	_tokenConfig.LastSellTimestamp = lastSellTimestamp

	//update the remote last sell timestamp
	_tokenConfigJSON := TokenConfigToJSON(*_tokenConfig)
	database.UpdateTokenConfig(walletChecksumAddress, tokenChecksumAddress, _tokenConfigJSON)

}
