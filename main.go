package main

import (
	"github.com/eko/slackbot"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	"github.com/eleven-labs/coach-bot/coach"
	"github.com/eleven-labs/coach-bot/config"
	"github.com/eleven-labs/coach-bot/google"
)

func main() {
	log.Info("Bot is starting...")
	slackbot.Token = config.Getenv("ELEVENBOT_SLACK_TOKEN")

	sheetsService := google.GetSheetsService()

	c := cron.New()

	// Every 23st day of the month at 09:00am (UTC so 11:00am in France)
	c.AddFunc("0 0 9 23 * *", func() {
		go coach.NotifyMeetingsOnSlack(sheetsService)
		go coach.NotifyCoachsByEmail(sheetsService)
	})

	c.Start()

	select {}
}
