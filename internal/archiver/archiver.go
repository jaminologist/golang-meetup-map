package archiver

import (
	"jaminologist/golangmeetupmap/internal/model"
	"time"
)

// AddToArchive returns a list of meetups that have dates after
func ArchiveMeetups(date time.Time, meetups []model.Meetup) (upcoming, archived []model.Meetup) {
	for _, meetup := range meetups {
		if meetup.Date.Before(date) {
			archived = append(archived, meetup)
			continue
		}
		upcoming = append(upcoming, meetup)
	}
	return upcoming, archived
}
