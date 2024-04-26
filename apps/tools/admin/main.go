package main

import (
	"log"

	"github.com/wscnd/go-service-boilerplate/apps/tools/admin/commands"
)

func main() {
	err := commands.GenKey()
	if err != nil {
		log.Fatalln(err)
	}
}
