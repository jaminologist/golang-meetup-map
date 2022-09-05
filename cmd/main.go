package main

import (
	_ "embed"
	"encoding/csv"
	"flag"
	"jaminologist/golangmeetupmap/internal/csvconvert"
	"jaminologist/golangmeetupmap/internal/templater"
	"log"
	"os"
	"path"
)

var root string

func init() {
	flag.StringVar(&root, "root", ".", "root directory of the project")
}

func main() {

	// Read meetups.csv and split into headers and rows.
	meetupsCSV, err := os.Open(path.Join(root, "docs", "meetups.csv"))
	defer meetupsCSV.Close()
	if err != nil {
		log.Fatalf("failed to read meetup.csv: %v", err)
	}
	records, err := csv.NewReader(meetupsCSV).ReadAll()
	if err != nil {
		log.Fatalf("failed to read csv records: %v", err)
	}

	headers := records[0]
	rows := records[1:]

	// Read icons directly and create map of saved icons
	files, err := os.ReadDir((path.Join(root, "docs", "icons")))
	if err != nil {
		log.Fatal(err)
	}

	icons := make(map[string]bool)
	for _, file := range files {
		icons[file.Name()] = true
	}

	meetups := csvconvert.ConvertRowsToMeetups(headers, rows, icons)
	if err := templater.Parse(root, path.Join(root, "docs", "templates", "index.html"), meetups); err != nil {
		log.Fatalf("failed to create template: %v", err)
	}
}
