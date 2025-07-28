package session

import (
	"github.com/vrianta/agai/v1/model"
)

/*
Model for session and it is holding the database record of the session
*/

// sessionModel is the private variable which is holding the sessions details
// it is goood idea to keep the name of the model in small caps
var SessionModel = model.New("sessions", struct {
	Id   *model.Field
	Data *model.Field
}{
	Id: &model.Field{
		Type:   model.FieldTypes.VarChar,
		Length: 100,
		Index: model.Index{
			PrimaryKey: true,
			Index:      true,
		},
		Nullable: false,
	},
	Data: &model.Field{
		Type:     model.FieldTypes.Text,
		Nullable: true,
	},
	// IsAuthenticated: model.Field{
	// 	Type:     model.FieldTypes.Bool,
	// 	Nullable: false,
	// },
	// ExpirationTime: model.Field{
	// 	Type:     model.FieldTypes.Time,
	// 	Nullable: true,
	// },
})
