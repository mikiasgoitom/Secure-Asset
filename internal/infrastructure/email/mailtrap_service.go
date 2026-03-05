package email

import (
	"fmt"
	"net/smtp"
)

// MailtrapService implements IEmailService using Mailtrap SMTP.
type MailtrapService struct {
	host   string
	port   string
	user   string
	pass   string
	from   string
}

func NewMailtrapService(host, port, user, pass, from string) *MailtrapService {
	return &MailtrapService{host: host, port: port, user: user, pass: pass, from: from}
}

// SendOTP sends a simple OTP email.
func (m *MailtrapService) SendOTP(toEmail string, otpCode string) error {
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)
	recipient := []string{toEmail}
	msg := []byte(fmt.Sprintf("Subject: Your Secure-Asset OTP\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n\r\n"+
		"Your one-time password is: %s. It expires in 10 minutes.\r\n", m.from, toEmail, otpCode))
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, recipient, msg)
}
