package coach

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/mattbaird/gochimp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"

	"github.com/eleven-labs/coach-bot/config"
)

// PlanningEntry defines a planning for a coach.
type PlanningEntry struct {
	Month      string
	CoachName  string
	CoachEmail string
	Subentries []PlanningSubentry
}

// PlanningSubentry defines an entry of a coach planning (for each users).
type PlanningSubentry struct {
	RecruiterName  string
	RecruiterEmail string
	UserName       string
}

// NotifyCoachsByEmail reads spreadsheet data and send a summary of coach meetings for next month.
func NotifyCoachsByEmail(sheetsService *sheets.Service) {
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, sheetRange).Do()

	if err != nil {
		log.Fatalf("[coach] Unable to retrieve data from sheet. %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Print("[coach] No data found.")
		return
	}

	planning := make(map[string]PlanningEntry)

	for key, row := range resp.Values {
		if 0 != key {
			nextMonth := findNextMonth()

			if strings.Contains(row[columnMeetingMonth].(string), nextMonth) {
				_, exists := planning[row[columnTechnicalEmail].(string)]

				subentry := PlanningSubentry{
					RecruiterName:  row[columnRecruiterFirstname].(string),
					RecruiterEmail: row[columnRecruiterEmail].(string),
					UserName:       fmt.Sprintf("%s %s", row[columnUserFirstname], row[columnUserLastname]),
				}

				if !exists {
					entry := PlanningEntry{
						Month:      nextMonth,
						CoachName:  row[columnTechnicalFirstname].(string),
						CoachEmail: row[columnTechnicalEmail].(string),
						Subentries: []PlanningSubentry{subentry},
					}

					planning[row[columnTechnicalEmail].(string)] = entry
				} else {
					subentries := planning[row[columnTechnicalEmail].(string)].Subentries
					subentries = append(subentries, subentry)

					entry := planning[row[columnTechnicalEmail].(string)]
					entry.Subentries = subentries

					planning[row[columnTechnicalEmail].(string)] = entry
				}
			}
		}
	}

	for _, entry := range planning {
		sendEmail(entry)
	}
}

// sendEmail sends an email to a coach with its planning entries for next month.
func sendEmail(entry PlanningEntry) {
	recipients := []gochimp.Recipient{
		gochimp.Recipient{Email: entry.CoachEmail, Name: entry.CoachName, Type: "cc"},
		gochimp.Recipient{Email: "mleifer@eleven-labs.com", Name: "Maxime Leifer", Type: "cc"},
	}

	recruiterAlreadyAdded := make(map[string]bool)
	var userList bytes.Buffer

	for _, subentry := range entry.Subentries {
		if _, exists := recruiterAlreadyAdded[subentry.RecruiterEmail]; !exists {
			recipients = append(recipients, gochimp.Recipient{Email: subentry.RecruiterEmail, Name: subentry.RecruiterName, Type: "cc"})
			recruiterAlreadyAdded[subentry.RecruiterEmail] = true
		}

		userList.WriteString(fmt.Sprintf("<li>Astronaute %s avec %s et %s</li>", subentry.UserName, entry.CoachName, subentry.RecruiterName))
	}

	mandrillApi, _ := gochimp.NewMandrill(config.Getenv("ELEVENBOT_MANDRILL_API_KEY"))

	message := gochimp.Message{
		Html: fmt.Sprintf(
			"Salut %s,<br><br>Voici ton planning de coaching du mois à venir :)<br><br>Bilans à effectuer :<br><br>%s<br><br>Je te conseille d'aller voir dès maintenant avec eux pour prévoir quand faire tout ça !<br><br>À bientôt,<br><br>Wilson",
			entry.CoachName, userList.String(),
		),
		Subject:            fmt.Sprintf("[Coaching] Bilans à venir en %s", entry.Month),
		FromEmail:          "wilson@eleven-labs.com",
		FromName:           "Wilson",
		To:                 recipients,
		PreserveRecipients: true,
	}

	_, err := mandrillApi.MessageSend(message, false)

	if err != nil {
		log.WithFields(log.Fields{
			"date":  time.Now().Format("2006-01-02"),
			"coach": entry.CoachEmail,
		}).Errorf("Email not sent.")
	}

	log.WithFields(log.Fields{
		"date":  time.Now().Format("2006-01-02"),
		"coach": entry.CoachEmail,
	}).Info("Email sent.")
}
