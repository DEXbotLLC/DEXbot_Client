package authentication

import (
	"bytes"
	"dexbot/internal/config"
	"dexbot/internal/handler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//The firebase username and password for refreshing the authentication token to enable connection to the database
var firebaseCreds *firebaseCredentials

//Authentication token to connect to the database
var FirebaseAuthToken *authToken

//API key for the database
var firebaseAPIKey = config.DatabaseConfig.DatabaseAPIKey

//Channel to notify the program when to refresh the authentication token
var AuthRefreshChannel chan bool

//Initializes firebaseCreds with the user name and password
func initalizeFirebaseCreds(username string, password string) {
	firebaseCreds = &firebaseCredentials{firebaseEmail: username, firebasePassword: password}
	//initialize auth refresh channel
	AuthRefreshChannel = make(chan bool)
	fmt.Println()
}

//Initializes the authentication token and refreshes the token every 50 minutes
func initializeAuthTokenWithRefresh() {
	//Set the authURL with the database API key
	authURL := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", firebaseAPIKey)
	requestBody := []byte(fmt.Sprintf(`{"email":"%v","password":"%v","returnSecureToken":true}`, firebaseCreds.firebaseEmail, firebaseCreds.firebasePassword))

	//Create a post request to the auth endpoint
	request, err := http.NewRequest("POST", authURL, bytes.NewBuffer(requestBody))
	handler.HandleError("authenticate: GetAuthToken: http.NewRequest", err)
	request.Header.Set("Content-Type", "application/json")

	//Initialize a new http client and send the auth request
	client := &http.Client{}
	response, err := client.Do(request)
	handler.HandleError("authenticate: GetAuthToken: client.Do", err)
	//Close the response after the function ends
	defer response.Body.Close()

	//Read in the response body and unmarshal to an auth response
	authBytes, err := ioutil.ReadAll(response.Body)
	handler.CheckAuthResponse(string(authBytes))
	handler.HandleError("authenticate: Token: ioutil.ReadAll", err)

	//Get the auth token from the authBytes and set the FirebaseAuthToken
	token := authToken{}
	json.Unmarshal(authBytes, &token)
	FirebaseAuthToken = &token

	//Start the cycle to refresh the token every 50 minutes
	go refreshAuthtoken()
}

//Refresh the authentication token every 50 minutes
func refreshAuthtoken() {
	//As an infinite loop
	for {
		//Sleep for 50 minutes
		time.Sleep(time.Second * 3000)

		//Set the authURL with the database API key
		authURL := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s", firebaseAPIKey)
		requestBody := []byte(fmt.Sprintf(`{"refresh_token":"%v","grant_type":"refresh_token"}`, FirebaseAuthToken.RefreshToken))

		//Create a post request to the auth endpoint
		request, err := http.NewRequest("POST", authURL, bytes.NewBuffer(requestBody))
		handler.HandleError("authenticate: GetAuthToken: http.NewRequest", err)
		request.Header.Set("Content-Type", "application/json")

		//Initialize a new http client and send the auth request
		client := &http.Client{}
		response, err := client.Do(request)
		handler.HandleError("authenticate: GetAuthToken: client.Do", err)

		//Read in the response body and unmarshal to an auth response
		authBytes, err := ioutil.ReadAll(response.Body)
		handler.CheckAuthResponse(string(authBytes))
		handler.HandleError("authenticate: Token: ioutil.ReadAll", err)

		//Get the auth token from the authBytes and set the FirebaseAuthToken.IdToken
		refreshToken := refreshToken{}
		json.Unmarshal(authBytes, &refreshToken)
		FirebaseAuthToken.IdToken = refreshToken.Id_token

		//Send a notification to the AuthRefreshChannel to refresh all of the authenticated connections
		AuthRefreshChannel <- true

		//Close the response body
		response.Body.Close()
	}
}
