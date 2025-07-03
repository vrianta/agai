package component

import "github.com/vrianta/Server/model"

type (
	Struct struct {
		model *model.Struct // Model associated with the component
		data  any           // Data to be used in the component
	}
)
