package data

import (
	"context"
	"log"
	"strings"
	"time"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/vemoxy/cc-helper/config"
	"github.com/vemoxy/cc-helper/model"
)

var cachedMerchantChannels []model.MerchantChannel
var cacheLastUpdated time.Time = time.Unix(0, 0)
var refreshInterval time.Duration = 120 * time.Second

func ReloadCache() {
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

	merchantChannelStore := freedb.NewGoogleSheetRowStore(
		auth,
		conf.FreeDB.SpreadsheetId,
		conf.FreeDB.Sheets.MerchantChannel,
		freedb.GoogleSheetRowStoreConfig{Columns: model.MerchantChannelFields},
	)

	var queryOutput []model.MerchantChannel

	err = merchantChannelStore.Select(&queryOutput).
		Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	cachedMerchantChannels = queryOutput
	cacheLastUpdated = time.Now()
}

func getMerchantChannels() []model.MerchantChannel {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Printf("Crashed, err: %v", r)
	// 	}
	// }()

	// if time.Since(cacheLastUpdated) >= refreshInterval {
	// 	ReloadCache()
	// }
	return cachedMerchantChannels
}

func QueryMerchantChannelsByMerchantName(merchantName string) []model.MerchantChannel {
	allMerchantChannels := getMerchantChannels()
	var result []model.MerchantChannel

	merchantName = strings.ToLower(merchantName)

	for _, mc := range allMerchantChannels {
		if strings.Contains(strings.ToLower(mc.MerchantName), merchantName) {
			result = append(result, mc)
		} else if fuzzy.Match(merchantName, strings.ToLower(mc.MerchantName)) {
			result = append(result, mc)
		}
	}

	return result
}
