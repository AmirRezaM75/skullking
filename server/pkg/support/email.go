package support

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {
	To      []string
	Subject string
	Body    string
}

func (mail Mail) Send() {
	var auth smtp.Auth

	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")

	if os.Getenv("APP_ENV") == "local" {
		auth = smtp.CRAMMD5Auth(username, password)
	} else {
		auth = smtp.PlainAuth(
			"",
			username,
			password,
			host,
		)
	}

	address := fmt.Sprintf("%s:%s", host, os.Getenv("MAIL_PORT"))

	body := []byte(mail.buildMessage())

	err := smtp.SendMail(address, auth, os.Getenv("MAIL_FROM_ADDRESS"), mail.To, body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email Sent Successfully!")
}

func (mail Mail) buildMessage() string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", os.Getenv("MAIL_FROM_ADDRESS"))
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
