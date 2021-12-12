// Package email handles the sending of emails
package email

import (
	"bytes"
	"fmt"
	"github.com/uberswe/golang-base-project/config"
	"github.com/uberswe/golang-base-project/util"
	"log"
	"mime/multipart"
	"net/smtp"
	"strings"
)

type Service struct {
	Config config.Config
}

func New(c config.Config) Service {
	return Service{
		Config: c,
	}
}

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
	htmlMessage := util.StringLinkToHTMLLink(message)
	htmlMessage = util.NL2BR(htmlMessage)
	_, _ = fmt.Fprintf(&b, "\r\n\r\n--%s\r\nContent-Type: %s; charset=UTF-8;\nContent-Transfer-Encoding: 8bit\r\n\r\n", writer.Boundary(), "text/html")
	b.Write([]byte(htmlMessage))

	_, _ = fmt.Fprintf(&b, "\r\n\r\n--%s--\r\n", writer.Boundary())

	sender := s.Config.SMTPSender
	if strings.Contains(sender, "<") {
		sender = util.GetStringBetweenStrings(sender, "<", ">")
	}

	// Sending email.
	err := smtp.SendMail(fmt.Sprintf("%s:%s", s.Config.SMTPHost, s.Config.SMTPPort), auth, sender, []string{to}, b.Bytes())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("Email sent to %s", to))
}
