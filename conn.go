package main

import (
	"crypto/tls"
	"github.com/emersion/go-imap/client"
	"log"
)

func connectIMAP(username, password, server string) (*client.Client, error) {
	c, err := client.DialTLS(server, &tls.Config{})
	if err != nil {
		return nil, err
	}

	if err := c.Login(username, password); err != nil {
		return nil, err
	}
	log.Println("Connected to IMAP server")
	return c, nil
}
