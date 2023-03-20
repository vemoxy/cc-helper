package main

import (
	"github.com/vemoxy/cc-helper/handler"
	"github.com/vemoxy/cc-helper/model"
)

func CheckMcc(merchantName string) []model.MerchantChannel {
	return handler.CheckMcc(merchantName)
}
