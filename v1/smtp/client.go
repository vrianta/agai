package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/log"
)

func init() {
	client_config := config.GetSmtpConfig() // smtp config

	if client_config.Host == "" || client_config.Username == "" || client_config.Password == "" {
		log.Debug("SMTP config is not settedup so we are moving forward without it")
		return
	}

	client.address = fmt.Sprintf("%s:%d", client_config.Host, client_config.Port)
	client.host = client_config.Host
	client.port = client_config.Port
	client.username = client_config.Username
	client.password = client_config.Password
	client.auth = smtp.PlainAuth("", client_config.Username, client_config.Password, client_config.Host)
	client.initialised = true
}

func (s *sMTPClient) SendMail(to []string, subject, body string) error {

	if !client.initialised {
		log.Error("SMTP Client is not initialised to send mail please initialise the error")
		return fmt.Errorf("SMTP Client is not initialised to send mail please initialise the error")
	}

	conn, err := net.Dial("tcp", s.address)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		fmt.Printf("failed to create client : %s", err.Error())
		return err
	}

	// Skip TLS setup to use an unencrypted connection
	if err := client.StartTLS(&tls.Config{InsecureSkipVerify: !config.GetSmtpConfig().UseTLS}); err != nil {
		log.Error("failed to start TLS: %s", err.Error())
		return err
	}

	if err := client.Auth(s.auth); err != nil {
		log.Error("failed to Authenticate client : %s", err.Error())
		return err
	}

	if err := client.Mail(s.username); err != nil {
		return err
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	defer writer.Close()

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.username, strings.Join(to, ","), subject, body)
	_, err = writer.Write([]byte(message))

	return err

}

func (s *sMTPClient) Close() error {
	if s.sMTPConfig == nil {
		return fmt.Errorf("smtp client not initialised")
	}
	return s.sMTPConfig.client.Quit()
}
