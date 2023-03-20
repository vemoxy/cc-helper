package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
	"github.com/vemoxy/cc-helper/config"
	"github.com/vemoxy/cc-helper/data"
	"github.com/vemoxy/cc-helper/model"
)

type CheckMccResult struct {
	Merchant    string
	ChannelName string
	Mcc         int
}

func CheckMcc(merchantName string) []CheckMccResult {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Crashed, err: %v", r)
		}
	}()
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	auth, err := auth.NewServiceFromFile(
		config.GoogleServiceAccountKeyFileName,
		freedb.FreeDBGoogleAuthScopes,
		auth.ServiceConfig{},
	)
	if err != nil {
		log.Fatal(err)
	}

	var checkMccResult []CheckMccResult

	merchantStore := freedb.NewGoogleSheetRowStore(
		auth,
		conf.FreeDB.SpreadsheetId,
		conf.FreeDB.Sheets.Merchant,
		freedb.GoogleSheetRowStoreConfig{Columns: model.MerchantFields},
	)

	var queriedMerchantName string
	var merchantId int
	var merchantQueryOutput []model.Merchant

	err = merchantStore.Select(&merchantQueryOutput).
		Where("name CONTAINS ?", merchantName).
		Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(merchantQueryOutput) == 0 {
		return checkMccResult
	} else {
		queriedMerchantName = merchantQueryOutput[0].Name
		merchantId = merchantQueryOutput[0].Id
	}

	channelStore := freedb.NewGoogleSheetRowStore(
		auth,
		conf.FreeDB.SpreadsheetId,
		conf.FreeDB.Sheets.Channel,
		freedb.GoogleSheetRowStoreConfig{Columns: model.ChannelFields},
	)
	var channelQueryOutput []model.Channel

	err = channelStore.Select(&channelQueryOutput).
		Where("merchantId = ?", merchantId).
		Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(channelQueryOutput) > 0 {
		for _, queryRow := range channelQueryOutput {
			checkMccResult = append(checkMccResult, CheckMccResult{
				Merchant:    queriedMerchantName,
				ChannelName: queryRow.ChannelName,
				Mcc:         queryRow.Mcc,
			})
		}
	}

	return checkMccResult
}

func CheckMccV2(merchantName string) []CheckMccResult {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Crashed, err: %v", r)
		}
	}()

	var result []CheckMccResult

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	auth, err := auth.NewServiceFromFile(
		config.GoogleServiceAccountKeyFileName,
		freedb.FreeDBGoogleAuthScopes,
		auth.ServiceConfig{},
	)
	if err != nil {
		log.Fatal(err)
	}

	merchantChannelStore := freedb.NewGoogleSheetRowStore(
		auth,
		conf.FreeDB.SpreadsheetId,
		conf.FreeDB.Sheets.MerchantChannel,
		freedb.GoogleSheetRowStoreConfig{Columns: model.MerchantChannelFields},
	)

	var queryOutput []model.MerchantChannel

	err = merchantChannelStore.Select(&queryOutput).
		Where("merchantName CONTAINS ?", merchantName).
		Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(queryOutput) == 0 {
		return result
	} else {
		for _, outputRow := range queryOutput {
			result = append(result, CheckMccResult{
				Merchant:    outputRow.MerchantName,
				ChannelName: outputRow.ChannelName,
				Mcc:         outputRow.Mcc,
			})
		}
	}

	return result
}

func GenerateMccResultMessage(results []CheckMccResult) string {
	var message string
	if len(results) == 0 {
		message = "Merchant not found in database."
	} else {
		textSlice := []string{"Found results:"}
		for _, result := range results {
			textSlice = append(
				textSlice,
				fmt.Sprintf("merchant: %s, channelName: %s, mcc: %d", result.Merchant, result.ChannelName, result.Mcc),
			)
		}
		message = strings.Join(textSlice, "\n")
	}
	return message
}

func CheckMccV3(merchantName string) []model.MerchantChannel {
	results := data.QueryMerchantChannelsByMerchantName(merchantName)

	return results
}

func GenerateMccV3ResultMessage(results []model.MerchantChannel) string {
	var message string
	if len(results) > 0 {
		messages := []string{"Found results:\n"}
		for _, match := range results {
			merchantDesc := fmt.Sprintf("*%s*", match.MerchantName)
			if len(match.ChannelName) > 0 {
				merchantDesc += fmt.Sprintf(" (%s)", match.ChannelName)
			}
			messages = append(messages, fmt.Sprintf("Merchant: %s, MCC: %d", merchantDesc, match.Mcc))
		}
		message = strings.Join(messages, "\n")
	} else {
		message = "No merchant found."
	}

	return message
}
