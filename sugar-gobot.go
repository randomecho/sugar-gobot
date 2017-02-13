package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("Could not read config/app file")
	} else {
		site_url := viper.GetString("site_url")

		fmt.Println("Connecting with", site_url)
	}
}
