package main

import (
	"context"
	"log"

	"github.com/vemoxy/cc-helper/config"
	"github.com/vemoxy/cc-helper/model"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
)

func AddMerchant(merchantName string) {
	conf, err := config.LoadConfig()
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
		conf.FreeDB.Sheets.Merchant,
		freedb.GoogleSheetRowStoreConfig{Columns: model.MerchantFields},
	)

	var maxId int
	var currentMaxIdOutput []model.Merchant
	ordering := []freedb.ColumnOrderBy{
		{Column: "id", OrderBy: freedb.OrderByDesc},
	}

	err = store.Select(&currentMaxIdOutput, "id").OrderBy(ordering).Limit(1).Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(currentMaxIdOutput) == 0 {
		maxId = 1
	} else {
		maxId = currentMaxIdOutput[0].Id + 1
	}

	merchantModel := model.Merchant{
		Id:   maxId,
		Name: merchantName,
	}
	err = store.Insert(merchantModel).Exec(context.Background())
	if err != nil {
		log.Fatalf("Error writing record to gsheet, err: %s, model: %v", err.Error(), merchantModel)
	}

	log.Printf("Inserted model: %v", merchantModel)

}
