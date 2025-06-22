package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckLuhn(t *testing.T) {
	t.Run("ValidLuhnNumbers", func(t *testing.T) {
		validNumbers := []string{
			"4532015112830366", // Valid Visa
			"5555555555554444", // Valid MasterCard
			"378282246310005",  // Valid American Express
			"6011111111111117", // Valid Discover
			"30569309025904",   // Valid Diners Club
			"38520000023237",   // Valid Diners Club
			"4111111111111111", // Valid test card
			"0",                // Single digit valid
			"18",               // Two digit valid
		}

		for _, number := range validNumbers {
			t.Run(number, func(t *testing.T) {
				result := CheckLuhn(number)
				assert.True(t, result, "Expected %s to be valid", number)
			})
		}
	})

	t.Run("InvalidLuhnNumbers", func(t *testing.T) {
		invalidNumbers := []string{
			"4532015112830367", // Invalid (last digit changed)
			"5555555555554445", // Invalid (last digit changed)
			"378282246310006",  // Invalid (last digit changed)
			"1234567890123456", // Invalid sequence
			"1111111111111111", // All ones (invalid)
			"1",                // Single digit invalid
			"19",               // Two digit invalid
		}

		for _, number := range invalidNumbers {
			t.Run(number, func(t *testing.T) {
				result := CheckLuhn(number)
				assert.False(t, result, "Expected %s to be invalid", number)
			})
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		invalidInputs := []string{
			"",            // Empty string
			"abc123",      // Contains letters
			"123-456-789", // Contains hyphens
			"123 456 789", // Contains spaces
			"12.34",       // Contains decimal point
			"1234a5678",   // Contains letter in middle
			"@#$%",        // Special characters
		}

		for _, input := range invalidInputs {
			t.Run(input, func(t *testing.T) {
				result := CheckLuhn(input)
				assert.False(t, result, "Expected %s to be invalid due to invalid characters", input)
			})
		}
	})

	t.Run("KnownTestCases", func(t *testing.T) {
		testCases := []struct {
			number   string
			expected bool
			name     string
		}{
			{"79927398713", true, "Valid 11-digit number"},
			{"79927398714", false, "Invalid 11-digit number"},
			{"4000000000000002", true, "Valid 16-digit test card"},
			{"4000000000000001", false, "Invalid 16-digit test card"},
			{"42424242424242424242", true, "Valid 20-digit number"},
			{"42424242424242424243", false, "Invalid 20-digit number"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := CheckLuhn(tc.number)
				assert.Equal(t, tc.expected, result, "Test case: %s", tc.name)
			})
		}
	})
}
