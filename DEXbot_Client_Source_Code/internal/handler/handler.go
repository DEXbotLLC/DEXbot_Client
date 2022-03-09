package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/term"
)

//Global wait group to manage goroutines
var WaitGroup sync.WaitGroup

//Add one goroutine to waitgroup
func AddGoroutineToWaitGroup() {
	WaitGroup.Add(1)
}

//Wait for every instance in the wait group to finish.
func WaitForWaitGroupToFinish() {
	WaitGroup.Wait()
}

//Prompt the user for a terminal input and return the value
func Input(message string) string {

	//Display the message prompt
	fmt.Printf(" %s", message)

	//Initialize a new reader
	reader := bufio.NewReader(os.Stdin)

	//Wait for the user input, ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	HandleError("handler: Input: ReadString", err)

	//Trim the new line delimiter from the input
	input = strings.Replace(input, "\r", "", -1)
	input = strings.Replace(input, "\n", "", -1)

	//Return the user input
	return input
}

//Prompt the user for a terminal input while obscuring the input and return the value
func InputObscured(message string) string {

	//Display the message prompt
	fmt.Print(message)

	//Wait for the user input, text is obscured while entering the input and returned as bytes
	byteInput, err := term.ReadPassword(int(os.Stdin.Fd()))
	HandleError("Error when entering value", err)

	//Return the input as a string
	return string(byteInput)
}

//Prompt the user for a terminal input while obscuring the input and return the value as bytes
//This allows for the input to be "zeroed", wiping the input
func InputPrivateKey(message string) []byte {

	//Display the message prompt
	fmt.Print(message)

	//Wait for the user input, text is obscured while entering the input and returned as bytes
	privateKeyBytes, err := term.ReadPassword(int(syscall.Stdin))
	HandleError("Error when entering value", err)

	//Trim the 0x from the private key and format the byte slice
	trimmedPrivateKeyBytes := []byte{}
	for i := 0; i < 64; i++ {
		trimmedPrivateKeyBytes = append(trimmedPrivateKeyBytes, privateKeyBytes[len(privateKeyBytes)-(i+1)])
	}
	for i, j := 0, len(trimmedPrivateKeyBytes)-1; i < j; i, j = i+1, j-1 {
		trimmedPrivateKeyBytes[i], trimmedPrivateKeyBytes[j] = trimmedPrivateKeyBytes[j], trimmedPrivateKeyBytes[i]
	}

	//Zero the byte slice containing the private key
	for i := range privateKeyBytes {
		privateKeyBytes[i] = 0
	}

	//Return the input as bytes
	return trimmedPrivateKeyBytes
}

//Remove leading zeros by getting the last 40 characters in an address and prepend 0x to the beginning of the string
func RemoveLeadingZeros(inputString string) string {
	strippedAddress := ""

	//Get the last 40 characters in the string
	for i := 0; i < 40; i++ {
		strippedAddress = string(inputString[len(inputString)-(i+1)]) + strippedAddress
	}

	//prepend 0x to the beginning of the string
	return "0x" + strippedAddress

}

//Error reporting boolean to toggle error reporting on and off
var errorReportingBool bool

//User Id variable to be sent with an error if error reporting is toggled on
var errorReportingUserId string

//Initiailze errror reporting setting and set userId
func InitializeErrorReporting(errorReporting bool, _userId string) {
	errorReportingBool = errorReporting
	errorReportingUserId = _userId
}

//Toggle error reporting on and off
func ToggleErrorReporting(errorReporting bool) {
	errorReportingBool = errorReporting
}

//Handle errors that occur, printing the issue in the terminal. If error reporting is toggled on, errors will be sent to DEXbot
func HandleError(context string, err error) {
	if err != nil {
		fmt.Println("------------------------")
		fmt.Println("Context: ", context)
		fmt.Println("Error: ", err)
		fmt.Println("------------------------")
		//If the user has error reporting enabled, send the error to the DEXbot team
		if errorReportingBool {
			SendErrorToDiscord(context, err.Error(), errorReportingUserId)
		}
		os.Exit(1)
	}
}

//Send errors to DEXbot development discord (only when error reporting is toggled on)
func SendErrorToDiscord(context string, errorMessage string, userId string) {

	//Create the message as a json string
	jsonString := []byte(fmt.Sprintf(`{"content": null,"embeds": [{"title": "%s, %s","description": "%s","color": 16734296}]}`, userId, context, errorMessage))

	//Create a post request
	request, err := http.NewRequest("POST", "https://discordapp.com/api/webhooks/919075558929883157/I7TFvchM_BVPSg9GNlmwGONxZoFDxKhesDx7hFHxoZuy3dyR-bZ4fDs7humrLpIhWxgw", bytes.NewBuffer(jsonString))
	HandleError("handler: SendErrorToDiscord: NewRequest", err)
	request.Header.Set("Content-Type", "application/json")

	//Initialize a new http client
	client := &http.Client{}

	//Send the payload
	response, err := client.Do(request)
	HandleError("handler: SendErrorToDiscord: client.Do", err)

	//Close the response body after the function finishes
	defer response.Body.Close()
}

//Handle the response from an authentication request
func CheckAuthResponse(authResponse string) {

	//If there is an error notifying an invalid email or password, display a message in the terminal
	if strings.Contains(authResponse, "error") {
		if strings.Contains(authResponse, "INVALID_EMAIL") || strings.Contains(authResponse, "INVALID_PASSWORD") {
			fmt.Println("Incorrect email or password, please check your credentials and try again.")
		}

		//Exit the program with a exit message
		Exit("Incorrect Firebase credentials, exiting program.")
	}
}

//Handle the response from an HTTP request
func CheckHTTPResponse(context string, httpResponse *http.Response, err error) error {
	//Check if there are any errors when sending the http request
	if err != nil {

		//If the user has error reporting enabled, send the error to the DEXbot team
		if errorReportingBool {
			SendErrorToDiscord(context, err.Error(), errorReportingUserId)
		}

		//return the error
		return err

	} else {
		//If there is an error in the HTTP response
		if strings.Contains(httpResponse.Status, "4") {

			//Print the error in the terminal
			fmt.Println(httpResponse)

			//If the user has error reporting enabled, send the error to the DEXbot team
			if errorReportingBool {
				SendErrorToDiscord(context, err.Error(), errorReportingUserId)
			}

			//Exit the program with an exit message
			Exit("HTTP response error, exiting program")
		}
	}
	return nil
}

//Print a message and terminate the program
func Exit(message string) {
	fmt.Println(message)
	os.Exit(0)
}
