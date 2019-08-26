package utils

import (
	"fmt"
	"github.com/natalizhy/blacklist/backend/models"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

var userError = map[string]map[string]string{
	"FirstName": {
		"required": "Обязательное поле",
		"cyr":      "Можна использовать только буквы",
		"min":      "Слишком мало символов",
	},
	"LastName": {
		"required": "Обязательное поле",
		"alpha":    "Можна использовать только буквы",
		"min":      "Слишком мало символов",
	},
	"Phone": {
		"required": "Обязательное поле",
		"numeric":  "Можна использовать только цифры",
		"min":      "Слишком мало символов",
	},
	"Info": {
		"required": "Обязательное поле",
		"min":      "Слишком мало символов",
	},
}

var validate = validator.New()

func ValidateCyr(fl validator.FieldLevel) bool {
	result := regexp.MustCompile("^[a-zA-ZА-Яа-я]+$")
	return result.MatchString(fl.Field().String())
}

func ValidateUser(user models.User) (errors map[string]map[string]string, err error) {
	err = validate.RegisterValidation("cyr", ValidateCyr)

	if err != nil {
		fmt.Println(err)
	}

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
