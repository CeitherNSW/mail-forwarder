package main

import "log"

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
	defer c.Logout()

	messages, err := fetchEmails(c)
	if err != nil {
		log.Fatal("Fetch emails error:", err)
	}

	if len(messages) == 0 {
		log.Println("No messages to forward")
		return
	}

	err = foewardMails(messages, smtpHost, smtpPort, smtpUser, smtpPass, destEmail)
	if err != nil {
		log.Fatal("Forward emails error:", err)
	}

	log.Println("Emails forwarded successfully")
}
