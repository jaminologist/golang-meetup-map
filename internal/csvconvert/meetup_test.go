package csvconvert

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"testing"
)

func TestReadMeetups(t *testing.T) {

	testMeetup1 := Meetup{
		Name:      "Name",
		Date:      "2022-01-29",
		Icon:      "Icon",
		Link:      "https://www.google.com",
		Latitude:  "24.9948056",
		Longitude: "-71.0351806",
	}

	testMeetup2 := Meetup{
		Name:      "Name2",
		Date:      "2022-01-30",
		Icon:      "Icon2",
		Link:      "https://www.bing.com",
		Latitude:  "-25.9948056",
		Longitude: "70.0351806",
	}

	reader := meetupsToCSVReader(testMeetup1, testMeetup2)

	type args struct {
		csvReader io.Reader
		icons     map[string]bool
	}
	tests := []struct {
		name    string
		args    args
		want    []Meetup
		wantErr bool
	}{
		{
			name: "when all meetups are valid, should run successfully",
			args: args{
				csvReader: reader,
				icons: map[string]bool{
					testMeetup1.Icon: true,
					testMeetup2.Icon: true,
				},
			},
			want: []Meetup{
				testMeetup1,
				testMeetup2,
			},
			wantErr: false,
		},
		{
			name: "when name is empty, should return err",
			args: args{
				csvReader: meetupsToCSVReader(Meetup{"", testMeetup1.Date, testMeetup1.Icon, testMeetup1.Link, testMeetup1.Latitude, testMeetup1.Longitude}),
				icons: map[string]bool{
					testMeetup1.Icon: true,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when date is incorrect format, should return err",
			args: args{
				csvReader: meetupsToCSVReader(Meetup{testMeetup1.Name, "29-01-2022", testMeetup1.Icon, testMeetup1.Link, testMeetup1.Latitude, testMeetup1.Longitude}),
				icons: map[string]bool{
					testMeetup1.Icon: true,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when icon is not found, should return err",
			args: args{
				csvReader: meetupsToCSVReader(Meetup{testMeetup1.Name, testMeetup1.Date, testMeetup1.Icon, testMeetup1.Link, testMeetup1.Latitude, testMeetup1.Longitude}),
				icons:     map[string]bool{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when link is not url, should return err",
			args: args{
				csvReader: meetupsToCSVReader(Meetup{testMeetup1.Name, testMeetup1.Date, testMeetup1.Icon, "hyrule", testMeetup1.Latitude, testMeetup1.Longitude}),
				icons:     map[string]bool{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when latitude is not a number, should return err",
			args: args{
				csvReader: meetupsToCSVReader(Meetup{testMeetup1.Name, testMeetup1.Date, testMeetup1.Icon, testMeetup1.Link, "latitude", testMeetup1.Longitude}),
				icons:     map[string]bool{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when longitude is not a number, should return err",
			args: args{
				csvReader: meetupsToCSVReader(Meetup{testMeetup1.Name, testMeetup1.Date, testMeetup1.Icon, testMeetup1.Link, testMeetup1.Latitude, "longitude"}),
				icons:     map[string]bool{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadMeetups(&DataParse{}, tt.args.csvReader, tt.args.icons)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMeetups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadMeetups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func meetupsToCSVReader(meetups ...Meetup) io.Reader {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"Name", "Date", "Icon", "Link", "Latitude", "Longitude"})
	for _, meetup := range meetups {
		writer.Write([]string{meetup.Name, meetup.Date, meetup.Icon, meetup.Link, meetup.Latitude, meetup.Longitude})
	}
	writer.Flush()
	return bytes.NewReader(buf.Bytes())
}
