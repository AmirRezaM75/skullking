package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"log"
	"strings"
)

type Validator struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func NewValidator() Validator {
	t := en.New()

	universalTranslator := ut.New(t, t)

	translator, found := universalTranslator.GetTranslator("en")

	if !found {
		log.Fatal("translator not found")
	}

	var validate = validator.New()

	if err := enTranslation.RegisterDefaultTranslations(validate, translator); err != nil {
		log.Fatal(err)
	}

	_ = validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} field is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("email", translator, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email address.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	return Validator{
		Validate:   validate,
		Translator: translator,
	}
}

func (v Validator) ValidateStruct(s interface{}) *ValidationError {
	err := v.Validate.Struct(s)

	if err != nil {
		e := ValidationError{}
		e.Message = "The given data is invalid."
		e.Errors = make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			e.Errors[strings.ToLower(err.Field())] = err.Translate(v.Translator)
		}

		return &e
	}

	return nil
}
