package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"log"
	"os"
	"text/template"
)

//go:embed meetups.csv
var b []byte

func main() {
	csvReader := csv.NewReader(bytes.NewReader(b))
	records, _ := csvReader.ReadAll()
	meetups := recordsToMeetups(records[1:])
	parse("./template/index.html", meetups)
}

type Meetup struct {
	Name      string
	Icon      string
	Link      string
	Latitude  string
	Longitude string
}

func recordsToMeetups(records [][]string) []Meetup {
	meetups := make([]Meetup, 0)
	for _, record := range records {

		iconURL := "./icons/" + record[1]
		if _, err := os.Stat("../docs/" + iconURL); err != nil {
			iconURL = "./icons/Go-Logo_Blue.png"
		}

		meetup := Meetup{
			Name:      record[0],
			Icon:      iconURL,
			Link:      record[2],
			Latitude:  record[3],
			Longitude: record[4],
		}
		meetups = append(meetups, meetup)
	}
	return meetups
}

func parse(path string, meetups []Meetup) {
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create("../docs/index.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, meetups)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}
