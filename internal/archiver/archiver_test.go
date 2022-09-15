package archiver

import (
	"jaminologist/golangmeetupmap/internal/model"
	"reflect"
	"testing"
	"time"
)

func TestArchiveMeetups(t *testing.T) {
	meetupJan2022 := model.Meetup{Date: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)}
	meetupFeb2022 := model.Meetup{Date: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)}
	meetupMar2022 := model.Meetup{Date: time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC)}
	type args struct {
		date    time.Time
		meetups []model.Meetup
	}
	tests := []struct {
		name         string
		args         args
		wantUpcoming []model.Meetup
		wantArchived []model.Meetup
	}{
		{
			"when a meetup is out of date, it should be added to the archive",
			args{
				meetupFeb2022.Date, []model.Meetup{meetupJan2022, meetupMar2022},
			},
			[]model.Meetup{meetupMar2022},
			[]model.Meetup{meetupJan2022},
		},
		{
			"when all meetups are out of date, upcoming should by empty",
			args{
				meetupMar2022.Date, []model.Meetup{meetupJan2022, meetupJan2022},
			},
			nil,
			[]model.Meetup{meetupJan2022, meetupJan2022},
		},
		{
			"when all meetups are upcoming, archived should by empty",
			args{
				meetupJan2022.Date, []model.Meetup{meetupMar2022, meetupMar2022},
			},
			[]model.Meetup{meetupMar2022, meetupMar2022},
			nil,
		},
		{
			"when a meetup is the same date as the archive date it shoudl not be archived",
			args{
				meetupJan2022.Date, []model.Meetup{meetupJan2022},
			},
			[]model.Meetup{meetupJan2022},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpcoming, gotArchived := ArchiveMeetups(tt.args.date, tt.args.meetups)
			if !reflect.DeepEqual(gotUpcoming, tt.wantUpcoming) {
				t.Errorf("ArchriveMeetups() gotUpcoming = %v, want %v", gotUpcoming, tt.wantUpcoming)
			}
			if !reflect.DeepEqual(gotArchived, tt.wantArchived) {
				t.Errorf("ArchriveMeetups() gotArchived = %v, want %v", gotArchived, tt.wantArchived)
			}
		})
	}
}
