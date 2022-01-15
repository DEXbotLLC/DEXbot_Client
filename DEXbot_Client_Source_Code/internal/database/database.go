package database

import (
	"bytes"
	"dexbot/internal/authentication"
	"dexbot/internal/handler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//HTTP client that connects to the database
var RTDBClient *http.Client

//Database URL
var DatabaseURL = "https://dexbot-90461-default-rtdb.firebaseio.com"

//Initialize the client that connects to the database
func initializeRTDBClient() {

	//Initialize a new HTTP client and set the RTDBClient
	RTDBClient = &http.Client{}
}

//Get data from the database through an authenticated connection
func Get(path string) map[string]interface{} {

	//Send a get request to the database reference with an authenticated url
	resp, err := RTDBClient.Get(fmt.Sprintf("%s/%s.json?auth=%s", DatabaseURL, path, authentication.FirebaseAuthToken.IdToken))
	handler.HandleError("database: Get: Get", err)
	//Close the response after the function ends
	defer resp.Body.Close()

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
	handler.HandleError("database: Update: client.Do", err)

	//Check the response for errors
	handler.CheckHTTPResponse(resp)

	//Close the body of the response
	defer resp.Body.Close()

}
