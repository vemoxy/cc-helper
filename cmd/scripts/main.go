package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Print("Command not recognized")
		return
	}

	switch args[0] {
	case "test":
		log.Print("Testing")
	case "importMcc":
		ImportMcc()
	case "addMerchant":
		if len(args[1:]) == 0 {
			log.Print("addMerchant needs merchant name as argument")
			return
		}
		AddMerchant(strings.Join(args[1:], " "))
	case "addChannel":
		// vars check
		if len(args[1:]) != 5 {
			log.Print("addChannel needs 5 arguments")
			return
		}

		merchantId, _ := strconv.Atoi(args[1])
		channelName := args[2]
		mcc, _ := strconv.Atoi(args[3])
		source := args[4]
		updateTime := int(time.Now().Unix())
		updateBy := args[5]

		AddChannel(merchantId, channelName, mcc, source, updateTime, updateBy)
	case "checkMcc":
		if len(args[1:]) == 0 {
			log.Print("checkMcc needs merchant name as argument")
			return
		}
		checkMccResult := CheckMcc(strings.Join(args[1:], " "))
		log.Printf("check mcc result: %+v", checkMccResult)
	case "checkMccV3":
		if len(args[1:]) == 0 {
			log.Print("checkMccV3 needs merchant name as argument")
			return
		}
		checkMccResult := CheckMccV3(strings.Join(args[1:], " "))
		log.Printf("check mcc result: %+v", checkMccResult)
	default:
		log.Print("Command not recognized")
	}
}
