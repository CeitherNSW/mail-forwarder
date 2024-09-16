package main

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/emersion/go-imap"
	// "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/go-gomail/gomail"
)

func forwardEmails(messages []*imap.Message, smtpHost string, smtpPort int, smtpUser, smtpPass, destinationEmail string) error {
	sem := make(chan struct{}, 10) // 限制最大并发数为10
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errList []error
	for _, msg := range messages {
		wg.Add(1)
		sem <- struct{}{} // 当并发超过10时会阻塞

		go func(msg *imap.Message) {
			defer wg.Done()
			defer func() { <-sem }() // 处理完释放一个并发槽

			if msg == nil {
				return
			}
			section := &imap.BodySectionName{}
			r := msg.GetBody(section)
			if r == nil {
				log.Println("server can't fetch message body")
				return
			}
			mr, err := mail.CreateReader(r)
			if err != nil {
				log.Println("can not read the mail", err)
				return
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
				log.Println("cannot send mail", err)
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				return
			}
		}(msg)
	}
	wg.Wait()
	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}
