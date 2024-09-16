package main

import (
	"github.com/emersion/go-imap/client"
	"log"
)

func main() {
	imapUser := ""
	imapPass := ""
	imapServer := "imap-mail.outlook.com:993"

	destEmail := ""

	smtpUser := ""
	smtpPass := ""
	smtpHost := "smtp-mail.outlook.com"
	smtpPort := 587

	c, err := connectIMAP(imapUser, imapPass, imapServer)
	if err != nil {
		log.Fatal("IMAP connection error:", err)
	}
	defer func(c *client.Client) {
		err := c.Logout()
		if err != nil {
			log.Fatal("Logout error:", err)
		}
	}(c)

	messages, err := fetchEmails(c)
	if err != nil {
		log.Fatal("Fetch emails error:", err)
	}

	if len(messages) == 0 {
		log.Println("No messages to forward")
		return
	}

	err = forwardEmails(messages, smtpHost, smtpPort, smtpUser, smtpPass, destEmail)
	if err != nil {
		log.Fatal("Forward emails error:", err)
	}

	log.Println("Emails forwarded successfully")
}
