package mapper

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestMapper_NumericToFloat(t *testing.T) {
	mapper := mapper{}

	tests := []struct {
		name     string
		input    pgtype.Numeric
		expected float64
	}{
		{
			name:     "zero value",
			input:    createNumeric(t, "0"),
			expected: 0.0,
		},
		{
			name:     "positive integer",
			input:    createNumeric(t, "123"),
			expected: 123.0,
		},
		{
			name:     "negative integer",
			input:    createNumeric(t, "-456"),
			expected: -456.0,
		},
		{
			name:     "positive decimal",
			input:    createNumeric(t, "123.45"),
			expected: 123.45,
		},
		{
			name:     "negative decimal",
			input:    createNumeric(t, "-78.90"),
			expected: -78.90,
		},
		{
			name:     "small decimal",
			input:    createNumeric(t, "0.01"),
			expected: 0.01,
		},
		{
			name:     "large number",
			input:    createNumeric(t, "999999.99"),
			expected: 999999.99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.NumericToFloat(tt.input)
			require.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestMapper_NumericToFloat_WithInvalidNumeric(t *testing.T) {
	mapper := mapper{}

	uninitializedNumeric := pgtype.Numeric{}
	result := mapper.NumericToFloat(uninitializedNumeric)

	require.NotNil(t, result, "Result should not be nil even for invalid numeric")
	assert.Equal(t, 0.0, *result)
}

func TestMapper_FloatToNumeric(t *testing.T) {
	mapper := mapper{}

	tests := []struct {
		name     string
		input    *float64
		expected float64 // We'll convert back to verify
	}{
		{
			name:     "zero value",
			input:    floatPtr(0.0),
			expected: 0.0,
		},
		{
			name:     "positive integer",
			input:    floatPtr(123.0),
			expected: 123.0,
		},
		{
			name:     "negative integer",
			input:    floatPtr(-456.0),
			expected: -456.0,
		},
		{
			name:     "positive decimal",
			input:    floatPtr(123.45),
			expected: 123.45,
		},
		{
			name:     "negative decimal",
			input:    floatPtr(-78.90),
			expected: -78.90,
		},
		{
			name:     "small decimal",
			input:    floatPtr(0.01),
			expected: 0.01,
		},
		{
			name:     "rounds to 2 decimal places",
			input:    floatPtr(123.456),
			expected: 123.46, // Should round to 2 decimal places
		},
		{
			name:     "rounds down",
			input:    floatPtr(123.454),
			expected: 123.45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.FloatToNumeric(tt.input)
			convertedBack := mapper.NumericToFloat(result)
			require.NotNil(t, convertedBack)
			assert.InDelta(t, tt.expected, *convertedBack, 0.001)
		})
	}
}

func TestMapper_RoundTripConversion(t *testing.T) {
	mapper := mapper{}

	testValues := []float64{
		0.0,
		123.45,
		-67.89,
		999.99,
		0.01,
		-0.01,
		1000000.50,
	}

	for _, value := range testValues {
		t.Run(formatFloat(value), func(t *testing.T) {
			// Convert float to numeric and back
			numeric := mapper.FloatToNumeric(&value)
			convertedBack := mapper.NumericToFloat(numeric)

			require.NotNil(t, convertedBack)
			assert.InDelta(t, value, *convertedBack, 0.01)
		})
	}
}

func TestMapper_SingletonInstance(t *testing.T) {
	assert.NotNil(t, Mapper)
	assert.IsType(t, mapper{}, Mapper)
}

// Helper functions

func createNumeric(t *testing.T, value string) pgtype.Numeric {
	var n pgtype.Numeric
	err := n.Scan(value)
	require.NoError(t, err, "Failed to create numeric from string: %s", value)
	return n
}

func floatPtr(f float64) *float64 {
	return &f
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
