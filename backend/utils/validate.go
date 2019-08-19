package utils

import (
	"github.com/natalizhy/blacklist/backend/models"
	"gopkg.in/go-playground/validator.v9"
)

var userError = map[string]map[string]string{
	"FirstName": {
		"required": "Обязательное поле1",
		"alpha":    "Можна использовать только буквы",
		"min":    "Слишком мало символов",
	},
	"LastName": {
		"required": "Обязательное поле2",
		"alpha":    "Можна использовать только буквы",
		"min":    "Слишком мало символов",
	},
	"Phone": {
		"required": "Обязательное поле3",
		"numeric":  "Можна использовать только цифры",
		"min":    "Слишком мало символов",
	},
	"Info": {
		"required": "Обязательное поле4",
		"min":    "Слишком мало символов",
	},
}

func ValidateUser(user models.User) (errors map[string]map[string]string, err error) {

	validate := validator.New()

	err = validate.Struct(user)

	errors = make(map[string]map[string]string)

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			if _, ok := errors[err.Field()]; !ok {
				errors[err.Field()] = make(map[string]string)
			}

			errors[err.Field()][err.Tag()] = userError[err.Field()][err.Tag()]
		}
	}

	return
}
