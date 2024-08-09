package validators

import (
	"github.com/go-playground/validator/v10"
	"log"
)

func ValidateNodeType(fl validator.FieldLevel) bool {
	nodeType := fl.Field().String()
	log.Printf("validation called on %s", fl.StructFieldName())
	for _, t := range []string{"file", "base"} {
		if t == nodeType {
			return true
		}
	}
	return false
}
