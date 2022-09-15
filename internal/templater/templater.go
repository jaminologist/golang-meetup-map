package templater

import (
	"fmt"
	"html/template"
	"jaminologist/golangmeetupmap/internal/model"
	"os"
	"path"
)

type MeetupMapPage struct {
	Meetups []model.Meetup
}

func Parse(root string, templatePath string, meetups []model.Meetup) error {
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
