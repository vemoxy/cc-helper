package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	FreeDB   FreeDBConfig   `json:"freedb"`
}

type TelegramConfig struct {
	ApiKey string `json:"api_key"`
	Debug  bool   `json:"debug"`
}

type FreeDBConfig struct {
	SpreadsheetId string     `json:"spreadsheet_id"`
	Sheets        SheetNames `json:"sheets"`
}

type SheetNames struct {
	Mcc             string `json:"mcc"`
	MerchantChannel string `json:"merchantChannel"`
}

var configFileName string = "config.json"

func LoadConfig() (*Config, error) {
	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &config, nil
}

var GoogleServiceAccountKeyFileName string = "google-service-account-key.json"
