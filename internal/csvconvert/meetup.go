package csvconvert

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"jaminologist/golangmeetupmap/internal/model"
	"net/url"
	"strconv"
	"time"
)

const dateLayoutYYYYMMDD = "2006-01-02"

var (
	headerName      = "Name"
	headerDate      = "Date"
	headerIcon      = "Icon"
	headerLink      = "Link"
	headerLatitude  = "Latitude"
	headerLongitude = "Longitude"
)

func ReadMeetups(csvReader io.Reader, icons map[string]bool) ([]model.Meetup, error) {
	records, err := csv.NewReader(csvReader).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv records: %w", err)
	}

	if len(records) < 1 {
		return nil, errors.New("length of csv file is < 1")
	}

	headers := records[0]
	rows := records[1:]

	meetups, err := convertRowsToMeetups(headers, rows, icons)
	if err != nil {
		return nil, fmt.Errorf("failed convert rows into Meetups: %w", err)
	}
	return meetups, nil
}

func convertRowsToMeetups(headers []string, rows [][]string, icons map[string]bool) ([]model.Meetup, error) {
	mappedRows := make([]map[string]string, 0)
	for i, row := range rows {
		mappedRows = append(mappedRows, make(map[string]string))
		for j, cell := range row {
			mappedRows[i][headers[j]] = cell
		}
	}

	meetups := make([]model.Meetup, 0)
	for i, row := range mappedRows {
		meetup := model.Meetup{
			Name:      row[headerName],
			Icon:      row[headerIcon],
			Link:      row[headerLink],
			Latitude:  row[headerLatitude],
			Longitude: row[headerLongitude],
		}

		index := i + 1 // Header row is deleted, for correct row number you need to add 1.
		if err := validateName(meetup.Name); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		date, err := time.Parse(dateLayoutYYYYMMDD, row[headerDate])
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}
		meetup.Date = date

		if err := validateIcon(meetup.Icon, icons); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateLink(meetup.Link); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateLatitude(meetup.Latitude); err != nil {
			return nil, fmt.Errorf("row %d: %w", index, err)
		}

		if err := validateLongitude(meetup.Longitude); err != nil {
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
	if _, err := time.Parse(dateLayoutYYYYMMDD, date); err != nil {
		return fmt.Errorf("unable to validate date: %s, error: %w", date, err)
	}
	return nil
}

func validateIcon(icon string, icons map[string]bool) error {
	if _, ok := icons[icon]; !ok {
		if icon != "" {
			return fmt.Errorf("unable to validate icon %s, it is present in docs/icons", icon)
		}
	}
	return nil
}

func validateLink(link string) error {
	if _, err := url.ParseRequestURI(link); err != nil {
		return fmt.Errorf("unable to validate url: %s, error: %w", link, err)
	}
	return nil
}

func validateLatitude(latitude string) error {
	if _, err := strconv.ParseFloat(latitude, 64); err != nil {
		return fmt.Errorf("unable to validate latitude: %s, error: %w", latitude, err)
	}
	return nil
}

func validateLongitude(longitude string) error {
	if _, err := strconv.ParseFloat(longitude, 64); err != nil {
		return fmt.Errorf("unable to validate longitude %s, %w", longitude, err)
	}
	return nil
}

func ConvertMeetupsToRows(meetups []model.Meetup) (rows [][]string) {
	rows = append(rows, []string{headerName, headerDate, headerIcon, headerLink, headerLatitude, headerLongitude})
	for _, meetup := range meetups {
		rows = append(rows,
			[]string{
				meetup.Name,
				meetup.Date.Format(dateLayoutYYYYMMDD),
				meetup.Icon,
				meetup.Link,
				meetup.Latitude,
				meetup.Longitude,
			},
		)
	}
	return rows
}
