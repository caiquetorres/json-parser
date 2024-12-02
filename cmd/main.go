package main

import (
	"fmt"
	"log"
	"os"

	jsonparser "github.com/caiquetorres/json-parser/parser"
)

func main() {
	file, err := os.Open("test.json")
	if err != nil {
		log.Fatal(err)
	}
	err = jsonparser.Parse(file)
	if err != nil {
		fmt.Println("invalid json")
		log.Fatal(err)
	} else {
		fmt.Println("valid json")
	}
}
