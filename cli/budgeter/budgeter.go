package budgeter

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

//go:embed usage.txt
var usage string

type Table interface {
	Insert(transaction.Transaction) error
	RangeTotal(start, end time.Time) (int, error)
	Remove(transactionID int) error
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
	// Config is a store where CLI can persist data in a key, value format.
	Config Store
	// DBPath is the filepath for the datastore being used. It does not have a
	// default, so it must be set.
	DBPath string
	// Err is used by CLI to log errors. By default, it writes to stderr with no
	// date prefix.
	Err *log.Logger
	// Transactions is a Transactions table, it allows the CLI app to interact
	// with a store of transactions. It does not have a default, so it must be set.
	Transactions Table
}

type command interface {
	Name() string
	Run(*CLI) int
}

// A commandFunc performs a budgeting action using the configured CLI. It returns an
// error code. A nonzero code is an error.
type commandFunc func(c *CLI) int

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
	if c.Err == nil {
		c.Err = log.New(os.Stderr, "", 0)
	}

	if len(args) < 2 {
		c.Err.Println(usage)
		return 1
	}

	commands := map[string]commandFunc{
		backupName: backup,
		exportName: export,
		ingestName: ingest,
		limitName:  limit,
		wipeName:   wipe,
		removeName: remove,
		reportName: report,
	}

	alias := args[1]
	c.args = args[2:]
	cmd, ok := commands[alias]
	if !ok {
		cmds := []command{&add{}, &recent{}}
		for _, cmd := range cmds {
			if cmd.Name() == alias {
				return cmd.Run(c)
			}
		}
		c.Err.Printf("command \"%s\" does not exist", alias)
		c.Err.Println()
		c.Err.Println(usage)
		return 1
	}
	return cmd(c)
}

func getFlagset(commandName string) *flag.FlagSet {
	fs := flag.NewFlagSet(commandName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}
