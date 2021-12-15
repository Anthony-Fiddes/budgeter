package transaction_test

import (
	"testing"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

func TestCents(t *testing.T) {
	tests := []struct {
		input    string
		expected transaction.Cent
	}{
		{
			input:    "5",
			expected: 500,
		},
		{
			input:    "5.00",
			expected: 500,
		},
		{
			input:    "5.30",
			expected: 530,
		},
		{
			input:    "5.3",
			expected: 530,
		},
		{
			input:    "500",
			expected: 50000,
		},
		{
			input:    "500,000",
			expected: 50000000,
		},
		{
			input:    "500,000.00",
			expected: 50000000,
		},
		{
			input:    "-5",
			expected: -500,
		},
		{
			input:    "-5.00",
			expected: -500,
		},
		{
			input:    "-5.30",
			expected: -530,
		},
		{
			input:    "-5.3",
			expected: -530,
		},
		{
			input:    "-500",
			expected: -50000,
		},
		{
			input:    "-500,000",
			expected: -50000000,
		},
		{
			input:    "-500,000.00",
			expected: -50000000,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := transaction.GetCents(test.input)
			if err != nil {
				t.Fatalf("err: %s\ntest: %+v", err, test)
			}
			if result != test.expected {
				t.Fatalf("received %d but expected %d", result, test.expected)
			}
		})
	}
}
