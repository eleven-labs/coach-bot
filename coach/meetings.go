package coach

import (
	"fmt"
	"strings"
	"time"

	"github.com/eko/slackbot"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

// NotifyMeetingsOnSlack reads spreadsheet data and prepare messages to be sent.
func NotifyMeetingsOnSlack(sheetsService *sheets.Service) {
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, sheetRange).Do()

	if err != nil {
		log.Fatalf("[coach - meetings] Unable to retrieve data from sheet. %v", err)
	}

	if len(resp.Values) == 0 {
		log.Infof("[coach - meetings] No data found")
		return
	}

	for key, row := range resp.Values {
		if 0 != key {
			nextMonth := findNextMonth()

			if strings.Contains(row[columnMeetingMonth].(string), nextMonth) {
				sendSlackNotification(row)
			}
		}
	}
}

// SendSlackNotification sens notification to concerned user.
func sendSlackNotification(row []interface{}) {
	var userID string
	var technicalID string
	var recruiterID string

	usersResponse, _ := slackbot.ListUsers()

	for i := 0; i < len(usersResponse.Members); i++ {
		if row[columnUserSlack] == usersResponse.Members[i].Name {
			userID = usersResponse.Members[i].ID
		} else if row[columnTechnicalSlack] == usersResponse.Members[i].Name {
			technicalID = usersResponse.Members[i].ID
		} else if row[columnRecruiterSlack] == usersResponse.Members[i].Name {
			recruiterID = usersResponse.Members[i].ID
		}

		if userID != "" && technicalID != "" && recruiterID != "" {
			break
		}
	}

	instantMessage := slackbot.MPInstantMessage{Users: fmt.Sprintf("%s,%s,%s", userID, technicalID, recruiterID)}
	conversationResponse, _ := slackbot.OpenMPIM(instantMessage)

	log.WithFields(log.Fields{
		"date":      time.Now().Format("2006-01-02"),
		"concerned": row[columnUserSlack],
		"technical": row[columnTechnicalSlack],
		"recruiter": row[columnRecruiterSlack],
	}).Info("Slack message sent.")

	message := slackbot.Message{
		AsUser:  true,
		Channel: conversationResponse.Group.ID,
		Text: fmt.Sprintf(
			"Salut tous les 3,\n\nÀ partir de la semaine prochaine, vous devez vous retrouver pour planifier le bilan OKR de %s.\n\nJe vous laisse le soin de trouver une date à laquelle vous êtes tous les 3 disponibles.\n\nÀ bientôt,\nWilson :male-astronaut:",
			fmt.Sprintf("%s %s", row[columnUserFirstname], row[columnUserLastname]),
		),
	}

	slackbot.PostMessage(message)
}
