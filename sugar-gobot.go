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

func main() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not read config/app file")
	}

	site_url := viper.GetString("site_url")
	username := viper.GetString("username")
	password := viper.GetString("password")
	var jsonBody map[string]interface{}

	fmt.Println("Connecting with", username, "to", site_url)

	creds := map[string]string{
		"username":      username,
		"password":      password,
		"grant_type":    "password",
		"client_id":     "sugar",
		"client_secret": "",
		"platform":      "api",
	}
	byteCreds, _ := json.Marshal(creds)
	responseOAuth, err := http.Post(site_url+"oauth2/token", "application/json", bytes.NewBuffer(byteCreds))

	if err != nil {
		log.Fatal(err)
	}

	defer responseOAuth.Body.Close()

	if responseOAuth.StatusCode != 200 {
		log.Fatal(responseOAuth)
	}

	bodyInBytes, _ := ioutil.ReadAll(responseOAuth.Body)
	json.Unmarshal(bodyInBytes, &jsonBody)
	accessToken := jsonBody["access_token"].(string)

	fmt.Println("access_token received", accessToken)

	record := map[string]string{
		"name": "Gallsaberry",
	}

	recordJSON, _ := json.Marshal(record)
	request, _ := http.NewRequest("POST", site_url+"Accounts", bytes.NewBuffer(recordJSON))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("oauth-token", accessToken)

	client := &http.Client{}
	responseCreate, err := client.Do(request)

	if err != nil {
		log.Fatal(responseCreate)
	}

	defer responseCreate.Body.Close()

	if responseCreate.StatusCode != 200 {
		log.Fatal(responseCreate)
	}

	responseCreateInBytes, _ := ioutil.ReadAll(responseCreate.Body)
	json.Unmarshal(responseCreateInBytes, &jsonBody)
	recordID := jsonBody["id"].(string)

	fmt.Println("Created Account with ID", recordID)

	recordUpdate := map[string]string{
		"name": "Chog",
	}

	urlUpdate := site_url + "Accounts/" + recordID

	recordUpdateJSON, _ := json.Marshal(recordUpdate)
	requestUpdate, _ := http.NewRequest("PUT", urlUpdate, bytes.NewBuffer(recordUpdateJSON))
	requestUpdate.Header.Set("Content-Type", "application/json")
	requestUpdate.Header.Add("oauth-token", accessToken)
	fmt.Println("Updating endpoint at", urlUpdate)

	clientUpdate := &http.Client{}
	responseUpdate, err := clientUpdate.Do(requestUpdate)

	if err != nil {
		log.Fatal("Response fail:", responseUpdate)
	}

	defer responseUpdate.Body.Close()

	if responseUpdate.StatusCode != 200 {
		log.Fatal("Update fail:", responseUpdate)
	}

	responseUpdateInBytes, _ := ioutil.ReadAll(responseUpdate.Body)
	json.Unmarshal(responseUpdateInBytes, &jsonBody)

	fmt.Println("Updated Account with ID", jsonBody["id"])
}
