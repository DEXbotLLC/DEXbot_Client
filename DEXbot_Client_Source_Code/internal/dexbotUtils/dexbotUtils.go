package dexbotUtils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"dexbot/internal/handler"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/shopspring/decimal"
	"github.com/umbracle/go-web3"
	"golang.org/x/crypto/sha3"
)

//Variable to hold the operating system Id for OS specific funcationality compatiblilty
var OperatingSystem string

var DEXbotClientChecksum string

//Intialize the checksum for the client application
func initalizeDEXbotChecksum() {
	//Get the version checksum so DEXbot can validate the integrity of the application
	DEXbotClientChecksum = RunChecksumOnDEXbotBinary()
}

//Function to initialize operating system
func InitializeOperatingSystem() {
	OperatingSystem = runtime.GOOS
}

//Printer that prints in white
var WhitePrinter = color.New(color.FgWhite, color.Bold)

//Printer that prints in blue
var BluePrinter = color.New(color.FgBlue, color.Bold)

//Printer that prints in green
var GreenPrinter = color.New(color.FgGreen, color.Bold)

//Printer that prints in yellow
var YellowPrinter = color.New(color.FgYellow, color.Bold)

//Convert a hex address to checksum address
func ToChecksumAddress(address string) (string, error) {

	//Check that the address is a valid Ethereum address
	re1 := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re1.MatchString(address) {
		return "", fmt.Errorf("given address '%s' is not a valid Ethereum Address", address)
	}

	//Convert the address to lowercase
	re2 := regexp.MustCompile("^0x")
	address = re2.ReplaceAllString(address, "")
	address = strings.ToLower(address)

	//Convert address to sha3 hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(address))
	sum := hasher.Sum(nil)
	addressHash := fmt.Sprintf("%x", sum)
	addressHash = re2.ReplaceAllString(addressHash, "")

	//Compile checksum address
	checksumAddress := "0x"
	for i := 0; i < len(address); i++ {
		indexedValue, err := strconv.ParseInt(string(rune(addressHash[i])), 16, 32)
		if err != nil {
			fmt.Println("Error when parsing addressHash during checksum conversion", err)
			return "", err
		}
		if indexedValue > 7 {
			checksumAddress += strings.ToUpper(string(address[i]))
		} else {
			checksumAddress += string(address[i])
		}
	}

	//Return the checksummed address
	return checksumAddress, nil
}

//Check if a hex address is a checksum address
func IsChecksumAddress(address string) (bool, error) {
	//Run a the address through the ToChecksumAddress to get the proper checksummed address
	checksumAddress, err := ToChecksumAddress(address)
	//Check if there were any errors when converting to checksum
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	//Check if the original address matches the checksummed address
	if checksumAddress == address {
		return true, nil
	} else {
		return false, nil
	}
}

// Parses private key bytes to an ECDSA Key
func BytesToECDSA(byteSlice []byte) (*ecdsa.PrivateKey, error) {

	//Decode the private key bytes
	n, err := hex.Decode(byteSlice, byteSlice)
	b := byteSlice[:n]

	//Check the byte slice for invalid characters
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}

	//Return an ECDSA key
	return crypto.ToECDSA(b)
}

//Check the checksummed value of the dexbot client application binary
func RunChecksumOnDEXbotBinary() string {
	if OperatingSystem == "linux" {
		dexbotBinaryChecksum, err := sha256sum("dexbot-linux")
		handler.HandleError("Error when getting checksum from application binary", err)
		return dexbotBinaryChecksum
	} else if OperatingSystem == "darwin" {
		dexbotBinaryChecksum, err := sha256sum("dexbot-mac")
		handler.HandleError("Error when getting checksum from application binary", err)
		return dexbotBinaryChecksum
	} else if OperatingSystem == "windows" {
		dexbotBinaryChecksum, err := sha256sum("dexbot-windows.exe")
		handler.HandleError("Error when getting checksum from application binary", err)
		return dexbotBinaryChecksum
	} else {
		return ""
	}
}

//Check if a specific string is in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//Get the sha256 hash of a file
func sha256sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

//Get the method signature from transaction input data in hex format
func GetMethodSignature(inputData []byte) string {
	inputDataHex := hex.EncodeToString(inputData)

	if len(inputDataHex) > 0 {
		return inputDataHex[:8]
	} else {
		return ""
	}
}

func UnpackTransactionPayload(transactionData string) *web3.Transaction {
	//Convert the transaction data from an encoded string to bytes
	transactionBytes, err := base64.StdEncoding.DecodeString(transactionData)
	handler.HandleError("Error when converting transaction to bytes", err)

	//Initialize an unsigned transaction variable
	transaction := &web3.Transaction{}

	//Unmarshal the unsigned transaction bytes to the unsignedTransaction variable
	err = json.Unmarshal(transactionBytes, transaction)

	// if there is no nonce, i.e. nonce is 0
	if err != nil && strings.Contains(err.Error(), "field 'nonce' not found") {
		//add the nonce to the transaction
		transactionBytesWithNonce := []byte(fmt.Sprintf("%s,%s}", string(transactionBytes[:len(string(transactionBytes))-1]), `"nonce":"0x0"`))

		//Initialize an unsigned transaction variable
		transaction = &web3.Transaction{}

		//Unmarshal the unsigned transaction bytes to the unsignedTransaction variable
		err = json.Unmarshal(transactionBytesWithNonce, transaction)
		handler.HandleError("Error when converting transaction with nonce=0 to bytes", err)

		return transaction

	} else if err != nil {
		handler.HandleError("Error when unmarshaling unsigned transaction", err)
	}

	return transaction
}

func StringToInt64(input string) int64 {
	newDecimal, err := decimal.NewFromString(input)
	handler.HandleError("handler: StringToInt64: decimal.NewFromString", err)
	newFloat64, _ := newDecimal.BigFloat().Int64()
	return newFloat64
}
