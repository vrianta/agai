package Config

import (
	"encoding/json"

	"github.com/vrianta/Server/Utils"
)

func New() {
	__config := class{}
	if err := json.Unmarshal([]byte(Utils.ReadFromFile("Config.json")), &__config); err != nil {
		panic(err)

	}

	Http = __config.Http
	Build = __config.Build
	if __config.StaticFolder != nil {
		StaticFolder = __config.StaticFolder
	}

	if __config.CssFolders != nil {
		CssFolder = __config.CssFolders
	}
	if __config.JsFolders != nil {
		JsFolders = __config.JsFolders
	}

	if __config.ViewFolder != "" {
		ViewFolder = __config.ViewFolder
	}
}
