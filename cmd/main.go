package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path"
	"text/template"
)

//go:embed meetups.csv
var b []byte

var root string

func init() {
	flag.StringVar(&root, "root", ".", "root directory of the project")
}

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

	meetups := convertDataToMeetups(root, data)
	parse(root, root+"/cmd/template/index.html", meetups)
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

func convertDataToMeetups(root string, records []map[string]string) MeetupMapPage {
	meetups := make([]Meetup, 0)
	for _, record := range records {

		iconURL := path.Join("icons", record["Icon"])
		if _, err := os.Stat(path.Join(root, "docs", iconURL)); err != nil {
			iconURL = path.Join("icons", "Go-Logo_Blue.png")
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

func parse(root string, templatePath string, meetupMapPage MeetupMapPage) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(path.Join(root, "docs", "index.html"))
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
