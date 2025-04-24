package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Email struct {
	From    string
	To      string
	Subject string
	Text    string
}

type Mailgun struct {
	ApiKey string `toml:"apiKey"`
	Sender string `toml:"sender"`
	Domain string `toml:"domain"`
}

type Config struct {
	Mailgun Mailgun `toml:"mailgun"`
}

func check(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func getConfig() Config {
	var config Config
	dat, err := os.ReadFile("config.toml")
	check(err)
	_, err = toml.Decode(string(dat), &config)
	check(err)
	return config
}

func FormatTo(name, email string) (string, error) {
	parsed, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}

	parts := strings.Split(parsed.Address, "@")
	if len(parts) != 2 {
		return "", errors.New("email missing domain part")
	}

	addr := &mail.Address{Name: name, Address: parsed.Address}
	return addr.String(), nil
}

func sendMail(email Email, config Config) {
	form := url.Values{}
	form.Add("from", email.From)
	form.Add("to", email.To)
	form.Add("subject", email.Subject)
	form.Add("html", email.Text)
	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/"+config.Mailgun.Domain+"/messages", strings.NewReader(form.Encode()))
	check(err)

	req.SetBasicAuth("api", config.Mailgun.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	check(err)

	if res.StatusCode != 200 {
		panic("Failed to send email. " + res.Status)
	}
}

func main() {
	config := getConfig()

	var name, email, subject, content, filePath string
	flag.StringVar(&name, "n", "", "name of the recepent")
	flag.StringVar(&email, "e", "", "email of the recepent")
	flag.StringVar(&subject, "s", "", "email subject")
	flag.StringVar(&content, "c", "", "email content")
	flag.StringVar(&filePath, "f", "", "[optional] email content from file. Will override -c")
	flag.Parse()

	if name == "" || email == "" || subject == "" || content == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	to, err := FormatTo(name, email)
	check(err)

	if filePath != "" {
		dat, err := os.ReadFile(filePath)
		check(err)
		content = string(dat)
	}

	sendMail(Email{From: config.Mailgun.Sender, To: to, Subject: subject, Text: content}, config)
}
