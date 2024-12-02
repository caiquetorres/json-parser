package main

import (
	"log"
	"os"

	jsonparser "github.com/caiquetorres/json-parser/parser"
)

func main() {
	file, err := os.Open("test.json")
	if err != nil {
		log.Fatal(err)
	}
	err = jsonparser.Run(file)
	if err != nil {
		log.Fatal(err)
	}
}
