package errs

import (
	"context"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/stoewer/go-strcase"
)

// MainErrorCode are codes for representing main errors in the services.
type MainErrorCode string

const (
	InputValidationError      MainErrorCode = "InputValidationError"
	PreSessionValidationError MainErrorCode = "PreSessionValidationError"
	SessionValidationError    MainErrorCode = "SessionValidationError"
	InvalidBasicAuth          MainErrorCode = "InvalidBasicAuth"
	InternalSystemError       MainErrorCode = "InternalSystemError"
)

const (
	userCreationErrorMessage                = "could not create a new user"
	passwordCreationErrorMessage            = "could not create a new password"
	passwordUpdateErrorMessage              = "could not update password"
	passwordDeleteErrorMessage              = "could not delete password"
	passwordAuthenticationErrorMessage      = "unauthorized password input"
	userPasswordsFetchErrorMessage          = "could not fetch user's passwords"
	userPasswordsAuthenticationErrorMessage = "unauthorized passwords fetch"
	signInErrorMessage                      = "could not sign in"
	existingEmailErrorMessage               = "the e-mail address is already taken"
	queryNonExistingEmailErrorMessage       = "user doesn't exist"
	wrongPasswordErrorMessage               = "wrong password"
)

func Add(ctx context.Context, err error) {
	graphql.AddError(ctx, err)
}

func BadCredencials(ctx context.Context) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "Email or password invalid.",
		Extensions: map[string]interface{}{
			"code": "BAD_CREDENCIALS",
		},
	}
}

func Forbidden(ctx context.Context) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "Operation forbidden.",
		Extensions: map[string]interface{}{
			"code": "FORBIDDEN",
		},
	}
}

func NotFound(ctx context.Context) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "No data found.",
		Extensions: map[string]interface{}{
			"code": "NOT_FOUND",
		},
	}
}

func Exists(ctx context.Context) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "Resource already exists.",
		Extensions: map[string]interface{}{
			"code": "EXISTS",
		},
	}
}

func Internal(ctx context.Context) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "Internal server error.",
		Extensions: map[string]interface{}{
			"code": "INTERNAL_ERROR",
		},
	}
}

func Validation(ctx context.Context, field string) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: fmt.Sprintf("Field \"%s\" is invalid.", strcase.LowerCamelCase(field)),
		Extensions: map[string]interface{}{
			"code": fmt.Sprintf("VALIDATION_ERROR_%s", strings.ToUpper(field)),
		},
	}
}

var translator *ut.UniversalTranslator

func init() {
	en := en.New()

	translator = ut.New(en)
}

// Translator returns the instance of translator for a given locale.
// If the given locale doesn't exist, a fallback locale will be used
func Translator(locale string) ut.Translator {
	trans, _ := translator.GetTranslator(locale)
	return trans
}

// Trans returns a translated message for a given locale.
// If the locale doesn't exist, a fallback locale will be used
func Trans(locale, msg string, params ...string) string {
	trans, _ := translator.GetTranslator(locale)
	message, _ := trans.T(msg, params...)

	return message
}

func ValidationError(ctx context.Context, erras validator.ValidationErrors) error {
	extensions := make(map[string]interface{}, len(erras))
	//translator := en.New()
	//uni := ut.New(translator, translator)
	trans, _ := translator.GetTranslator("en")

	for key, value := range erras.Translate(trans) {
		extensions[key] = value
	}

	return &gqlerror.Error{
		Message:    "Validation error",
		Extensions: extensions,
	}
}

func GetHumanReadableError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required_without_all":
		return fmt.Sprintf("The %s field is required if the following fields are not set: %s", fe.Field(), strings.Join(strings.Split(fe.Param(), " "), ", "))
	case "required_without":
		return fmt.Sprintf("The %s field is required if any of the following fields are not set: %s", fe.Field(), strings.Join(strings.Split(fe.Param(), " "), ", "))
	case "url":
		return fmt.Sprintf("The %s field should be a valid URL.", fe.Field())
	case "min":
		return fmt.Sprintf("The %s field must have a minimum length of %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("The %s field must have a maximum length of %s", fe.Field(), fe.Param())
	case "length":
		return fmt.Sprintf("The %s field must have a length of %s", fe.Field(), fe.Param())
	case "WellKnownOAuthProviders":
		return fmt.Sprintf("%s is not a well-known OAuth Provider.", fe.Value())
	default:
		return fe.Error()
	}
}

func FieldsHaveValidationErrors(ctx context.Context, err error) bool {
	fieldErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}

	for _, fieldError := range fieldErrors {
		errPath := graphql.GetPath(ctx)
		errPath = append(errPath, ast.PathName(fieldError.StructField()))

		graphql.AddError(ctx, &gqlerror.Error{
			Path:       errPath,
			Message:    GetHumanReadableError(fieldError),
			Extensions: map[string]interface{}{},
		})
	}

	return true
}
