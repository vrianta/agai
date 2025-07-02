package component

import (
	"github.com/vrianta/Server/models"
)

/*
LOGIC: user will send a struct in form of T(Type) and it will return a Component obj
This Component will have a model associated with it and a data field which will be used to store the data
Later in initialise function we will loop through all the components and populate them with the data
*/
func New[T any](m *models.Struct) *Struct {
	data_obj := new(T)
	response := &Struct{
		model: m,
		data:  &data_obj,
	}

	component_storage = append(component_storage, response)
	return response
}
