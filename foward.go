package main

import (
	"crypto/tls"
	"io"
	"log"

	"github.com/emersion/go-imap"
	// "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/go-gomail/gomail"
)

func foewardMails(messages []*imap.Message, smtpHost string, smtpPort int, smtpUser, smtpPass, destinationEmail string) error {
	for _, msg := range messages {
		if msg == nil {
			continue
		}
		section := &imap.BodySectionName{}
		r := msg.GetBody(section)
		if r == nil {
			log.Println("server can't fetch message body")
			continue
		}
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Println("can not read the mail", err)
			continue
		}
		m := gomail.NewMessage()
		m.SetHeader("From", smtpUser)
		m.SetHeader("To", destinationEmail)
		m.SetHeader("Subject", "Fwd: "+msg.Envelope.Subject)

		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("read mail error", err)
				break
			}
			switch h := part.Header.(type) {
			case *mail.InlineHeader:
				b, _ := io.ReadAll(part.Body)
				currentType, _, _ := h.ContentType()
				if currentType == "text/plain" {
					m.SetBody("text/plain", string(b))
				} else if currentType == "text/html" {
					m.SetBody("text/html", string(b))
				}
			case *mail.AttachmentHeader:
				fileName, _ := h.Filename()
				m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
					_, err := io.Copy(w, part.Body)
					return err
				}))
			}
		}
		d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: false}

		if err := d.DialAndSend(m); err != nil {
			log.Println("can not send mail", err)
			continue
		}

		log.Println("mail forwarded", msg.Envelope.Subject)
	}
	return nil
}
