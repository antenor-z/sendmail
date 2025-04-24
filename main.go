package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"
	"sendmail/config"
	"sendmail/util"
	"strings"
	"time"
)

type Email struct {
	From    string
	To      string
	Subject string
	Text    string
}

func sendMail(email Email, config config.Config) {
	form := url.Values{}
	form.Add("from", email.From)
	form.Add("to", email.To)
	form.Add("subject", email.Subject)
	form.Add("html", email.Text)
	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/"+config.Mailgun.Domain+"/messages", strings.NewReader(form.Encode()))
	util.Check(err)

	req.SetBasicAuth("api", config.Mailgun.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	util.Check(err)

	if res.StatusCode != 200 {
		panic("Failed to send email. " + res.Status)
	}
}

func main() {
	config := config.GetConfig()

	var name, email, subject, content, filePath string
	flag.StringVar(&name, "n", "", "name of the recipient")
	flag.StringVar(&email, "e", "", "email of the recipient")
	flag.StringVar(&subject, "s", "", "email subject")
	flag.StringVar(&content, "c", "", "email content")
	flag.StringVar(&filePath, "f", "", "[optional] email content from file. Will override -c")
	flag.Parse()

	if name == "" || email == "" || subject == "" || content == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	to, err := util.FormatTo(name, email)
	util.Check(err)

	if filePath != "" {
		dat, err := os.ReadFile(filePath)
		util.Check(err)
		content = string(dat)
	}
	sendMail(Email{From: config.Mailgun.Sender, To: to, Subject: subject, Text: content}, config)
}
