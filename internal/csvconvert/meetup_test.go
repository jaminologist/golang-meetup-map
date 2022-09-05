package csvconvert

import (
	"reflect"
	"testing"
)

func TestConvertRowsToMeetups(t *testing.T) {
	type args struct {
		headers []string
		rows    [][]string
		icons   map[string]bool
	}

	testMeetup1 := Meetup{
		Name:      "Name",
		Date:      "Date",
		Icon:      "Icon",
		Link:      "Link",
		Latitude:  "Latitude",
		Longitude: "Longitude",
	}

	testMeetup2 := Meetup{
		Name:      "Name2",
		Date:      "Date2",
		Icon:      "Icon2",
		Link:      "Link2",
		Latitude:  "Latitude2",
		Longitude: "Longitude2",
	}

	tests := []struct {
		name string
		args args
		want []Meetup
	}{
		{
			name: "Two meetups should be created with the correct information",
			args: args{
				headers: []string{headerName, headerDate, headerIcon, headerLink, headerLatitude, headerLongitude},
				rows: [][]string{
					{testMeetup1.Name, testMeetup1.Date, testMeetup1.Icon, testMeetup1.Link, testMeetup1.Latitude, testMeetup1.Longitude},
					{testMeetup2.Name, testMeetup2.Date, testMeetup2.Icon, testMeetup2.Link, testMeetup2.Latitude, testMeetup2.Longitude},
				},
				icons: map[string]bool{
					"Icon":  true,
					"Icon2": true,
				},
			},
			want: []Meetup{testMeetup1, testMeetup2},
		},
		{
			name: "icon can not be found. Should be set to default.png",
			args: args{
				headers: []string{headerName, headerDate, headerIcon, headerLink, headerLatitude, headerLongitude},
				rows: [][]string{
					{testMeetup1.Name, testMeetup1.Date, testMeetup1.Icon, testMeetup1.Link, testMeetup1.Latitude, testMeetup1.Longitude},
					{testMeetup2.Name, testMeetup2.Date, testMeetup2.Icon, testMeetup2.Link, testMeetup2.Latitude, testMeetup2.Longitude},
				},
				icons: map[string]bool{
					"Icon": true,
				},
			},
			want: []Meetup{
				testMeetup1,
				{
					testMeetup2.Name,
					testMeetup2.Date,
					"default.png",
					testMeetup2.Link,
					testMeetup2.Latitude,
					testMeetup2.Longitude,
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertRowsToMeetups(tt.args.headers, tt.args.rows, tt.args.icons); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertRowsToMeetups() = %v, want %v", got, tt.want)
			}
		})
	}
}
