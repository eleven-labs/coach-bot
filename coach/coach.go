package coach

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eko/slackbot"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"

	"../config"
)

var (
	spreadsheetID = config.Getenv("ELEVENBOT_COACH_SPREADSHEET_ID")
	sheetRange    = "Planning Officiel!A1:AX"

	columnLastname                  = 0
	columnFirstname                 = 1
	columnSlack                     = 2
	columnEntryDate                 = 3
	columnRecruiter                 = 4
	columnCommercial                = 5
	columnManager                   = 6
	columnMeetingRecruiter          = 7
	columnMeetingTechnicalFirstname = 8
	columnMeetingTechnicalSlack     = 9

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

// NotifyPassedMeetings reads spreadsheet data and prepare messages to be sent.
func NotifyPassedMeetings(sheetsService *sheets.Service) {
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, sheetRange).Do()

	if err != nil {
		log.Fatalf("[coach] Unable to retrieve data from sheet. %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Print("[coach] No data found.")
		return
	}

	var remindedColumnKey int

	for key, row := range resp.Values {
		if 0 == key {
			remindedColumnKey = FindRemindedColumn(row)
		} else {
			information := row[remindedColumnKey].(string)

			if strings.Contains(information, "Bilan OKR") {
				fmt.Print(information)
				SendSlackNotification(row)
			}
		}
	}
}

// FindRemindedColumn returns the reminded column at current date (current month -2 months)
func FindRemindedColumn(row []interface{}) int {
	var value int

	now := time.Now()
	reminder := now.AddDate(0, -2, 0) // reminder set for 2 after meeting

	for key, column := range row {
		monthInt, _ := strconv.Atoi(reminder.Format("01"))
		monthFrench = monthMapping[monthInt]

		if column == fmt.Sprintf("%s %s", monthFrench, reminder.Format("2006")) {
			value = key
			break
		}
	}

	return value
}

// SendSlackNotification sens notification to concerned user.
func SendSlackNotification(row []interface{}) {
	var userId string

	usersResponse, _ := slackbot.ListUsers()

	for i := 0; i < len(usersResponse.Members); i++ {
		if row[columnMeetingTechnicalSlack] == usersResponse.Members[i].Name {
			userId = usersResponse.Members[i].ID
			break
		}
	}

	channel := slackbot.Channel{User: userId}
	slackbot.OpenIM(channel)

	log.WithFields(log.Fields{
		"date":      time.Now().Format("2006-01-02"),
		"coach":     row[columnMeetingTechnicalSlack],
		"concerned": fmt.Sprintf("%s %s", row[columnFirstname], row[columnLastname]),
	}).Info("Slack message sent.")

	message := slackbot.Message{
		AsUser:  true,
		Channel: userId,
		Text: fmt.Sprintf(
			"Hello %s !\nEn %s tu as fait passer (avec %s) le bilan OKR de %s.\nIl est temps de prendre de ses nouvelles :g11rocket:",
			row[columnMeetingTechnicalFirstname], monthFrench, row[columnMeetingRecruiter], fmt.Sprintf("%s %s", row[columnFirstname], row[columnLastname]),
		),
	}

	slackbot.PostMessage(message)
}
