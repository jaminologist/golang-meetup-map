package templater

import (
	"fmt"
	"html/template"
	"jaminologist/golangmeetupmap/internal/model"
	"os"
	"path"
)

type MeetupMapPage struct {
	Meetups []templateMeetup
}

type templateMeetup struct {
	Name      string
	Date      string
	Icon      string
	Link      string
	Latitude  string
	Longitude string
}

func Parse(root string, templatePath string, ms []model.Meetup) error {

	meetups := make([]templateMeetup, 0)
	for _, m := range ms {
		meetups = append(meetups, templateMeetup{
			Name:      m.Name,
			Date:      m.Date.Format("Monday, 2 January 2006"),
			Icon:      m.Icon,
			Link:      m.Link,
			Latitude:  m.Latitude,
			Longitude: m.Longitude,
		})
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(path.Join(root, "docs", "index.html"))
	defer f.Close()
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}

	err = t.Execute(f, MeetupMapPage{Meetups: meetups})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
