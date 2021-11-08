package validator

import (
	"context"
	"reflect"
	"strings"
	"unicode"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/graphql"

	"github.com/go-playground/locales/en"

	ut "github.com/go-playground/universal-translator"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/hayashiki/audiy-api/cmd/errs"

	"github.com/go-playground/validator/v10"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	en := en.New()
	uni = ut.New(en, en)
	ent, _ := uni.GetTranslator("en")
	ent.Add("internal error", "An unexpected error has occurred.", true)
	ent.Add("login error", "Incorrect username or password.", true)

	trans = ent
	validate = validator.New()
	//validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	//	fieldName := fld.Tag.Get("jaFieldName")
	//	if fieldName == "-" {
	//		return ""
	//	}
	//	return fieldName
	//})
	//ja_translations.RegisterDefaultTranslations(validate, trans)
}

// Validate validates given struct using go-playground/validator
func Validate(src interface{}) error {
	return validate.Struct(src)
}

// GetErrorMessages エラーメッセージ群の取得
func GetErrorMessages(err error) error {
	if err == nil {
		return nil
	}
	var messages []string
	errs := err.(validator.ValidationErrors)
	extensions := make(map[string]interface{}, len(errs))
	for key, value := range errs.Translate(trans) {
		messages = append(messages, value)
		extensions[key] = value
		//User.Email
		//Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag
		println(key)
		println(value)
		//graphql.AddError(ctx, gqlerror.Errorf("field '%s' with value '%s' violates constraint: %s", err.Field(), err.Value(), err.Tag()))

	}
	return &gqlerror.Error{
		Message:    "Validation error",
		Extensions: extensions,
	}
}

func ManageValidationsErrors(ctx context.Context, validationErrors error) (gqlError *gqlerror.Error) {
	if validationErrors != nil {
		for i, err := range validationErrors.(validator.ValidationErrors) {
			//graphql.AddError(ctx, gqlerror.Errorf("field '%s' with value '%s' violates constraint: %s", err.Field(), err.Value(), err.Tag()))
			rule := err.ActualTag()
			param := err.Param()
			if param != "" {
				rule += "=" + param
			}

			// Recreate graphql path for input argument.
			path := graphql.GetPath(ctx)
			println(graphql.GetErrors(ctx).Error())
			namespace := strings.Split(err.Namespace(), ".")[1:]
			for _, name := range namespace {
				path = append(path, ast.PathName(makeFirstCharLowercase(name)))
			}

			//
			if i == 0 {
				gqlError = &gqlerror.Error{
					Extensions: map[string]interface{}{
						"code":  errs.InputValidationError,
						"value": err.Value(),
						"rule":  rule,
					},
				}
			} else {
				// Add error to response.
				graphql.AddError(ctx, &gqlerror.Error{
					Path:    path,
					Message: err.Error(),
					Extensions: map[string]interface{}{
						"code":  errs.InputValidationError,
						"value": err.Value(),
						"rule":  rule,
					},
				})
			}
		}
	}

	return
}

func makeFirstCharLowercase(s string) string {
	for _, v := range s {
		return string(unicode.ToLower(v)) + s[1:]
	}

	return s
}

//v := validator.New()
type ErrField struct {
	Field   string
	Message string
}

type ErrFields []ErrField

func ValidateErrors(value interface{}) error {
	if err := validate.Struct(value); err != nil {
		errs := err.(validator.ValidationErrors)
		extensions := make(map[string]interface{}, len(errs))
		//translator := en.New()
		//uni := ut.New(translator, translator)
		//trans, _ :=
		//translator.GetTranslator("en")

		for _, value := range errs {
			extensions[value.Field()] = value.Value()
		}

		return &gqlerror.Error{
			Message:    "Validation error",
			Extensions: extensions,
		}
	}
	return nil
}

type Validator struct {
	*validator.Validate
}

func New() *Validator {
	validate := validator.New()
	v := &Validator{
		validate,
	}
	return v
}

func (v Validator) CheckStruct(ctx context.Context, s interface{}, canBeNil bool) bool {
	if s == nil || (reflect.ValueOf(s).Kind() == reflect.Ptr && reflect.ValueOf(s).IsNil()) {
		return canBeNil
	}
	err := v.Struct(s)
	if err != nil {
		v.AddErrs(ctx, err)
		return false
	}
	return true
}

func (v Validator) AddErrs(ctx context.Context, err error) {
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs.Add(ctx, errs.Validation(ctx, e.Field()))
		}
	}
}
