package model

type Merchant struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

var MerchantFields []string = []string{
	"id",
	"name",
}

type Channel struct {
	Id          int    `db:"id"`
	MerchantId  int    `db:"merchantId"`
	ChannelName string `db:"channelName"`
	Mcc         int    `db:"mcc"`
	Source      string `db:"source"`
	UpdateTime  int    `db:"updateTime"`
	UpdateBy    string `db:"updateBy"`
}

var ChannelFields []string = []string{
	"id",
	"merchantId",
	"channelName",
	"mcc",
	"source",
	"updateTime",
	"updateBy",
}

type MerchantCategoryCode struct {
	Mcc         int    `db:"mcc"`
	Description string `db:"description"`
	Metadata    string `db:"metadata"`
}

var MerchantCategoryCodeFields []string = []string{
	"mcc",
	"description",
	"metadata",
}

type channelType string

func (channelType) Default() channelType        { return channelType("default") }
func (channelType) OverTheCounter() channelType { return channelType("overTheCounter") }
func (channelType) InAppPayment() channelType   { return channelType("inAppPayment") }
func (channelType) Custom(s string) channelType { return channelType("custom: " + s) }
