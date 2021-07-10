package transaction_test

import (
	"bytes"
	"testing"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

type csvTest struct {
	name         string
	transactions []transaction.Transaction
	text         string
}

func equal(tx, other transaction.Transaction) bool {
	if tx.Entity != other.Entity || tx.Note != other.Note {
		return false
	}
	if tx.DateString() != other.DateString() {
		return false
	}
	if tx.Amount != other.Amount {
		return false
	}
	return true
}

func csvTestData() []csvTest {
	return []csvTest{
		{
			name: "single transaction",
			transactions: []transaction.Transaction{
				{
					Entity: "Apossumtheosis",
					Amount: 400000,
					Date:   -1,
					Note:   "it has begun.",
				},
			},
			text: "12/31/1969,Apossumtheosis,$4000.00,it has begun.\n",
		},
		{
			name: "single negative transaction",
			transactions: []transaction.Transaction{
				{
					Entity: "Apossumtheosis",
					Amount: -400000,
					Date:   -1,
					Note:   "it has begun.",
				},
			},
			text: "12/31/1969,Apossumtheosis,-$4000.00,it has begun.\n",
		},
		{
			name: "single modern transaction",
			transactions: []transaction.Transaction{
				{
					Entity: "Apossumtheosis",
					Amount: 400000,
					Date:   1625784806,
					Note:   "it has begun.",
				},
			},
			text: "7/8/2021,Apossumtheosis,$4000.00,it has begun.\n",
		},
		{
			name: "duplicate modern transactions",
			transactions: []transaction.Transaction{
				{
					Entity: "Apossumtheosis",
					Amount: 400000,
					Date:   1625784806,
					Note:   "it has begun.",
				},
				{
					Entity: "Apossumtheosis",
					Amount: 400000,
					Date:   1625784806,
					Note:   "it has begun.",
				},
			},
			text: "7/8/2021,Apossumtheosis,$4000.00,it has begun.\n" +
				"7/8/2021,Apossumtheosis,$4000.00,it has begun.\n",
		},
	}
}

func TestCSVWriter(t *testing.T) {
	tests := csvTestData()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			cw := transaction.NewCSVWriter(&b)
			if err := cw.WriteAll(test.transactions); err != nil {
				t.Error(err)
			}

			result := b.String()
			expected := test.text
			if result != expected {
				t.Logf("result: %q", result)
				t.Errorf("expected: %q", expected)
			}
		})
	}
}

func TestCSVReader(t *testing.T) {
	tests := csvTestData()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := bytes.NewBufferString(test.text)
			cw := transaction.NewCSVReader(b)
			results, err := cw.ReadAll()
			if err != nil {
				t.Error(err)
			}

			for i := range test.transactions {
				expected := test.transactions[i]
				if i >= len(results) {
					t.Errorf("result has no transaction where there should be %+v", expected)
					continue
				}
				result := results[i]
				if !equal(result, expected) {
					t.Logf("Result: %+v", result)
					t.Errorf("Expected: %+v", expected)
				}
			}
		})
	}
}
