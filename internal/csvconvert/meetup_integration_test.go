package csvconvert_test

import (
	"jaminologist/golangmeetupmap/internal/csvconvert"
	"os"
	"path"
	"testing"
)

func TestReadMeetups(t *testing.T) {
	t.Run("validate meetups.csv", func(t *testing.T) {
		root := "../../"

		files, err := os.ReadDir((path.Join(root, "docs", "icons")))
		if err != nil {
			t.Errorf("failed to read docs/icons: %v", err)
		}

		icons := make(map[string]bool)
		for _, file := range files {
			icons[file.Name()] = true
		}

		meetupsCSV, err := os.Open(path.Join(root, "docs", "meetups.csv"))
		defer meetupsCSV.Close()
		if err != nil {
			t.Errorf("failed to open meetups.csv: %v", err)
		}

		_, err = csvconvert.ReadMeetups(&csvconvert.DataParse{}, meetupsCSV, icons)
		if err != nil {
			t.Errorf("failed to read meetups.csv %v", err)
		}
	})
}
