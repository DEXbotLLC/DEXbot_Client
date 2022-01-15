package eventListener

import (
	"dexbot/internal/handler"
	"encoding/json"
	"strings"
)

//Convert eventData from bytes to a map
func unmarshalEventData(eventData []byte) map[string]interface{} {

	//Initalize a new map to store the converted eventData
	var eventDataMap map[string]interface{}

	//Unmarshal eventData bytes to eventDataMap
	err := json.Unmarshal(eventData, &eventDataMap)
	handler.HandleError("eventHandler: unmarshalEventData: json.Unmarshal (eventData)", err)

	//If the event data exists, return the map
	if eventData != nil {
		return eventDataMap
	} else {
		return nil
	}
}

//Get the reference path from the payload
func getReferenceTree(eventData []byte) []string {

	//Intialize a map to unmarshal the payload path data into
	var eventDataMap map[string]interface{}

	//Unmarshal the payload path data
	err := json.Unmarshal(eventData, &eventDataMap)
	handler.HandleError("eventHandler: getEventPath: json.Unmarshal (eventData)", err)

	//Get the eventData path as a string
	eventPath := eventDataMap["path"].(string)

	//Parse the reference path
	referenceTree := strings.Split(eventPath, "/")
	filteredReferenceTree := []string{}

	//Filter out zero values from reference path and append the path to a list of strings
	for _, v := range referenceTree {
		if v != "" {
			filteredReferenceTree = append(filteredReferenceTree, v)
		}
	}

	//Return the reference tree
	return filteredReferenceTree
}
