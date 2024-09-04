package validator

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testValidator = New()
var ctx = context.Background()

func TestNameValidation(t *testing.T) {
	t.Parallel()

	shouldPass := []string{
		"Pismo", "Jyotsna Gupta", "Chad",
		"Visa Seturam", "A B C D E", "gg",
		"This name has exactly hundred characters but it is not really a name at this point but a short story",
		"V P",
		"a",
		"a a a",
		"a a a a a a",
	}

	for _, input := range shouldPass {
		name := input
		t.Run(fmt.Sprintf("should pass for valid name: %s", name), func(t *testing.T) {
			t.Parallel()

			res, err := testValidator.IsValidString(ctx, name, "name")
			assert.Nil(t, err)
			assert.Equal(t, true, res.Valid, name)
		})
	}

	shouldFail := []string{
		"Mr. Decentro",                     // punctuation
		"",                                 // minimum character = 1
		"F1 Engine",                        // numbers
		"Yes I have ✨ emojis ✨ in my name", // non-ascii characters (emojis)
		"すべてを一度にどこでも",                      // other non-ascii characters
		"       IRIS",                      // leading whitespace
		"IRIS      ",                       // trailing whitespace
		"      IRIS      ",                 // leading and trailing whitespace
		"hello\\xbd\\xb2=\\xbc ⌘",          // random bytes
		"\u200cEvil",                       // unicode for the invisible character
		"Tabs	are	better	but	not	here",     // tabs
		" ",                                // only whitespace
		"    ",                             // Multiple spaces
	}

	for _, input := range shouldFail {
		name := input
		t.Run(fmt.Sprintf("should fail for invalid name: %s", name), func(t *testing.T) {
			t.Parallel()

			res, err := testValidator.IsValidString(ctx, name, "name")
			assert.Nil(t, err)
			assert.Equal(t, false, res.Valid, name)
		})
	}
}
