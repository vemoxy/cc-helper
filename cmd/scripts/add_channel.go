package main

import (
	"context"
	"log"

	"github.com/vemoxy/cc-helper/config"
	"github.com/vemoxy/cc-helper/model"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
)

func AddChannel(merchantId int, channelName string, mcc int, source string, updateTime int, updateBy string) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// FreeLeh auth
	auth, err := auth.NewServiceFromFile(
		"google-service-account-key.json",
		freedb.FreeDBGoogleAuthScopes,
		auth.ServiceConfig{},
	)
	if err != nil {
		log.Fatal(err)
	}

	store := freedb.NewGoogleSheetRowStore(
		auth,
		config.FreeDB.SpreadsheetId,
		config.FreeDB.Sheets.Channel,
		freedb.GoogleSheetRowStoreConfig{Columns: model.ChannelFields},
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

	channelModel := model.Channel{
		Id:          maxId,
		MerchantId:  merchantId,
		ChannelName: channelName,
		Mcc:         mcc,
		Source:      source,
		UpdateTime:  updateTime,
		UpdateBy:    updateBy,
	}
	err = store.Insert(channelModel).Exec(context.Background())
	if err != nil {
		log.Fatalf("Error writing record to gsheet, err: %s, model: %v", err.Error(), channelModel)
	}

	log.Printf("Inserted model: %v", channelModel)

}
