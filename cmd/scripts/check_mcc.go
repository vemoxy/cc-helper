package main

import (
	"github.com/vemoxy/cc-helper/handler"
	"github.com/vemoxy/cc-helper/model"
)

func CheckMcc(merchantName string) []handler.CheckMccResult {
	return handler.CheckMcc(merchantName)
}

func CheckMccV3(merchantName string) []model.MerchantChannel {
	return handler.CheckMccV3(merchantName)
}
