package config

import (
	_ "embed"
	"encoding/json"
	"strconv"

	"github.com/vrianta/agai/v1/utils"
)

func init() {
	__config := applicationConfig{}
	json.Unmarshal([]byte(utils.ReadFromFile(webConfigFile)), &__config)

	// Use environment variables if present, else fallback to config.json values
	if envPort := utils.GetEnvString("SERVER_PORT"); envPort != nil && *envPort != "" {
		webConfig.Port = *envPort
	} else if __config.Port != "" {
		webConfig.Port = __config.Port
	}
	if envHost := utils.GetEnvString("SERVER_HOST"); envHost != nil && *envHost != "" {
		webConfig.Host = *envHost
	} else if __config.Host != "" {
		webConfig.Host = __config.Host
	}
	if envHttp := utils.GetEnvString("SERVER_HTTPS"); envHttp != nil && *envHttp != "" {
		webConfig.Https = *envHttp == "true"
	} else {
		webConfig.Https = __config.Https
	}

	if envBuild := utils.GetEnvString("BUILD"); envBuild != nil && *envBuild != "" {
		webConfig.Build = *envBuild == "true"
	} else {
		webConfig.Build = __config.Build
	}

	// MaxSessionCount: environment variable takes precedence
	if envMax := utils.GetEnvString("MAX_SESSION_COUNT"); envMax != nil && *envMax != "" {
		if v, err := strconv.Atoi(*envMax); err == nil {
			webConfig.MaxSessionCount = v
		}
	} else if __config.MaxSessionCount > 0 {
		webConfig.MaxSessionCount = __config.MaxSessionCount
	}

	// SessionStoreType: environment variable takes precedence
	if envStoreType := utils.GetEnvString("SESSION_STORE_TYPE"); envStoreType != nil && *envStoreType != "" {
		webConfig.SessionStoreType = *envStoreType
	} else if __config.SessionStoreType != "" {
		webConfig.SessionStoreType = __config.SessionStoreType
	}

	if __config.StaticFolders != nil {
		webConfig.StaticFolders = __config.StaticFolders
	}
	if __config.CssFolders != nil {
		webConfig.CssFolders = __config.CssFolders
	}
	if __config.JsFolders != nil {
		webConfig.JsFolders = __config.JsFolders
	}
	if __config.ViewFolder != "" {
		webConfig.ViewFolder = __config.ViewFolder
	}

	if __config.PassordCost != 0 {
		webConfig.PassordCost = __config.PassordCost
	} else {
		webConfig.PassordCost = 10 // default cost
	}
}

// Load The Web Config at runtime
func Load_Web() {
	__config := applicationConfig{}
	json.Unmarshal([]byte(utils.ReadFromFile(webConfigFile)), &__config)

	// Use environment variables if present, else fallback to config.json values
	if envPort := utils.GetEnvString("SERVER_PORT"); envPort != nil && *envPort != "" {
		webConfig.Port = *envPort
	} else if __config.Port != "" {
		webConfig.Port = __config.Port
	}
	if envHost := utils.GetEnvString("SERVER_HOST"); envHost != nil && *envHost != "" {
		webConfig.Host = *envHost
	} else if __config.Host != "" {
		webConfig.Host = __config.Host
	}
	if envHttp := utils.GetEnvString("SERVER_HTTPS"); envHttp != nil && *envHttp != "" {
		webConfig.Https = *envHttp == "true"
	} else {
		webConfig.Https = __config.Https
	}

	if envBuild := utils.GetEnvString("BUILD"); envBuild != nil && *envBuild != "" {
		webConfig.Build = *envBuild == "true"
	} else {
		webConfig.Build = __config.Build
	}

	// MaxSessionCount: environment variable takes precedence
	if envMax := utils.GetEnvString("MAX_SESSION_COUNT"); envMax != nil && *envMax != "" {
		if v, err := strconv.Atoi(*envMax); err == nil {
			webConfig.MaxSessionCount = v
		}
	} else if __config.MaxSessionCount > 0 {
		webConfig.MaxSessionCount = __config.MaxSessionCount
	}

	// SessionStoreType: environment variable takes precedence
	if envStoreType := utils.GetEnvString("SESSION_STORE_TYPE"); envStoreType != nil && *envStoreType != "" {
		webConfig.SessionStoreType = *envStoreType
	} else if __config.SessionStoreType != "" {
		webConfig.SessionStoreType = __config.SessionStoreType
	}

	if __config.StaticFolders != nil {
		webConfig.StaticFolders = __config.StaticFolders
	}
	if __config.CssFolders != nil {
		webConfig.CssFolders = __config.CssFolders
	}
	if __config.JsFolders != nil {
		webConfig.JsFolders = __config.JsFolders
	}
	if __config.ViewFolder != "" {
		webConfig.ViewFolder = __config.ViewFolder
	}
}
