package Config

import (
	_ "embed"
	"encoding/json"

	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Utils"
)

func Init() {
	__config := class{}
	if err := json.Unmarshal([]byte(Utils.ReadFromFile("Config.json")), &__config); err != nil {
		Log.WriteLogf("Warning:  Failed to Load Config File: %s", err.Error())
		return
	}

	// Use environment variables if present, else fallback to config.json values
	if envPort := Utils.GetEnvString("SERVER_PORT"); envPort != nil && *envPort != "" {
		Port = *envPort
	} else if __config.Port != "" {
		Port = __config.Port
	}
	if envHost := Utils.GetEnvString("SERVER_HOST"); envHost != nil && *envHost != "" {
		Host = *envHost
	} else if __config.Host != "" {
		Host = __config.Host
	}
	if envHttp := Utils.GetEnvString("SERVER_HTTPS"); envHttp != nil && *envHttp != "" {
		Https = *envHttp == "true"
	} else {
		Https = __config.Https
	}

	if envBuild := Utils.GetEnvString("BUILD"); envBuild != nil && *envBuild != "" {
		Build = *envBuild == "true"
	} else {
		Build = __config.Build
	}

	if __config.StaticFolders != nil {
		StaticFolders = __config.StaticFolders
	}
	if __config.CssFolders != nil {
		CssFolders = __config.CssFolders
	}
	if __config.JsFolders != nil {
		JsFolders = __config.JsFolders
	}
	if __config.ViewFolder != "" {
		ViewFolder = __config.ViewFolder
	}
}
