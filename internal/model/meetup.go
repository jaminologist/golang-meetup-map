package model

import "time"

type Meetup struct {
	Name      string
	Date      time.Time
	Icon      string
	Link      string
	Latitude  string
	Longitude string
}
