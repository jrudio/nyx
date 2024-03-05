package main

import (
	"flag"
	"fmt"
	"github.com/jrudio/nyx"
	"log"
	"os"
)

func main() {
	var lineup nyx.Lineup
	inputFile := flag.String("i", "", "path to .html file")
	lineupURL := flag.String("u", "", "URL from the Insomniac app")

	flag.Parse()

	if *inputFile == "" && *lineupURL == "" {
		fmt.Println("html file or a url is required to parse. please use the -i or -u flag")

		return
	}

	if *inputFile != "" {
		// read a file
		data, err := os.ReadFile(*inputFile)

		if err != nil {
			log.Fatal(err)
		}

		lineup, err = nyx.ParseLineup(data)

		if err != nil {
			log.Fatalf("failed parsing the lineup: %v", err)
		}

		printLineup(lineup)

	} else if *lineupURL != "" {
		var err error

		lineup, err = nyx.Get(*lineupURL)

		if err != nil {
			log.Fatalf("failed parsing the lineup: %v", err)
		}

		printLineup(lineup)
	} else {
		log.Fatalf("something went wrong...")
	}
}

func printLineup(lineup nyx.Lineup) {
		if lineup.Size < 1 {
		log.Println("did not find any artists at the given url or html file")

		return
	}

	log.Println("found the following lineup...")

	for i := 0; i < lineup.Size; i++ {
		artist := lineup.Artists[i]

		log.Println(artist.Name, artist.Img)
	}
}
