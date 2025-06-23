package mapper

import (
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWithdrawal_ToDomainModel(t *testing.T) {
	mapper := withdrawal{}

	tests := []struct {
		name     string
		input    db.WithdrawalsRow
		expected model.Withdrawal
	}{
		{
			name: "basic withdrawal",
			input: db.WithdrawalsRow{
				OrderNumber: "12345",
				Withdrawn:   createNumeric(t, "100.50"),
				ProcessedAt: createTimestamp(t, "2023-12-01T10:30:00Z"),
			},
			expected: model.Withdrawal{
				Order:       "12345",
				Sum:         100.50,
				ProcessedAt: parseTime(t, "2023-12-01T10:30:00Z"),
			},
		},
		{
			name: "zero amount withdrawal",
			input: db.WithdrawalsRow{
				OrderNumber: "00000",
				Withdrawn:   createNumeric(t, "0"),
				ProcessedAt: createTimestamp(t, "2023-11-15T14:45:30Z"),
			},
			expected: model.Withdrawal{
				Order:       "00000",
				Sum:         0.0,
				ProcessedAt: parseTime(t, "2023-11-15T14:45:30Z"),
			},
		},
		{
			name: "large amount withdrawal",
			input: db.WithdrawalsRow{
				OrderNumber: "999999999",
				Withdrawn:   createNumeric(t, "9999.99"),
				ProcessedAt: createTimestamp(t, "2023-01-01T00:00:00Z"),
			},
			expected: model.Withdrawal{
				Order:       "999999999",
				Sum:         9999.99,
				ProcessedAt: parseTime(t, "2023-01-01T00:00:00Z"),
			},
		},
		{
			name: "decimal amount withdrawal",
			input: db.WithdrawalsRow{
				OrderNumber: "ABC123",
				Withdrawn:   createNumeric(t, "45.67"),
				ProcessedAt: createTimestamp(t, "2023-06-15T12:00:00Z"),
			},
			expected: model.Withdrawal{
				Order:       "ABC123",
				Sum:         45.67,
				ProcessedAt: parseTime(t, "2023-06-15T12:00:00Z"),
			},
		},
		{
			name: "small decimal amount",
			input: db.WithdrawalsRow{
				OrderNumber: "MICRO001",
				Withdrawn:   createNumeric(t, "0.01"),
				ProcessedAt: createTimestamp(t, "2023-08-20T16:30:45Z"),
			},
			expected: model.Withdrawal{
				Order:       "MICRO001",
				Sum:         0.01,
				ProcessedAt: parseTime(t, "2023-08-20T16:30:45Z"),
			},
		},
		{
			name: "empty order number",
			input: db.WithdrawalsRow{
				OrderNumber: "",
				Withdrawn:   createNumeric(t, "25.00"),
				ProcessedAt: createTimestamp(t, "2023-03-10T09:15:00Z"),
			},
			expected: model.Withdrawal{
				Order:       "",
				Sum:         25.00,
				ProcessedAt: parseTime(t, "2023-03-10T09:15:00Z"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToDomainModel(tt.input)

			assert.Equal(t, tt.expected.Order, result.Order)
			assert.Equal(t, tt.expected.Sum, result.Sum)
			assert.True(t, tt.expected.ProcessedAt.Equal(result.ProcessedAt),
				"Expected time %v, got %v", tt.expected.ProcessedAt, result.ProcessedAt)
		})
	}
}

func TestWithdrawal_ToDomainModel_WithInvalidNumeric(t *testing.T) {
	mapper := withdrawal{}

	// Test with uninitialized numeric (should not panic and return 0)
	input := db.WithdrawalsRow{
		OrderNumber: "TEST001",
		Withdrawn:   pgtype.Numeric{},
		ProcessedAt: createTimestamp(t, "2023-12-01T10:30:00Z"),
	}

	result := mapper.ToDomainModel(input)

	assert.Equal(t, "TEST001", result.Order)
	assert.Equal(t, 0.0, result.Sum) // Should default to 0 when conversion fails
	assert.Equal(t, parseTime(t, "2023-12-01T10:30:00Z"), result.ProcessedAt)
}

func TestWithdrawal_ToDomainModel_WithDifferentTimeZones(t *testing.T) {
	mapper := withdrawal{}

	tests := []struct {
		name     string
		timeStr  string
		orderNum string
	}{
		{
			name:     "UTC time",
			timeStr:  "2023-12-01T10:30:00Z",
			orderNum: "UTC001",
		},
		{
			name:     "different UTC time",
			timeStr:  "2024-01-15T23:59:59Z",
			orderNum: "UTC002",
		},
		{
			name:     "beginning of epoch",
			timeStr:  "1970-01-01T00:00:00Z",
			orderNum: "EPOCH001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := db.WithdrawalsRow{
				OrderNumber: tt.orderNum,
				Withdrawn:   createNumeric(t, "100.00"),
				ProcessedAt: createTimestamp(t, tt.timeStr),
			}

			result := mapper.ToDomainModel(input)

			assert.Equal(t, tt.orderNum, result.Order)
			assert.Equal(t, 100.00, result.Sum)
			assert.Equal(t, parseTime(t, tt.timeStr), result.ProcessedAt)
		})
	}
}

func TestWithdrawal_SingletonInstance(t *testing.T) {
	assert.NotNil(t, Withdrawal)
	assert.IsType(t, withdrawal{}, Withdrawal)
}

func TestWithdrawal_ToDomainModel_FieldMapping(t *testing.T) {
	mapper := withdrawal{}

	// Test that all fields are properly mapped
	input := db.WithdrawalsRow{
		OrderNumber: "FIELD_TEST",
		Withdrawn:   createNumeric(t, "123.45"),
		ProcessedAt: createTimestamp(t, "2023-07-20T14:30:00Z"),
	}

	result := mapper.ToDomainModel(input)

	// Verify each field is mapped correctly
	assert.Equal(t, input.OrderNumber, result.Order, "OrderNumber should map to Order")

	expectedSum, _ := input.Withdrawn.Float64Value()
	assert.Equal(t, expectedSum.Float64, result.Sum, "Withdrawn should map to Sum")

	assert.Equal(t, input.ProcessedAt.Time, result.ProcessedAt, "ProcessedAt should map correctly")
}

// Helper functions

func createTimestamp(t *testing.T, timeStr string) pgtype.Timestamptz {
	parsedTime := parseTime(t, timeStr)
	return pgtype.Timestamptz{
		Time:  parsedTime,
		Valid: true,
	}
}

func parseTime(t *testing.T, timeStr string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err, "Failed to parse time: %s", timeStr)
	return parsedTime
}
