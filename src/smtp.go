package main

import (
	"os"
	"fmt"
	// "strconv"
	"time"
	"gopkg.in/mail.v2"
	"log"
)

type SmtpService struct {
	config *Config
}

func (ss *SmtpService) SendMail(attachmentPath string) error {

	_, err := os.Stat(attachmentPath)
	if err != nil {
		log.Fatalln("Attachment not statable:", attachmentPath)
	}

	m := mail.NewMessage()
	m.SetHeader("From", ss.config.Smtp.Sender)
	m.SetHeader("To", ss.config.Smtp.Recipient)
	m.SetHeader("Subject", ss.config.Smtp.Subject)
	m.Attach(attachmentPath)
	d := mail.NewDialer(
		ss.config.Smtp.Host.String(), 
		ss.config.Smtp.Port, 
		ss.config.Smtp.User, 
		ss.config.Smtp.Pass,
	)
	d.Timeout = 10 * time.Second

	return d.DialAndSend(m)
}

func NewSmtpService(cfg *Config) (*SmtpService, error) {

	c := cfg.Smtp

	if c.Host.String() == "" {
		return nil, fmt.Errorf("SMTP Host missing")
	}
	
	/*	
	if c.Port == nil {
		return nil, fmt.Errorf("SMTP_PORT missing")
	}
	*/
	/*
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT invalid: %q", portStr)
	}
	*/
	
	if c.User == "" {
		return nil, fmt.Errorf("SMTP User missing")
	}

	if c.Pass == "" {
		return nil, fmt.Errorf("SMTP Pass missing")
	}

	if c.Recipient == "" {
		return nil, fmt.Errorf("SMTP Recipient missing")
	}

	if c.Sender == "" {
		return nil, fmt.Errorf("SMTP Sender missing")
	}

	
	if c.Subject == "" {
		return nil, fmt.Errorf("SMTP Subject missing")
	}

	return &SmtpService{config: cfg}, nil
}