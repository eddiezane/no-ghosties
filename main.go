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

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	unix := today.Unix()

	api := slack.New(token)
	users, err := api.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	deleted := make([]string, 0)
	for _, u := range users {
		if u.Deleted && (u.Updated.Time().Unix() > unix) {
			s := fmt.Sprintf("%s - %s - @%s\n", u.Profile.RealName, u.Profile.Title, u.Profile.DisplayNameNormalized)
			deleted = append(deleted, s)
		}
	}

	_, _, err = api.PostMessage(channel, slack.MsgOptionText(strings.Join(deleted, "\n"), false))
	if err != nil {
		log.Fatal(err)
	}
}
