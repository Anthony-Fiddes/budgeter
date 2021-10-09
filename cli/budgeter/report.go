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

const reportName = "report"

//go:embed reportUsage.txt
var reportUsage string

// report tells the user how much they've spent over the last few months.
func report(c *CLI) int {
	// defaultReportMonths determines how many months to query for when calling
	// the report command.
	const defaultReportMonths = 6

	if len(c.args) != 0 {
		c.Log.Printf("%s takes no arguments", reportName)
		c.Log.Println()
		c.Log.Println(reportUsage)
		return 1
	}

	type total struct {
		month  time.Month
		amount int
	}

	sPrintTotals := func(totals []total) string {
		buf := &bytes.Buffer{}
		// this tab writer uses the same settings as tabby, except obviously for
		// where it writes to.
		tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
		tab := tabby.NewCustom(tw)
		tab.AddHeader("Month", "Spent")
		for _, t := range totals {
			tab.AddLine(t.month.String(), transaction.Dollars(t.amount))
		}
		tab.Print()
		return buf.String()
	}

	var totals []total
	start := month.Start(time.Now().UTC())
	start = month.Add(start, -defaultReportMonths+1)
	for i := 0; i < defaultReportMonths; i++ {
		end := month.End(start)
		amount, err := c.Transactions.RangeTotal(start, end)
		if err != nil {
			c.Log.Printf("could not get totals for all of the requested months: %v", err)
			c.Log.Println("correctly collected totals: ")
			c.Log.Println(sPrintTotals(totals))
			return 1
		}
		totals = append(totals, total{month: start.Month(), amount: amount})
		start = month.Add(start, 1)
	}

	fmt.Println(sPrintTotals(totals))

	return 0
}
