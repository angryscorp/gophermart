package mapper

import (
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestBalance_NumericToFloat(t *testing.T) {
	mapper := balance{}

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
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBalance_FloatToNumeric(t *testing.T) {
	mapper := balance{}

	tests := []struct {
		name     string
		input    float64
		expected string // To compare the string representation
	}{
		{
			name:     "zero value",
			input:    0.0,
			expected: "0.00",
		},
		{
			name:     "positive integer",
			input:    123.0,
			expected: "123.00",
		},
		{
			name:     "negative integer",
			input:    -456.0,
			expected: "-456.00",
		},
		{
			name:     "positive decimal",
			input:    123.45,
			expected: "123.45",
		},
		{
			name:     "negative decimal",
			input:    -78.90,
			expected: "-78.90",
		},
		{
			name:     "small decimal",
			input:    0.01,
			expected: "0.01",
		},
		{
			name:     "rounds to 2 decimal places",
			input:    123.456,
			expected: "123.46",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.FloatToNumeric(tt.input)
			convertedBack := mapper.NumericToFloat(result)
			expectedFloat, err := parseFloat(tt.expected)
			require.NoError(t, err)
			assert.InDelta(t, expectedFloat, convertedBack, 0.001)
		})
	}
}

func TestBalance_ToDomainModel(t *testing.T) {
	mapper := balance{}

	tests := []struct {
		name     string
		input    db.BalanceRow
		expected model.Balance
	}{
		{
			name: "zero balances",
			input: db.BalanceRow{
				Balance:   createNumeric(t, "0"),
				Withdrawn: createNumeric(t, "0"),
			},
			expected: model.Balance{
				Current:   0.0,
				Withdrawn: 0.0,
			},
		},
		{
			name: "positive balances",
			input: db.BalanceRow{
				Balance:   createNumeric(t, "1234.56"),
				Withdrawn: createNumeric(t, "789.12"),
			},
			expected: model.Balance{
				Current:   1234.56,
				Withdrawn: 789.12,
			},
		},
		{
			name: "negative balance, positive withdrawn",
			input: db.BalanceRow{
				Balance:   createNumeric(t, "-100.50"),
				Withdrawn: createNumeric(t, "250.75"),
			},
			expected: model.Balance{
				Current:   -100.50,
				Withdrawn: 250.75,
			},
		},
		{
			name: "large numbers",
			input: db.BalanceRow{
				Balance:   createNumeric(t, "999999.99"),
				Withdrawn: createNumeric(t, "888888.88"),
			},
			expected: model.Balance{
				Current:   999999.99,
				Withdrawn: 888888.88,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToDomainModel(tt.input)
			assert.Equal(t, tt.expected.Current, result.Current)
			assert.Equal(t, tt.expected.Withdrawn, result.Withdrawn)
		})
	}
}

func TestBalance_RoundTripConversion(t *testing.T) {
	mapper := balance{}

	testValues := []float64{
		0.0,
		123.45,
		-67.89,
		999.99,
		0.01,
		-0.01,
	}

	for _, value := range testValues {
		t.Run(fmt.Sprintf("round_trip_%.2f", value), func(t *testing.T) {
			// Convert float to numeric and back
			numeric := mapper.FloatToNumeric(value)
			convertedBack := mapper.NumericToFloat(numeric)

			// Should be equal within small delta due to floating point precision
			assert.InDelta(t, value, convertedBack, 0.001)
		})
	}
}

func TestBalance_SingletonInstance(t *testing.T) {
	assert.NotNil(t, Balance)
	assert.IsType(t, balance{}, Balance)
}

// Helper functions

func createNumeric(t *testing.T, value string) pgtype.Numeric {
	var n pgtype.Numeric
	err := n.Scan(value)
	require.NoError(t, err, "Failed to create numeric from string: %s", value)
	return n
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
