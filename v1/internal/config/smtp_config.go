package config

import (
	"encoding/json"
	"strconv"

	Utils "github.com/vrianta/agai/v1/utils"
)

func init() {
	__smtp := smtpConfigStruct{}
	json.Unmarshal([]byte(Utils.ReadFromFile(smtpConfigFile)), &__smtp)

	if envHost := Utils.GetEnvString("SMTP_HOST"); envHost != nil && *envHost != "" {
		smtpConfig.Host = *envHost
	} else if __smtp.Host != "" {
		smtpConfig.Host = __smtp.Host
	}

	if envPort := Utils.GetEnvString("SMTP_PORT"); envPort != nil && *envPort != "" {
		if port, err := strconv.Atoi(*envPort); err == nil {
			smtpConfig.Port = port
		}
	} else if __smtp.Port > 0 {
		smtpConfig.Port = __smtp.Port
	}

	if envUser := Utils.GetEnvString("SMTP_USERNAME"); envUser != nil && *envUser != "" {
		smtpConfig.Username = *envUser
	} else if __smtp.Username != "" {
		smtpConfig.Username = __smtp.Username
	}

	if envPass := Utils.GetEnvString("SMTP_PASSWORD"); envPass != nil && *envPass != "" {
		smtpConfig.Password = *envPass
	} else if __smtp.Password != "" {
		smtpConfig.Password = __smtp.Password
	}

	if envTLS := Utils.GetEnvString("SMTP_USE_TLS"); envTLS != nil && *envTLS != "" {
		smtpConfig.UseTLS = *envTLS == "true"
	} else {
		smtpConfig.UseTLS = __smtp.UseTLS
	}
}
