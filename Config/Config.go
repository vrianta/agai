package Config

import (
	"encoding/json"
	"fmt"

	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Utils"
)

func New() {
	__config := class{}
	if err := json.Unmarshal([]byte(Utils.ReadFromFile("Config.json")), &__config); err != nil {
		fmt.Errorf("Error Loading Config File: %v", err)
	}

	Log.WriteLog("Config Of the Server Loaded Successfully: ", __config)
	Http = __config.Http
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
