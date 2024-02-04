package validate

import (
	"github.com/go-playground/validator"
	"regexp"
)

const (
	ruAlphaRegexString                 = `^[а-яА-ЯёЁa-zA-Z]+$`
	ruAlphaNumericRegexString          = `^[а-яА-ЯёЁa-zA-Z0-9]+$`
	ruPrintableASCIIRegexString        = `^[а-яА-ЯёЁ\x20-\x7E]*$`
	usernameRegexString                = `^[a-zA-Z][a-zA-Z0-9_.-]*$`
	passwordRegexString                = `^[\x20-\x7E]+$`
	containsLowercaseLetterRegexString = `[a-z]`
	containsUppercaseLetterRegexString = `[A-Z]`
	containsDigitRegexString           = `\d`
	containsSymbolRegexString          = "[ !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]"
)

var (
	ruAlphaRegex                 = regexp.MustCompile(ruAlphaRegexString)
	ruAlphaNumericRegex          = regexp.MustCompile(ruAlphaNumericRegexString)
	ruPrintableASCIIRegex        = regexp.MustCompile(ruPrintableASCIIRegexString)
	usernameRegex                = regexp.MustCompile(usernameRegexString)
	passwordRegex                = regexp.MustCompile(passwordRegexString)
	containsLowercaseLetterRegex = regexp.MustCompile(containsLowercaseLetterRegexString)
	containsUppercaseLetterRegex = regexp.MustCompile(containsUppercaseLetterRegexString)
	containsDigitRegex           = regexp.MustCompile(containsDigitRegexString)
	containsSymbolRegex          = regexp.MustCompile(containsSymbolRegexString)
)

var customValidators = map[string]validator.Func{
	"rualpha":      isRuAlpha,
	"rualphanum":   isRuAlphaNumeric,
	"ruprintascii": isRuPrintableASCII,
	"username":     isUsername,
	"password":     isPassword,
}

type Validator struct {
	validator *validator.Validate
}

// NewValidator is constructor.
func NewValidator(validator *validator.Validate) *Validator {
	for tag, fn := range customValidators {
		if err := validator.RegisterValidation(tag, fn); err != nil {
			panic(err)
		}
	}
	return &Validator{validator: validator}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func isRuAlpha(fl validator.FieldLevel) bool {
	return ruAlphaRegex.MatchString(fl.Field().String())
}

func isRuAlphaNumeric(fl validator.FieldLevel) bool {
	return ruAlphaNumericRegex.MatchString(fl.Field().String())
}

func isRuPrintableASCII(fl validator.FieldLevel) bool {
	return ruPrintableASCIIRegex.MatchString(fl.Field().String())
}

func isUsername(fl validator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

func isPassword(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if containsLowercaseLetterRegex.MatchString(str) &&
		containsUppercaseLetterRegex.MatchString(str) &&
		containsDigitRegex.MatchString(str) &&
		containsSymbolRegex.MatchString(str) {
		return passwordRegex.MatchString(str)
	}
	return false
}
