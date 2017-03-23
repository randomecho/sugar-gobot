package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
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
			fmt.Println("Response", response)
		}
	}
}
