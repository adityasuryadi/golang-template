package pkg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidation(vl *validator.Validate) *Validation {
	vl.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validation{
		validate: vl,
	}
}

func removeFirstNameSpace(namespace string) string {
	s := strings.Split(namespace, ".")
	if len(s) > 1 {
		arr := make([]string, 0, len(s))
		for i := 1; i < len(s); i++ {
			arr = append(arr, s[i])
		}
		result := strings.Join([]string(arr), ".")
		return result
	}
	return namespace
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field " + fe.Field() + " tidak boleh kosong"
	case "lte":
		return "harus lebih kecil dari " + fe.Param()
	case "gtenow":
		return "harus lebih besar dari tanggal hari ini"
	case "gte":
		return "harus lebih besar dari " + fe.Param()
	case "email":
		return "format email salah"
	case "unique":
		return "data exist"
	case "min":
		return "minimal " + fe.Param() + " karakter"
	case "max":
		return "max " + fe.Param() + " Kb"
	case "image_validation":
		return "Harus Image"
	case "number":
		return "harus numeric"
	}
	return "Unknown error"
}

func (v *Validation) ValidateRequest(request interface{}) error {
	err := v.validate.Struct(request)
	if err != nil {
		return err
	}
	return nil
}

func (v *Validation) ErrorJson(err error) []ErrorMessage {
	fmt.Println("validationssssssss ", err)
	validationErrors := err.(validator.ValidationErrors)

	out := make([]ErrorMessage, len(validationErrors))
	for i, fieldError := range validationErrors {
		out[i] = ErrorMessage{
			Field:   removeFirstNameSpace(fieldError.Namespace()),
			Message: GetErrorMsg(fieldError),
		}
	}
	return out
}
