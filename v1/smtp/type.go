package smtp

import "net/smtp"

type (
	sMTPConfig struct {
		client *smtp.Client
	}
	sMTPClient struct {
		sMTPConfig *sMTPConfig
		auth       smtp.Auth
		address    string

		host               string
		port               int
		username, password string
		initialised        bool
	}
)
