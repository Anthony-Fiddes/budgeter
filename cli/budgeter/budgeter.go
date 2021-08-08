package budgeter

import (
	"log"
	"os"
	"time"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

//go:embed usage.txt
var usage string

// period is an enum representing the lengths of time that budgeter allows
type period int

const (
	unknown period = iota
	day
	week
	month
)

func (p period) Unknown() bool {
	if p <= 0 || p > month {
		return true
	}
	return false
}

func (p period) String() string {
	if p < 0 || p > month {
		return "unknown"
	}
	return [...]string{"unknown", "day", "week", "month"}[int(p)]
}

var periods = map[string]period{
	unknown.String(): unknown,
	day.String():     day,
	week.String():    week,
	month.String():   month,
}

func getPeriod(s string) period {
	p, ok := periods[s]
	if !ok {
		return unknown
	}
	return p
}

type Table interface {
	Insert(transaction.Transaction) error
	Range(start, end time.Time, limit int) (*transaction.Rows, error)
	Search(query string, limit int) (*transaction.Rows, error)
	Total() (int, error)
}

type Store interface {
	// Put puts a value into the Store. If it is already present, it's overwritten.
	Put(Key, Value string) error
	// Get gets a value from the Store. If the value is not present, "" is
	// returned with a nil error.
	Get(Key string) (string, error)
}

type CLI struct {
	args []string
	// Config a store where CLI can persist data in a key, value format.
	Config Store
	// DBPath is the filepath for the datastore being used. It does not have a
	// default, i.e. it must be set.
	DBPath string
	// Log is used by CLI to log errors. By default, it writes to stderr with no
	// date prefix.
	Log *log.Logger
	// Transactions is a Transactions table, it allows the CLI app to interact
	// with a store of transactions.
	Transactions Table
}

// A command performs a budgeting action using the configured CLI. It returns an
// error code. A nonzero code is an error.
type command func(c *CLI) int

// Run runs the budgeter CLI with the given arguments.
//
// Run returns an error code. A nonzero code is an error, and 0 means success.
func (c *CLI) Run(args []string) int {
	if c.DBPath == "" {
		panic("budgeter: DBPath must be set on CLI")
	}
	if c.Transactions == nil {
		panic("budgeter: Transactions must be set on CLI")
	}
	if c.Log == nil {
		c.Log = log.New(os.Stderr, "", 0)
	}

	if len(args) < 2 {
		c.Log.Println(usage)
		return 1
	}

	commands := map[string]command{
		addName:    add,
		backupName: backup,
		exportName: export,
		ingestName: ingest,
		limitName:  limit,
		wipeName:   wipe,
		recentName: recent,
	}

	alias := args[1]
	c.args = args[2:]
	cmd, ok := commands[alias]
	if !ok {
		c.Log.Printf("command \"%s\" does not exist", alias)
		c.Log.Println()
		c.Log.Println(usage)
		return 1
	}
	return cmd(c)
}