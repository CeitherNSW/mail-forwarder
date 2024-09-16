package main

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)

func fetchEmails(c *client.Client) ([]*imap.Message, error) {
	mBox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}
	log.Println("Mailbox status:", mBox.Messages, mBox.Recent)
	if mBox.Messages == 0 {
		log.Println("No messages in mailbox")
		return nil, nil
	}

	criteria := imap.NewSearchCriteria()
	ids, err := c.Search(criteria)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		log.Println("No messages found in mailbox")
		return nil, nil
	}

	seqs := new(imap.SeqSet)
	seqs.AddNum(ids...)

	message := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqs, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody}, message)
	}()

	var res []*imap.Message
	for msg := range message {
		res = append(res, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}
	return res, nil
}
