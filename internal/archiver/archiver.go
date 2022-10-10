package archiver

import (
	"jaminologist/golangmeetupmap/internal/model"
	"time"
)

// ArchiveMeetups splits a list of meetups into ones that are before and after the given date
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
