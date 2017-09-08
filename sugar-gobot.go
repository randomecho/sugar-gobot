package main

import (
	"bytes"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var jsonBody map[string]interface{}
var accessToken, siteURL string

type SugarRecord struct {
	Id           string `json:"id"`
	DateEntered  string `json:"date_entered"`
	DateModified string `json:"date_modified"`
}

type SearchResult struct {
	NextOffset int           `json:"next_offset"`
	Records    []SugarRecord `json:"records"`
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

func getRecordByFields(module string, record map[string]string) string {
	var result SearchResult
	filterBy := url.Values{}

	for k, v := range record {
		filterBy.Add("filter[]["+k+"]", v)
	}

	request, _ := http.NewRequest("GET", siteURL+module+"/filter?"+filterBy.Encode(), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("oauth-token", accessToken)

	log.Println("Search", module, "by /filter?"+filterBy.Encode())

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal("Response fail:", response)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal("Search fail:", response)
	}

	responseInBytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseInBytes, &result)

	if len(result.Records) == 0 {
		log.Println("No", module, "record found matching search criteria")

		return "0"
	} else {

		log.Println("Found", module, "record with ID:", result.Records[0].Id)

		return result.Records[0].Id
	}
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
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not read config/app file")
	}

	siteURL = viper.GetString("site_url")
	username := viper.GetString("username")
	password := viper.GetString("password")

	connect(siteURL, username, password)

	accountData := map[string]string{
		"name": strings.Title(randomdata.Adjective() + " " + randomdata.Noun() + " " + randomdata.City()),
	}

	accountID := createRecord("Accounts", accountData)

	contactData := map[string]string{
		"first_name": randomdata.FirstName(randomdata.RandomGender),
		"last_name": randomdata.LastName(),
	}

	contactID := createRecord("Contacts", contactData)

	linkRecords("Accounts", accountID, "Contacts", contactID)

	recordLookup := map[string]string{
		"name": "Poyo",
	}

	getRecordByFields("Accounts", recordLookup)

	recordUpdate := map[string]string{
		"name": "Chog",
	}

	updateRecord("Accounts", recordID, recordUpdate)

	deleteRecord("Accounts", recordID)
}
