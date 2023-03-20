package handler

import (
	"fmt"
	"strings"

	"github.com/vemoxy/cc-helper/data"
	"github.com/vemoxy/cc-helper/model"
)

func CheckMcc(merchantName string) []model.MerchantChannel {
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
