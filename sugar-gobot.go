package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/Pallinder/go-randomdata"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var jsonBody map[string]interface{}
var accessToken, siteURL string

type MassRecords struct {
	Params struct {
		UID []string `json:"uid"`
	} `json:"massupdate_params"`
}

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

func massDelete(module string, records []string) bool {
	var massRecords MassRecords
	massRecords.Params.UID = records
	recordJSON, _ := json.Marshal(massRecords)
	request, _ := http.NewRequest("DELETE", siteURL+module+"/MassUpdate", bytes.NewBuffer(recordJSON))
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

	log.Println("Delete", len(records), module, "records")

	return true
}

func linkRecords(module, recordID, relatedModule, relatedID string) bool {
	request, _ := http.NewRequest("POST", siteURL+module+"/"+recordID+"/link/"+strings.ToLower(relatedModule)+"/"+relatedID, nil)
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

	log.Println("Linked", module, "ID:", recordID, "to", relatedModule, "ID:", relatedID)

	return true
}

func connect(siteURL, username, password string) {
	log.Println("Connect", username, "to", siteURL)

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

	log.Println("Use access_token", accessToken)
}

func main() {
	recordsToCreate := flag.Int("num", 10, "Count of records to create on each run")
	purgeAfterCreation := flag.Bool("delete", false, "Whether to delete the records after creation")
	flag.Parse()

	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	var createdAccounts []string
	var createdContacts []string

	if err != nil {
		log.Fatal("Could not read config/app file")
	}

	siteURL = viper.GetString("site_url")
	username := viper.GetString("username")
	password := viper.GetString("password")

	connect(siteURL, username, password)

	for i := 0; i < *recordsToCreate; i++ {
		accountData := map[string]string{
			"name": strings.Title(randomdata.Adjective() + " " + randomdata.Noun() + " " + randomdata.City()),
		}

		accountID := createRecord("Accounts", accountData)
		createdAccounts = append(createdAccounts, accountID)
		log.Println("Accounts created:", len(createdAccounts))

		contactData := map[string]string{
			"first_name": randomdata.FirstName(randomdata.RandomGender),
			"last_name":  randomdata.LastName(),
		}

		contactID := createRecord("Contacts", contactData)
		createdContacts = append(createdContacts, contactID)
		log.Println("Contacts created:", len(createdContacts))

		linkRecords("Accounts", accountID, "Contacts", contactID)
	}

	if *purgeAfterCreation {
		log.Println("Time to purge newly created records")
		massDelete("Accounts", createdAccounts)
		massDelete("Contacts", createdContacts)
	}
}
