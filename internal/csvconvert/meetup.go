package csvconvert

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

var (
	headerName      = "Name"
	headerDate      = "Date"
	headerIcon      = "Icon"
	headerLink      = "Link"
	headerLatitude  = "Latitude"
	headerLongitude = "Longitude"
)

type Meetup struct {
	Name      string
	Date      string
	Icon      string
	Link      string
	Latitude  string
	Longitude string
}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name dataParser -inpkg --filename data_parser_mock.go
type dataParser interface {
	parseRequestURI(rawURL string) (*url.URL, error)
	parseFloat(s string, bitSize int) (float64, error)
	newReader(r io.Reader) (records [][]string, err error)
}

// DataParse struct implements dataParser to extend test coverage
type DataParse struct{}

func (*DataParse) parseRequestURI(rawURL string) (*url.URL, error) {
	return url.ParseRequestURI(rawURL)
}

func (*DataParse) parseFloat(s string, bitSize int) (float64, error) {
	return strconv.ParseFloat(s, bitSize)
}

func (*DataParse) newReader(r io.Reader) (records [][]string, err error) {
	return csv.NewReader(r).ReadAll()
}

func ReadMeetups(d dataParser, csvReader io.Reader, icons map[string]bool) ([]Meetup, error) {
	records, err := d.newReader(csvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read csv records: %w", err)
	}

	if len(records) < 1 {
		return nil, errors.New("length of csv file is < 1")
	}

	headers := records[0]
	rows := records[1:]

	meetups, err := convertRowsToMeetups(d, headers, rows, icons)
	if err != nil {
		return nil, fmt.Errorf("failed convert rows into Meetups: %w", err)
	}
	return meetups, nil
}

func convertRowsToMeetups(d dataParser, headers []string, rows [][]string, icons map[string]bool) ([]Meetup, error) {
	mappedRows := make([]map[string]string, 0)
	for i, row := range rows {
		mappedRows = append(mappedRows, make(map[string]string))
		for j, cell := range row {
			mappedRows[i][headers[j]] = cell
		}
	}

	meetups := make([]Meetup, 0)
	for i, row := range mappedRows {
		meetup := Meetup{
			Name:      row[headerName],
			Date:      row[headerDate],
			Icon:      row[headerIcon],
			Link:      row[headerLink],
			Latitude:  row[headerLatitude],
			Longitude: row[headerLongitude],
		}

		index := i + 1 // Header row is deleted, for correct row number you need to add 1.
		if err := validateName(meetup.Name); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateDate(meetup.Date); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateIcon(meetup.Icon, icons); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateLink(d, meetup.Link); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateLatitude(d, meetup.Latitude); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateLongitude(d, meetup.Longitude); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}
		meetups = append(meetups, meetup)
	}
	return meetups, nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("unable to validate date: %s, error: %w", date, err)
	}
	return nil
}

func validateIcon(icon string, icons map[string]bool) error {
	fmt.Println(icons)
	if _, ok := icons[icon]; !ok {
		if icon != "" {
			return fmt.Errorf("unable to validate icon %s, it is present in docs/icons", icon)
		}
	}
	return nil
}

func validateLink(d dataParser, link string) error {
	if _, err := d.parseRequestURI(link); err != nil {
		return fmt.Errorf("unable to validate url=%s, error=%w", link, err)
	}
	return nil
}

func validateLatitude(d dataParser, latitude string) error {
	if _, err := d.parseFloat(latitude, 64); err != nil {
		return fmt.Errorf("unable to validate latitude=%s err=%w", latitude, err)
	}
	return nil
}

func validateLongitude(d dataParser, longitude string) error {
	if _, err := d.parseFloat(longitude, 64); err != nil {
		return fmt.Errorf("unable to validate longitude=%s err=%w", longitude, err)
	}
	return nil
}
