package csvconvert

import (
	"fmt"
	io "io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type meetupTest struct {
	suite.Suite

	mockParser *mockDataParser

	reader      io.Reader
	testMeetup1 Meetup
	testMeetup2 Meetup
	icons       map[string]bool
	headers     []string
	firstRow    []string
}

func (mt *meetupTest) SetupTest() {
	mt.mockParser = new(mockDataParser)

	mt.testMeetup1 = Meetup{
		Name:      "Name",
		Date:      "2022-01-29",
		Icon:      "Icon",
		Link:      "https://www.google.com",
		Latitude:  "24.9948056",
		Longitude: "-71.0351806",
	}

	mt.testMeetup2 = Meetup{
		Name:      "Name2",
		Date:      "2022-01-30",
		Icon:      "Icon2",
		Link:      "https://www.bing.com",
		Latitude:  "-25.9948056",
		Longitude: "70.0351806",
	}

	mt.reader = meetupsToCSVReader(mt.testMeetup1, mt.testMeetup2)

	mt.icons = map[string]bool{
		mt.testMeetup1.Icon: true,
		mt.testMeetup2.Icon: true,
	}
	mt.headers = []string{
		headerName,
		headerDate,
		headerIcon,
		headerLink,
		headerLatitude,
		headerLongitude,
	}
	mt.firstRow = []string{
		"DURIAN",
		"2006-01-02",
		mt.testMeetup1.Icon,
		"link",
		"1.0",
		"1.0",
	}
}

func Test_workflowTest(t *testing.T) {
	suite.Run(t, new(meetupTest))
}

func (mt *meetupTest) Test_ReadMeetups_newReader_Error() {
	mt.mockParser.On("newReader", mt.reader).Return(nil, fmt.Errorf("new reader error"))

	mu, err := ReadMeetups(mt.mockParser, mt.reader, nil)
	mt.Nil(mu)
	mt.EqualError(err, "failed to read csv records: new reader error")
}

func (mt *meetupTest) Test_ReadMeetups_newReader_ZeroRecords() {
	mt.mockParser.On("newReader", mt.reader).Return([][]string{}, nil)

	mu, err := ReadMeetups(mt.mockParser, mt.reader, nil)
	mt.Nil(mu)
	mt.EqualError(err, "length of csv file is < 1")
}

func (mt *meetupTest) Test_validateLongitude_Error() {
	mt.mockParser.On("parseFloat", "EPIC", 64).Return(-1.0, fmt.Errorf("parse error"))

	err := validateLongitude(mt.mockParser, "EPIC")
	mt.EqualError(err, "unable to validate longitude=EPIC err=parse error")
}

func (mt *meetupTest) Test_validateLatitude_Error() {
	mt.mockParser.On("parseFloat", "EPIC", 64).Return(-1.0, fmt.Errorf("parse error"))

	err := validateLatitude(mt.mockParser, "EPIC")
	mt.EqualError(err, "unable to validate latitude=EPIC err=parse error")
}

func (mt *meetupTest) Test_validateLink_Error() {
	mt.mockParser.On("parseRequestURI", "EPIC").Return(nil, fmt.Errorf("parse url error"))
	err := validateLink(mt.mockParser, "EPIC")
	mt.EqualError(err, "unable to validate url=EPIC, error=parse url error")
}

func (mt *meetupTest) Test_ReadMeetups_validateIcon_Error() {
	mt.mockParser.On("newReader", mt.reader).Return([][]string{mt.headers, mt.firstRow}, nil)
	mt.mockParser.On("parseRequestURI", "link").Return(nil, fmt.Errorf("parse url error"))

	got, err := ReadMeetups(mt.mockParser, mt.reader, mt.icons)
	mt.Nil(got)
	mt.EqualError(err, "failed convert rows into Meetups: row 1: unable to validate url=link, error=parse url error")
}

func (mt *meetupTest) Test_ReadMeetups_validateLatitude_Error() {
	mt.mockParser.On("newReader", mt.reader).Return([][]string{mt.headers, mt.firstRow}, nil)
	mt.mockParser.On("parseRequestURI", "link").Return(nil, nil)
	mt.mockParser.On("parseFloat", "1.0", 64).Return(-1.0, fmt.Errorf("parse error"))

	got, err := ReadMeetups(mt.mockParser, mt.reader, mt.icons)
	mt.Nil(got)
	mt.EqualError(err, "failed convert rows into Meetups: row 1: unable to validate latitude=1.0 err=parse error")
}

func (mt *meetupTest) Test_ReadMeetups_validateLongitude_Error() {
	mt.mockParser.On("newReader", mt.reader).Return([][]string{mt.headers, mt.firstRow}, nil)
	mt.mockParser.On("parseRequestURI", "link").Return(nil, nil)
	mt.mockParser.On("parseFloat", "1.0", 64).Return(1.0, nil).Once()
	mt.mockParser.On("parseFloat", "1.0", 64).Return(-1.0, fmt.Errorf("parse error")).Once()

	got, err := ReadMeetups(mt.mockParser, mt.reader, mt.icons)
	mt.Nil(got)
	mt.EqualError(err, "failed convert rows into Meetups: row 1: unable to validate longitude=1.0 err=parse error")
}
