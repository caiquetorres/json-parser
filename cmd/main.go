package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	jsonparser "github.com/caiquetorres/json-parser/parser"
)

func main() {
	flag.Parse()
	paths := flag.Args()
	type reader struct {
		id string
		r  io.Reader
	}
	var readers []reader
	if len(paths) == 0 {
		readers = append(readers, reader{
			id: "stdin",
			r:  os.Stdin,
		})
	} else {
		for _, path := range paths {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			readers = append(readers, reader{
				id: path,
				r:  file,
			})
		}
	}
	for _, reader := range readers {
		err := jsonparser.Parse(reader.r)
		isValid := err == nil
		if isValid {
			fmt.Println(fmt.Sprintf("%s valid", reader.id))
		} else {
			fmt.Println(fmt.Sprintf("%s invalid", reader.id))
		}
	}
}
