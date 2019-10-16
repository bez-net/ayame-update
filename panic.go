package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

const cobot_dbs = "https://hooks.slack.com/services/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"

func getStringPanicStack() (str string) {
	if r := recover(); r != nil {
		str = string(debug.Stack())
	}
	return
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendSlackNotification(webhookUrl string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}

	return nil
}

func sendDebugStackToSlack() {
	str := getStringPanicStack()
	log.Printf(str)
	SendSlackNotification(cobot_dbs, str)
}
