package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func main() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("Could not read config/app file")
	} else {
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
		response, err := http.Post(site_url+"oauth2/token", "application/json", bytes.NewBuffer(byteCreds))

		if err != nil {
			fmt.Println("Error", err)
		} else {
			if response.StatusCode == 200 {
				bodyInBytes, _ := ioutil.ReadAll(response.Body)
				json.Unmarshal(bodyInBytes, &jsonBody)
				accessToken := jsonBody["access_token"]

				fmt.Println("access_token received", accessToken)

				record := map[string]string{
					"name": "Gallsaberry",
				}

				recordJSON, _ := json.Marshal(record)
				request, _ := http.NewRequest("POST", site_url+"Accounts", bytes.NewBuffer(recordJSON))
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("oauth-token", accessToken.(string))

				client := &http.Client{}
				response, err := client.Do(request)

				if err != nil {
					fmt.Println("Error on create", err)
				} else {
					if response.StatusCode == 200 {
						bodyInBytes, _ := ioutil.ReadAll(response.Body)
						json.Unmarshal(bodyInBytes, &jsonBody)

						fmt.Println("Created Account with ID", jsonBody["id"])
					}

				}

			}
		}
	}
}
