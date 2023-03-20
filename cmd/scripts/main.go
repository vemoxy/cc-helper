package main

import (
	"log"
	"os"
	"strings"
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
	case "checkMccV3":
		if len(args[1:]) == 0 {
			log.Print("checkMccV3 needs merchant name as argument")
			return
		}
		checkMccResult := CheckMcc(strings.Join(args[1:], " "))
		log.Printf("check mcc result: %+v", checkMccResult)
	default:
		log.Print("Command not recognized")
	}
}
