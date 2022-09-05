package main

import (
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

	// Read icons directly and create map of saved icons
	files, err := os.ReadDir((path.Join(root, "docs", "icons")))
	if err != nil {
		log.Fatal(err)
	}

	icons := make(map[string]bool)
	for _, file := range files {
		icons[file.Name()] = true
	}

	// Read meetups.csv
	meetupsCSV, err := os.Open(path.Join(root, "docs", "meetups.csv"))
	defer meetupsCSV.Close()
	if err != nil {
		log.Fatalf("failed to open meetups.csv: %v", err)
	}

	meetups, err := csvconvert.ReadMeetups(meetupsCSV, icons)
	if err != nil {
		log.Fatalf("failed to read meetups.csv: %v", err)
	}
	if err := templater.Parse(root, path.Join(root, "docs", "templates", "index.html"), meetups); err != nil {
		log.Fatalf("failed to create template: %v", err)
	}
}
