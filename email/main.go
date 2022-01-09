// Package email handles the sending of emails
package email

import (
	"bytes"
	"fmt"
	"github.com/uberswe/golang-base-project/config"
	"github.com/uberswe/golang-base-project/text"
	"log"
	"mime/multipart"
	"net/smtp"
	"strings"
)

// Service holds a golang-base-project config.Config and provides functions to send emails
type Service struct {
	Config config.Config
}

// New takes a golang-base-project config.Config and returns an instance of Service
func New(c config.Config) Service {
	return Service{
		Config: c,
	}
}

// Send sends an email with the provided subject and message to the provided email.
func (s Service) Send(to string, subject string, message string) {
	// Authentication.
	auth := smtp.PlainAuth("", s.Config.SMTPUsername, s.Config.SMTPPassword, s.Config.SMTPHost)

	// RFC #822 Standard
	writer := multipart.NewWriter(bytes.NewBufferString(""))
	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, "From: %s\r\nTo: %s\r\nSubject: %s\r\n", s.Config.SMTPSender, to, subject)
	_, _ = fmt.Fprintf(&b, "MIME-Version: 1.0\r\n")
	_, _ = fmt.Fprintf(&b, "Content-Type: multipart/alternative; charset=\"UTF-8\"; boundary=\"%s\"\r\n", writer.Boundary())
	_, _ = fmt.Fprintf(&b, "\r\n\r\n--%s\r\nContent-Type: %s; charset=UTF-8;\nContent-Transfer-Encoding: 8bit\r\n\r\n", writer.Boundary(), "text/plain")
	b.Write([]byte(message))
	htmlMessage := text.LinkToHTMLLink(message)
	htmlMessage = text.Nl2Br(htmlMessage)
	_, _ = fmt.Fprintf(&b, "\r\n\r\n--%s\r\nContent-Type: %s; charset=UTF-8;\nContent-Transfer-Encoding: 8bit\r\n\r\n", writer.Boundary(), "text/html")
	b.Write([]byte(htmlMessage))

	_, _ = fmt.Fprintf(&b, "\r\n\r\n--%s--\r\n", writer.Boundary())

	sender := s.Config.SMTPSender
	if strings.Contains(sender, "<") {
		sender = text.BetweenStrings(sender, "<", ">")
	}

	// Sending email.
	err := smtp.SendMail(fmt.Sprintf("%s:%s", s.Config.SMTPHost, s.Config.SMTPPort), auth, sender, []string{to}, b.Bytes())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("Email sent to %s", to))
}
