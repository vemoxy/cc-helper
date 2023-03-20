package model

type MerchantChannel struct {
	Id           int    `db:"id"`
	MerchantName string `db:"merchantName"`
	ChannelName  string `db:"channelName"`
	Mcc          int    `db:"mcc"`
	Source       string `db:"source"`
	UpdateTime   int    `db:"updateTime"`
	UpdateBy     string `db:"updateBy"`
}

var MerchantChannelFields []string = []string{
	"id",
	"merchantName",
	"channelName",
	"mcc",
	"source",
	"updateTime",
	"updateBy",
}
