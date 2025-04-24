package util

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"
)

func Check(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
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
