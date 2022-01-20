package handlers

import (
	. "backend/model"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateTagErrorMsg recieve err=struct tag error, and check if required field/empty field/invalid email format
//is the cause of the error , and return str= informative error.
// if the fields are not the problem, return str= the non informative error.
func ValidateTagErrorMsg(err error) string {
	str := ""
	if veErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range veErr {
			switch {
			case fieldErr.Tag() == "required":
				str = str + fmt.Sprintf("%s is required and is mandatory field\n", fieldErr.Field())
			case fieldErr.Tag() == "email":
				str = str + fmt.Sprintf("Invalid email\n")
			case fieldErr.Tag() == "notemptystring":
				str = str + fmt.Sprintf("%s is empty, please dont leave it like that.\n", fieldErr.Field())
			}

		}
	}
	if str == "" {
		str = fmt.Sprintf("Invalid data. Check this error:\n%s", err.Error())
	}

	return str
}

//ValidateDbErrorMsg recieve err=db error, params=what cause the error(e.g. id,email),Type= to which type it belongs(e.g. person,task)
// and return str =informative error.  return the non informative error if id and email are not the cause of the error.
func ValidateDbErrorMsg(err error, Type string, params *OptParams) string {
	str := ""
	if strings.Contains(err.Error(), "UNIQUE constraint") {
		str = fmt.Sprintf("A person with email '%s' already exists.", *params.Email)

	} else {
		if err.Error() == "record not found" {
			str = fmt.Sprintf("A %s with the id '%s'does not exist. ", Type, params.ID)
		} else {
			str = fmt.Sprintf("DataBase error:\n%s", err.Error())
		}
	}

	return str
}

//NotEmptyString
var NotEmptyString validator.Func = func(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if str == "" {
		return false
	}
	tr := strings.TrimSpace(str)

	if ok {
		if tr == "" {
			return false
		}
		return true
	}
	return false
}
