package budgeter

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/month"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	"github.com/cheynewallace/tabby"
)

type report struct {
	Transactions Table
}

func newReport(c *CLI) *report {
	return &report{Transactions: c.Transactions}
}

func (r report) Name() string { return "report" }

//go:embed reportUsage.txt
var reportUsage string

func (r report) Usage() string {
	return reportUsage
}

// report tells the user how much they've spent over the last few months.
func (r report) Run(cmdArgs []string) error {
	// defaultReportMonths determines how many months to query for when calling
	// the report command.
	const defaultReportMonths = 6

	type total struct {
		month  time.Month
		amount transaction.Cent
	}

	sPrintTotals := func(totals []total) string {
		buf := &bytes.Buffer{}
		// this tab writer uses the same settings as tabby, except obviously for
		// where it writes to.
		tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
		tab := tabby.NewCustom(tw)
		tab.AddHeader("Month", "Spent")
		for _, t := range totals {
			amtStr := t.amount.String()
			// Align all the amount cells
			if t.amount >= 0 {
				amtStr = " " + amtStr
			}
			tab.AddLine(t.month.String(), amtStr)
		}
		tab.Print()
		return buf.String()
	}

	var totals []total
	start := month.Start(time.Now().UTC())
	start = month.Add(start, -defaultReportMonths+1)
	for i := 0; i < defaultReportMonths; i++ {
		end := month.End(start)
		amount, err := r.Transactions.RangeTotal(start, end)
		if err != nil {
			return fmt.Errorf("could not get totals for all of the requested months: %w", err)
		}
		totals = append(totals, total{month: start.Month(), amount: amount})
		start = month.Add(start, 1)
	}

	fmt.Println(sPrintTotals(totals))

	return nil
}
