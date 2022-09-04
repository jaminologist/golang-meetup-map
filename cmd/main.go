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

	headers := records[0]
	data := make([]map[string]string, 0)
	for i, row := range records[1:] {
		data = append(data, make(map[string]string))
		for j, cell := range row {
			data[i][headers[j]] = cell
		}
	}

	meetups := convertDataToMeetups(data)
	parse("./template/index.html", meetups)
}

var (
	headerName      = "Name"
	headerDate      = "Date"
	headerIcon      = "Icon"
	headerLink      = "Link"
	headerLatitude  = "Latitude"
	headerLongitude = "Longitude"
)

type MeetupMapPage struct {
	Meetups []Meetup
}

type Meetup struct {
	Name      string
	Date      string
	Icon      string
	Link      string
	Latitude  string
	Longitude string
}

func convertDataToMeetups(records []map[string]string) MeetupMapPage {
	meetups := make([]Meetup, 0)
	for _, record := range records {

		iconURL := "./icons/" + record["Icon"]
		if _, err := os.Stat("../docs/" + iconURL); err != nil {
			iconURL = "./icons/Go-Logo_Blue.png"
		}

		meetup := Meetup{
			Name:      record[headerName],
			Icon:      iconURL,
			Date:      record[headerDate],
			Link:      record[headerLink],
			Latitude:  record[headerLatitude],
			Longitude: record[headerLongitude],
		}
		meetups = append(meetups, meetup)
	}
	return MeetupMapPage{
		Meetups: meetups,
	}
}

func parse(path string, meetupMapPage MeetupMapPage) {
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create("../docs/index_test.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, meetupMapPage)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}
