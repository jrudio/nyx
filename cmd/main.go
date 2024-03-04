package main

import (
	"flag"
	"fmt"
	"github.com/jrudio/nyx"
	"log"
	"os"
)

func main() {
	inputFile := flag.String("i", "", "path to .html file")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("html file is required to parse. please use the -i flag")

		return
	}

	// read a file
	data, err := os.ReadFile(*inputFile)

	if err != nil {
		log.Fatal(err)
	}

	lineup, err := nyx.ParseLineup(data)

	if err != nil {
		log.Fatalf("failed parsing the lineup: %v", err)
	}

	log.Println("Found the following lineup...")

	for i := 0; i < lineup.Size; i++ {
		artist := lineup.Artists[i]

		log.Println(artist.Name, artist.Img)
	}
}
