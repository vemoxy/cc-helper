package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vemoxy/cc-helper/config"
	"github.com/vemoxy/cc-helper/model"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
)

var mccFileName string = "mcc_codes.json"

/*
   "mcc": "0742",
   "edited_description": "Veterinary Services",
   "combined_description": "Veterinary Services",
   "usda_description": "Veterinary Services",
   "irs_description": "Veterinary Services",
   "irs_reportable": "Yes",
   "id": 0
*/

type TransactionTypeFromJson struct {
	Mcc                 string `json:"mcc"`
	EditedDescription   string `json:"edited_description"`
	CombinedDescription string `json:"combined_description"`
	UsdaDescription     string `json:"usda_description"`
	IrsDescription      string `json:"irs_description"`
	IrsReportable       string `json:"irs_reportable"`
	Id                  int    `json:"id"`
}

func ImportMcc() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(mccFileName); err != nil {
		req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/greggles/mcc-codes/main/mcc_codes.json", nil)
		if err != nil {
			panic(err)
		}

		// Send the request and get the response
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// Create a new file to save the downloaded data to
		file, err := os.Create(mccFileName)
		if err != nil {
			panic(err)
		}

		// Copy the downloaded data to the file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			panic(err)
		}

		file.Close()
	}

	// Read the JSON file
	file, err := os.ReadFile(mccFileName)
	if err != nil {
		log.Fatal(err)
	}

	var fileContents []TransactionTypeFromJson
	err = json.Unmarshal(file, &fileContents)
	if err != nil {
		log.Fatal(err)
	}

	// FreeLeh auth
	auth, err := auth.NewServiceFromFile(
		config.GoogleServiceAccountKeyFileName,
		freedb.FreeDBGoogleAuthScopes,
		auth.ServiceConfig{},
	)
	if err != nil {
		log.Fatal(err)
	}

	store := freedb.NewGoogleSheetRowStore(
		auth,
		conf.FreeDB.SpreadsheetId,
		conf.FreeDB.Sheets.Mcc,
		freedb.GoogleSheetRowStoreConfig{Columns: model.MerchantCategoryCodeFields},
	)

	for _, record := range fileContents {
		rawData, err := json.Marshal(record)
		if err != nil {
			log.Printf("Error unmarshal record, err: %s, record: %v", err.Error(), record)
			continue
		}

		mccInt, err := strconv.Atoi(record.Mcc)
		if err != nil {
			log.Printf("Error convert mcc to int, err: %s, mcc: %s", err.Error(), record.Mcc)
			continue
		}

		mccModel := model.MerchantCategoryCode{
			Mcc:         mccInt,
			Description: record.EditedDescription,
			Metadata:    string(rawData),
		}
		err = store.Insert(mccModel).Exec(context.Background())
		if err != nil {
			log.Fatalf("Error writing record to gsheet, err: %s, model: %v", err.Error(), mccModel)
		}

		log.Printf("Inserted record: %v", record)
		time.Sleep(time.Second)
	}
}
