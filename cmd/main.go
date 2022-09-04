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
	Link      string
	Latitude  string
	Longitude string
}

func recordsToMeetups(records [][]string) []Meetup {
	meetups := make([]Meetup, 0)
	for _, record := range records {
		meetup := Meetup{
			Name:      record[0],
			Link:      record[1],
			Latitude:  record[2],
			Longitude: record[3],
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
