package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator has functions for validating struct and variables
type Validator struct {
	validate *validator.Validate
}

// New returns a new Validator
func New() *Validator {
	v := buildValidator()
	vl := &Validator{validate: v}

	// We have to register the custom tags before using them.
	// Check the registerCustomTags() method for all the custom tags we have added.
	vl.registerCustomTags()

	return vl
}

// Result contains the details of the validation
type Result struct {
	// Valid is true when validation succeeds
	Valid bool
	// Fields contains the list of invalid fields. This is only set when Valid is false
	Fields []string
}

// IsValidStruct validates a struct and returns the results
func (v *Validator) IsValidStruct(ctx context.Context, data interface{}) (*Result, error) {
	if err := v.validate.StructCtx(ctx, data); err != nil {
		return HandleValidationError(err)
	}

	return &Result{
		Valid:  true,
		Fields: nil,
	}, nil
}

// IsValidString validates a string and returns the results
func (v *Validator) IsValidString(ctx context.Context, str string, tags string) (*Result, error) {
	if err := v.validate.VarCtx(ctx, str, tags); err != nil {
		return HandleValidationError(err)
	}

	return &Result{
		Valid:  true,
		Fields: nil,
	}, nil
}

type (
	validatorFn      func(fl validator.FieldLevel) bool
	stringModifierFn func(str string) string
)

// AddCustomValidator adds a custom validator tag
func (v *Validator) AddCustomValidator(name string, f validatorFn) error {
	if err := v.validate.RegisterValidation(name, validator.Func(f)); err != nil {
		return fmt.Errorf("failed to add custom  validator: %w", err)
	}

	return nil
}

// AddStructLevelValidation adds a custom struct level validation
func (v *Validator) AddStructLevelValidation(fn validator.StructLevelFunc, structType interface{}) {
	v.validate.RegisterStructValidation(fn, structType)
}

// AddStringModifier adds a tag that modifies the string value
func (v *Validator) AddStringModifier(name string, fn stringModifierFn) error {
	return v.AddCustomValidator(name, func(fl validator.FieldLevel) bool {
		if fl.Field().Type().String() == "string" {
			str := fn(fl.Field().String())
			fl.Field().SetString(str)
		}

		return true
	})
}

// HandleValidationError takes the validation error and read the failed fields
func HandleValidationError(err error) (*Result, error) {
	fields := make([]string, 0)

	if invalidErr, ok := err.(*validator.InvalidValidationError); ok {
		return nil, fmt.Errorf("HandleValidationError: %w", invalidErr)
	}

	for _, validationErr := range err.(validator.ValidationErrors) {
		fields = append(fields, validationErr.Field())
	}

	return &Result{
		Valid:  false,
		Fields: fields,
	}, nil
}

// buildValidator builds the validator.Validate.
// It also adds the function to read json tags from struct to use it for reposting errors.
func buildValidator() *validator.Validate {
	v := validator.New()

	// register function to get tag name from json tags.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("schema"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

var (
	//Only alphabets and spaces allowed
	//Minimum of 1 and maximum of 100 characters
	//leading or trailing whitespace not allowed
	nameRegex = regexp.MustCompile(`(?i)^((?:[a-z] ?[a-z ]{0,98}[a-z])|([a-z]))$`)
)

func (v *Validator) registerCustomTags() {
	//trim:  Trim the string i.e. remove leading and trailing whitespace
	_ = v.AddStringModifier("trim", strings.TrimSpace)

	// name: Validate name
	_ = v.AddCustomValidator("name", func(fl validator.FieldLevel) bool {
		return nameRegex.MatchString(fl.Field().String())
	})

	// true: Add custom validator for always true
	_ = v.AddCustomValidator("true", func(fl validator.FieldLevel) bool {
		return fl.Field().Bool()
	})

}
