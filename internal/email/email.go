package email

import (
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
	"sync"
)

func SendMail(wg *sync.WaitGroup, to string, subject string, content string) error {
	defer wg.Done()
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("SMTP_USERNAME"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
