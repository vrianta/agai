package component

import "github.com/vrianta/Server/models"

type (
	Struct struct {
		model *models.Struct // Model associated with the component
		data  any            // Data to be used in the component
	}
)
