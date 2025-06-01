package Config

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Utils"
)

func Init() {
	__config := class{}
	if err := json.Unmarshal([]byte(Utils.ReadFromFile("Config.json")), &__config); err != nil {
		Log.WriteLogf("Error Loading Config File: %s", err.Error())
		os.Exit(1)
		return
	}

	Log.WriteLog("Config Of the Server Loaded Successfully: ", __config)

	if __config.Port != "" {
		Port = __config.Port
	}
	if __config.Host != "" {
		Host = __config.Host
	}
	Http = __config.Https
	Build = __config.Build
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
