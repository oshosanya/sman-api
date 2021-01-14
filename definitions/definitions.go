package definitions

import validation "github.com/go-ozzo/ozzo-validation"

type IDCard struct {
	Name     string `json:"name" form:"name"`
	Position string `json:"position" form:"position"`
	Branch   string `json:"branch" form:"branch"`
	IDNumber string `json:"id_number" form:"id_number"`
}

func (id IDCard) Validate() error {
	return validation.ValidateStruct(&id,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&id.Name, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&id.Position, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be id string consisting of two letters in upper case
		validation.Field(&id.Branch, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be id string consisting of five digits
		validation.Field(&id.IDNumber, validation.Required, validation.Length(3, 50)),
	)
}
