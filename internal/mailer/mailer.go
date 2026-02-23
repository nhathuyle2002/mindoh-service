package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// Mailer sends emails via SMTP.
type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewMailer creates a new Mailer.
func NewMailer(host string, port int, username, password, from string) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// Send sends an HTML email.
func (m *Mailer) Send(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	headers := strings.Join([]string{
		"From: " + m.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
	}, "\r\n")

	msg := []byte(headers + "\r\n\r\n" + htmlBody)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	// Use TLS for port 465, STARTTLS for 587/25
	if m.Port == 465 {
		tlsCfg := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         m.Host,
		}
		conn, err := tls.Dial("tcp", addr, tlsCfg)
		if err != nil {
			return fmt.Errorf("tls dial: %w", err)
		}
		client, err := smtp.NewClient(conn, m.Host)
		if err != nil {
			return fmt.Errorf("smtp new client: %w", err)
		}
		defer client.Close()
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
		if err = client.Mail(m.From); err != nil {
			return err
		}
		if err = client.Rcpt(to); err != nil {
			return err
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		return w.Close()
	}

	// STARTTLS (port 587)
	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}

// IMailer is the interface for sending email (allows mocking in tests).
type IMailer interface {
	Send(to, subject, htmlBody string) error
}

// NoopMailer is a no-op mailer used when SMTP is not configured (dev mode).
type NoopMailer struct{}

func (n *NoopMailer) Send(to, subject, htmlBody string) error {
	fmt.Printf("[mailer] (noop) to=%s subject=%s\n", to, subject)
	return nil
}
