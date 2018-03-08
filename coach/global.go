package coach

import (
	"strconv"
	"time"

	"../config"
)

var (
	spreadsheetID = config.Getenv("ELEVENBOT_COACH_SPREADSHEET_ID")
	sheetRange    = "Feuille 1!A1:AX"

	columnUserLastname       = 0
	columnUserFirstname      = 1
	columnUserSlack          = 2
	columnTechnicalFirstname = 3
	columnTechnicalSlack     = 4
	columnTechnicalEmail     = 5
	columnRecruiterFirstname = 6
	columnRecruiterSlack     = 7
	columnRecruiterEmail     = 8
	columnMeetingMonth       = 9
	columnLastMeetingMonth   = 10

	monthFrench  string
	monthMapping = map[int]string{
		1:  "Janvier",
		2:  "Février",
		3:  "Mars",
		4:  "Avril",
		5:  "Mai",
		6:  "Juin",
		7:  "Juillet",
		8:  "Août",
		9:  "Septembre",
		10: "Octobre",
		11: "Novembre",
		12: "Décembre",
	}
)

// findNextMonth returns the next month label corresponding to the actual one
func findNextMonth() string {
	now := time.Now()

	intMonth, _ := strconv.Atoi(now.Format("01"))

	if intMonth == 12 { // Special case for December month
		intMonth = 0
	}

	return monthMapping[intMonth+1]
}
