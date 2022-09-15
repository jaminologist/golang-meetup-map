package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"jaminologist/golangmeetupmap/internal/archiver"
	"jaminologist/golangmeetupmap/internal/csvconvert"
	"log"
	"os"
	"path"
	"time"
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

	upcomingMeetups, _ := archiver.ArchiveMeetups(time.Now(), meetups)

	var upcomingRows [][]string
	if upcomingMeetups != nil {
		upcomingRows = csvconvert.ConvertMeetupsToRows(upcomingMeetups)
	}

	var upcomingBuf bytes.Buffer
	writer := csv.NewWriter(&upcomingBuf)
	writer.WriteAll(upcomingRows)
	writer.Flush()

	isMeetupsSame, err := isFileTheSame(upcomingBuf.Bytes(), path.Join(root, "docs", "meetups.csv"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("unable to open meetups.csv %v", err)
	}
	if !isMeetupsSame {
		file, err := os.Create(path.Join(root, "docs", "meetups.csv"))
		defer file.Close()
		if err != nil {
			log.Fatalf("unable to create meetups.csv %v", err)
		}
		file.Write(upcomingBuf.Bytes())
	}
}

func isFileTheSame(compareBytes []byte, path string) (bool, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return false, err
	}
	var fileBuf []byte
	file.Read(fileBuf)
	return bytes.Equal(fileBuf, compareBytes), nil
}
