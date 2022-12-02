package mail

import (
	"fmt"

	"github.com/ilhamfzri/pendek.in/config"
	"gopkg.in/gomail.v2"
)

type MailClient struct {
	Dialer     *gomail.Dialer
	SenderName string
}

var subjectVerificationEmail = "Verify code to activate your pendek.in account"

func NewMailClient(cfg config.MailConfig) *MailClient {
	fmt.Println(cfg)
	dialer := gomail.NewDialer(
		cfg.StmpHost,
		cfg.StmpPort,
		cfg.AuthEmail,
		cfg.AuthPassword,
	)

	return &MailClient{
		Dialer:     dialer,
		SenderName: cfg.SenderName,
	}
}

func (client *MailClient) SendVerificationEmail(email string, code string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", client.SenderName)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", subjectVerificationEmail)
	mailer.SetBody("text/html", fmt.Sprintf("Your email verification code : %s", code))

	err := client.Dialer.DialAndSend(mailer)
	return err
}
