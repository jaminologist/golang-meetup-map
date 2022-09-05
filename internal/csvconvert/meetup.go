package csvconvert

var (
	headerName      = "Name"
	headerDate      = "Date"
	headerIcon      = "Icon"
	headerLink      = "Link"
	headerLatitude  = "Latitude"
	headerLongitude = "Longitude"
)

type Meetup struct {
	Name      string
	Date      string
	Icon      string
	Link      string
	Latitude  string
	Longitude string
}

func ConvertRowsToMeetups(headers []string, rows [][]string, icons map[string]bool) []Meetup {
	mappedRows := make([]map[string]string, 0)
	for i, row := range rows {
		mappedRows = append(mappedRows, make(map[string]string))
		for j, cell := range row {
			mappedRows[i][headers[j]] = cell
		}
	}

	meetups := make([]Meetup, 0)
	for _, row := range mappedRows {

		// Sets icon to default, if icon can not be found.
		icon := row[headerIcon]
		if _, ok := icons[icon]; !ok {
			icon = "default.png"
		}

		meetup := Meetup{
			Name:      row[headerName],
			Icon:      icon,
			Date:      row[headerDate],
			Link:      row[headerLink],
			Latitude:  row[headerLatitude],
			Longitude: row[headerLongitude],
		}
		meetups = append(meetups, meetup)
	}
	return meetups
}
