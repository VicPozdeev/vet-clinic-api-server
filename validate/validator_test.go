package validate

import (
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	enAlphasTestSting         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ruAlphasTestSting         = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
	numericTestSting          = "0123456789"
	asciiSymbolsTestString    = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	usernameSymbolsTestString = "_.-"
)

type ruAlphaTag struct {
	S string `validate:"rualpha"`
}

type ruAlphaNumericTag struct {
	S string `validate:"rualphanum"`
}

type ruPrintableASCIITag struct {
	S string `validate:"ruprintascii"`
}

type usernameTag struct {
	S string `validate:"username"`
}

type passwordTag struct {
	S string `validate:"password"`
}

func TestValidator(t *testing.T) {
	v := NewValidator(validator.New())
	validationErrs := &validator.ValidationErrors{}

	assert.Empty(t, v.Validate(ruAlphaTag{S: ruAlphasTestSting + enAlphasTestSting}))
	assert.ErrorAs(t, v.Validate(ruAlphaTag{S: ruAlphasTestSting + enAlphasTestSting + numericTestSting}), validationErrs)
	assert.ErrorAs(t, v.Validate(ruAlphaTag{S: ruAlphasTestSting + enAlphasTestSting + asciiSymbolsTestString}), validationErrs)

	assert.Empty(t, v.Validate(ruAlphaNumericTag{S: ruAlphasTestSting + enAlphasTestSting + numericTestSting}))
	assert.ErrorAs(t, v.Validate(ruAlphaNumericTag{S: ruAlphasTestSting + enAlphasTestSting + numericTestSting + asciiSymbolsTestString}), validationErrs)

	assert.Empty(t, v.Validate(ruPrintableASCIITag{S: ruAlphasTestSting + enAlphasTestSting + numericTestSting + asciiSymbolsTestString}))
	assert.ErrorAs(t, v.Validate(ruPrintableASCIITag{S: "\n\t\b"}), validationErrs)

	assert.Empty(t, v.Validate(usernameTag{S: enAlphasTestSting + numericTestSting + usernameSymbolsTestString}))
	assert.ErrorAs(t, v.Validate(usernameTag{S: ""}), validationErrs)
	assert.ErrorAs(t, v.Validate(usernameTag{S: "1user"}), validationErrs)
	assert.ErrorAs(t, v.Validate(usernameTag{S: ruAlphasTestSting}), validationErrs)
	assert.ErrorAs(t, v.Validate(usernameTag{S: asciiSymbolsTestString}), validationErrs)

	assert.Empty(t, v.Validate(passwordTag{S: enAlphasTestSting + numericTestSting + asciiSymbolsTestString}))
	assert.ErrorAs(t, v.Validate(passwordTag{S: ""}), validationErrs)
	assert.ErrorAs(t, v.Validate(passwordTag{S: "PASS1!"}), validationErrs)
	assert.ErrorAs(t, v.Validate(passwordTag{S: "pass1!"}), validationErrs)
	assert.ErrorAs(t, v.Validate(passwordTag{S: "Pass!"}), validationErrs)
	assert.ErrorAs(t, v.Validate(passwordTag{S: "Pass1"}), validationErrs)
	assert.ErrorAs(t, v.Validate(passwordTag{S: ruAlphasTestSting + "Pass1!"}), validationErrs)
}
