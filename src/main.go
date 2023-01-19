package main

import (
	"log"

	"github.com/yazdsr/master-api/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
