package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatal("SLACK_TOKEN must be set")
	}
	channel, ok := os.LookupEnv("SLACK_CHANNEL")
	if !ok {
		log.Fatal("SLACK_CHANNEL must be set")
	}

	api := slack.New(token)

	lastTick := time.Now().Unix()

	ticker := time.NewTicker(15 * time.Minute)

	log.Info("starting loop")

	errorCount := 0
	for {
		if errorCount >= 3 {
			log.Fatal(fmt.Errorf("error threshold exceeded: %d", errorCount))
		}

		tick := <-ticker.C

		log.Info("tick")

		users, err := api.GetUsers()
		if err != nil {
			errorCount++
			log.Error(err)
			continue
		}

		var deleted []string
		for _, u := range users {
			if u.Deleted && (u.Updated.Time().Unix() > lastTick) {
				s := fmt.Sprintf("%s - %s - @%s\n", u.Profile.RealName, u.Profile.Title, u.Profile.DisplayNameNormalized)
				deleted = append(deleted, s)
			}
		}

		l := len(deleted)
		log.Infof("%d new", l)
		if l > 0 {
			if _, _, err = api.PostMessage(channel, slack.MsgOptionText(strings.Join(deleted, "\n"), false)); err != nil {
				errorCount++
				log.Error(err)
				continue
			}
		}

		errorCount = 0
		lastTick = tick.Unix()
	}
}
