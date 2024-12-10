package email

import (
	"gopkg.in/gomail.v2"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
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
func IsVail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if !re.MatchString(email) {
		return false
	}
	at := strings.LastIndex(email, "@")
	if at == -1 {
		return false
	}
	domain := email[at+1:]

	// Perform DNS lookup for MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}
	return true
}
