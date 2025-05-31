package RenderEngine

import "github.com/vrianta/Server/Template"

var (
	templateRecords = make(map[string]Template.Struct) // keep the reocrd of all the templated which are already templated

)
