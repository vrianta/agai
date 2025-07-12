package config

import (
	"encoding/json"

	Utils "github.com/vrianta/agai/v1/utils"
)

func init() {
	__config := databaseConfigStruct{}
	json.Unmarshal([]byte(Utils.ReadFromFile(dBConfigFile)), &__config)

	if envHost := Utils.GetEnvString("DB_HOST"); envHost != nil && *envHost != "" {
		databaseConfig.Host = *envHost
	} else if __config.Host != "" {
		databaseConfig.Host = __config.Host
	}

	if envPort := Utils.GetEnvString("DB_PORT"); envPort != nil && *envPort != "" {
		databaseConfig.Port = *envPort
	} else if __config.Port != "" {
		databaseConfig.Port = __config.Port
	}

	if envUser := Utils.GetEnvString("DB_USER"); envUser != nil && *envUser != "" {
		databaseConfig.User = *envUser
	} else if __config.User != "" {
		databaseConfig.User = __config.User
	}

	if envPassword := Utils.GetEnvString("DB_PASSWORD"); envPassword != nil && *envPassword != "" {
		databaseConfig.Password = *envPassword
	} else if __config.Password != "" {
		databaseConfig.Password = __config.Password
	}

	if envDatabase := Utils.GetEnvString("DB_DATABASE"); envDatabase != nil && *envDatabase != "" {
		databaseConfig.Database = *envDatabase
	} else if __config.Database != "" {
		databaseConfig.Database = __config.Database
	}

	if envProtocol := Utils.GetEnvString("DB_PROTOCOL"); envProtocol != nil && *envProtocol != "" {
		databaseConfig.Protocol = *envProtocol
	} else if __config.Protocol != "" {
		databaseConfig.Protocol = __config.Protocol
	}

	if envDriver := Utils.GetEnvString("DB_DRIVER"); envDriver != nil && *envDriver != "" {
		databaseConfig.Driver = *envDriver
	} else if __config.Driver != "" {
		databaseConfig.Driver = __config.Driver
	}

	if envSSLMode := Utils.GetEnvString("DB_SSLMODE"); envSSLMode != nil && *envSSLMode != "" {
		databaseConfig.SSLMode = *envSSLMode
	} else if __config.SSLMode != "" {
		databaseConfig.SSLMode = __config.SSLMode
	}
}
