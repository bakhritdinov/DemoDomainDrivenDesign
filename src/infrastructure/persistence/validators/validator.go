package validators

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func InitValidator(lang string) {
	validate = validator.New()

	english := en.New()
	russian := ru.New()
	uni := ut.New(english, english, russian)

	var found bool
	translator, found = uni.GetTranslator(lang)
	if !found {
		translator, _ = uni.GetTranslator("en") // fallback to English
	}

	switch lang {
	case "ru":
		_ = rutranslations.RegisterDefaultTranslations(validate, translator)
	default:
		_ = entranslations.RegisterDefaultTranslations(validate, translator)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	registerCustomMessages()
}

func ValidateStruct(model interface{}) error {
	err := validate.Struct(model)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessages := make([]string, len(validationErrors))

		for i, e := range validationErrors {
			errorMessages[i] = e.Translate(translator)
		}

		return fmt.Errorf(strings.Join(errorMessages, "; "))
	}

	return err
}

func RegisterValidationCallbacks(db *gorm.DB) {
	callback := db.Callback()

	_ = callback.Create().Before("gorm:create").Register("validations:validate", validateModel)
	_ = callback.Update().Before("gorm:update").Register("validations:validate", validateModel)
}

func validateModel(db *gorm.DB) {
	if db.Statement.Schema == nil || db.Statement.Dest == nil {
		return
	}

	if _, ok := db.InstanceGet("skip_validations"); ok {
		return
	}

	model := db.Statement.Dest
	if err := ValidateStruct(model); err != nil {
		_ = db.AddError(err)
	}
}

func registerCustomMessages() {
	_ = validate.RegisterTranslation("gte", translator,
		func(ut ut.Translator) error {
			return ut.Add("gte", "{0} должен быть не менее {1}", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("gte", fe.Field(), fe.Param())
			return t
		},
	)
}
