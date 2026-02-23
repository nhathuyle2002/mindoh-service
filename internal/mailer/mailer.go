package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"strings"
)

// ─── Brevo HTTP API mailer ────────────────────────────────────────────────────

// BrevoMailer sends email via Brevo's transactional email HTTP API.
// Use this in environments where outbound SMTP ports are blocked (e.g. Railway).
type BrevoMailer struct {
	APIKey string
	From   string
}

func NewBrevoMailer(apiKey, from string) *BrevoMailer {
	return &BrevoMailer{APIKey: apiKey, From: from}
}

func (b *BrevoMailer) Send(to, subject, htmlBody string) error {
	slog.Info("sending email via Brevo API", "to", to, "subject", subject)

	payload := map[string]interface{}{
		"sender":      map[string]string{"email": b.From},
		"to":          []map[string]string{{"email": to}},
		"subject":     subject,
		"htmlContent": htmlBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("brevo: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("brevo: create request: %w", err)
	}
	req.Header.Set("api-key", b.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("brevo: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errBody)
		slog.Error("brevo API error", "status", resp.StatusCode, "body", errBody)
		return fmt.Errorf("brevo: unexpected status %d", resp.StatusCode)
	}
	slog.Info("email sent via Brevo API", "to", to)
	return nil
}

// ─── SMTP mailer ──────────────────────────────────────────────────────────────

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
	slog.Info("sending email", "to", to, "subject", subject, "smtp_host", m.Host, "smtp_port", m.Port)
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
	slog.Debug("[mailer] noop", "to", to, "subject", subject)
	return nil
}
