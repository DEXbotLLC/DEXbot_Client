package database

import (
	"bytes"
	"dexbot/internal/authentication"
	"dexbot/internal/config"
	"dexbot/internal/handler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//HTTP client that connects to the database
var RTDBClient *http.Client

//Database URL
var DatabaseURL = config.DatabaseConfig.DatabaseURL

//Initialize the client that connects to the database
func initializeRTDBClient() {

	//Initialize a new HTTP client and set the RTDBClient
	RTDBClient = &http.Client{}
}

//Get data from the database through an authenticated connection
func Get(path string) map[string]interface{} {

	//Send a get request to the database reference with an authenticated url
	resp, err := RTDBClient.Get(fmt.Sprintf("%s/%s.json?auth=%s", DatabaseURL, path, authentication.FirebaseAuthToken.IdToken))
	httpErr := handler.CheckHTTPResponse("database: Get: RTDBClient.Get", resp, err)
	//Close the body of the response after the function ends
	defer resp.Body.Close()

	//if there is an error, try again
	if httpErr != nil {
		return Get(path)
	}

	//Convert the response body to bytes
	dataBytes, err := ioutil.ReadAll(resp.Body)
	handler.HandleError("database: Get: ioutil.ReadAll", err)

	//Unmarshal the bytes to a map
	data := make(map[string]interface{})
	json.Unmarshal(dataBytes, &data)

	//Return the data as a map
	return data

}

func Update(path string, data map[string]interface{}) {
	//Convert the data payload to JSON bytes
	jsonBytes, err := json.Marshal(data)
	handler.HandleError("database: Update: Marshal", err)

	//Create a PUT request to the database reference with an authenticated url
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s.json?auth=%s", DatabaseURL, path, authentication.FirebaseAuthToken.IdToken), bytes.NewBuffer(jsonBytes))
	handler.HandleError("database: Update: NewRequest", err)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	//Send the PUT request to the database
	resp, err := RTDBClient.Do(req)
	httpErr := handler.CheckHTTPResponse("database: Update: client.Do", resp, err)

	//Close the body of the response after the function ends
	defer resp.Body.Close()

	//if there is an error, try again
	if httpErr != nil {
		Update(path, data)
	}
}

func Delete(path string) {
	//send a DELETE request to the database reference with an authenticated url

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s.json?auth=%s", config.DatabaseConfig.DatabaseURL, path, authentication.FirebaseAuthToken.IdToken), nil)
	handler.HandleError("database: Delete: NewRequest", err)

	resp, err := RTDBClient.Do(req)
	httpErr := handler.CheckHTTPResponse("database: Delete: client.Do", resp, err)
	defer resp.Body.Close()

	//if there is an EOF error, try again
	if httpErr != nil && httpErr.Error() == "EOF" {
		Delete(path)
	}

}
