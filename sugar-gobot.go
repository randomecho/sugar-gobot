package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

var jsonBody map[string]interface{}
var accessToken, siteURL string

func createRecord(module string, record map[string]string) string {
	recordJSON, _ := json.Marshal(record)
	request, _ := http.NewRequest("POST", siteURL+module, bytes.NewBuffer(recordJSON))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("oauth-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(response)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal(response)
	}

	responseInBytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseInBytes, &jsonBody)
	recordID := jsonBody["id"].(string)

	log.Println("Create", module, "record with ID:", recordID)

	return recordID
}

func updateRecord(module string, recordID string, record map[string]string) string {
	recordJSON, _ := json.Marshal(record)
	request, _ := http.NewRequest("PUT", siteURL+module+"/"+recordID, bytes.NewBuffer(recordJSON))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("oauth-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal("Response fail:", response)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal("Update fail:", response)
	}

	responseInBytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseInBytes, &jsonBody)
	recordID = jsonBody["id"].(string)

	log.Println("Update", module, "record with ID:", recordID)

	return recordID
}

func deleteRecord(module, recordID string) bool {
	request, _ := http.NewRequest("DELETE", siteURL+module+"/"+recordID, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("oauth-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal("Response fail:", response)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal("Delete fail:", response)
	}

	responseInBytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseInBytes, &jsonBody)

	log.Println("Delete", module, "record with ID:", recordID)

	return true
}

func main() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not read config/app file")
	}

	siteURL = viper.GetString("site_url")
	username := viper.GetString("username")
	password := viper.GetString("password")

	fmt.Println("Connecting with", username, "to", siteURL)

	creds := map[string]string{
		"username":      username,
		"password":      password,
		"grant_type":    "password",
		"client_id":     "sugar",
		"client_secret": "",
		"platform":      "api",
	}
	byteCreds, _ := json.Marshal(creds)
	responseOAuth, err := http.Post(siteURL+"oauth2/token", "application/json", bytes.NewBuffer(byteCreds))

	if err != nil {
		log.Fatal(err)
	}

	defer responseOAuth.Body.Close()

	if responseOAuth.StatusCode != 200 {
		log.Fatal(responseOAuth)
	}

	bodyInBytes, _ := ioutil.ReadAll(responseOAuth.Body)
	json.Unmarshal(bodyInBytes, &jsonBody)
	accessToken = jsonBody["access_token"].(string)

	fmt.Println("access_token received", accessToken)

	recordData := map[string]string{
		"name": "Gallsaberry",
	}

	recordID := createRecord("Accounts", recordData)

	recordUpdate := map[string]string{
		"name": "Chog",
	}

	updateRecord("Accounts", recordID, recordUpdate)

	deleteRecord("Accounts", recordID)
}
