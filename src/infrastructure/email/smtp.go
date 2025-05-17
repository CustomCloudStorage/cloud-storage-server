package email

import (
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/CustomCloudStorage/utils"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Timeout  time.Duration
}

type SMTPMailer struct {
	cfg SMTPConfig
}

func NewSMTPMailer(cfg SMTPConfig) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	addr := net.JoinHostPort(m.cfg.Host, strconv.Itoa(m.cfg.Port))
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	headers := map[string]string{
		"From":         m.cfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `text/plain; charset="utf-8"`,
	}

	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s", k, v))
	}
	sb.WriteString("")
	sb.WriteString(body)
	msg := []byte(sb.String())

	dialer := &net.Dialer{Timeout: m.cfg.Timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "SMTP dial failed")
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.cfg.Host)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "create SMTP client failed")
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return utils.ErrInternal.Wrap(err, "SMTP auth failed")
	}

	if err := client.Mail(m.cfg.From); err != nil {
		return utils.ErrInternal.Wrap(err, "SMTP MAIL FROM failed")
	}

	if err := client.Rcpt(to); err != nil {
		return utils.ErrInternal.Wrap(err, "SMTP RCPT TO failed")
	}

	writer, err := client.Data()
	if err != nil {
		return utils.ErrInternal.Wrap(err, "SMTP DATA command failed")
	}

	_, err = writer.Write(msg)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "writing message failed")
	}
	if err := writer.Close(); err != nil {
		return utils.ErrInternal.Wrap(err, "closing DATA writer failed")
	}

	if err := client.Quit(); err != nil {
		return utils.ErrInternal.Wrap(err, "SMTP QUIT failed")
	}
	return nil
}
