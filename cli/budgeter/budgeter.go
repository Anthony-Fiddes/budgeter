package budgeter

import (
	"log"
	"os"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

//go:embed usage.txt
var usage string

type Tabler interface {
	Insert(transaction.Transaction) error
	Search(query string, limit int) (*transaction.Rows, error)
	Total() (int, error)
}

type CLI struct {
	args []string
	// DBPath is the filepath for the datastore being used. It does not have a
	// default.
	DBPath string
	// Log is used by CLI to log errors. By default, it writes to stderr with no
	// date prefix.
	Log *log.Logger
	// Transactions is a Transactions table, it allows the CLI app to interact
	// with a store of transactions.
	Transactions Tabler
}

// A command performs a budgeting action using the configured CLI. It returns an error code.
type command func(c *CLI) int

// Run runs the budgeter CLI with the given arguments.
//
// Run returns an error code. 1 is an error, and 0 means success.
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
		wipeName:   wipe,
		recentName: recent,
	}

	alias := args[1]
	cmdArgs := args[2:]
	cmd, ok := commands[alias]
	if !ok {
		c.Log.Printf("command \"%s\" does not exist", alias)
		c.Log.Println()
		c.Log.Println(usage)
		return 1
	}
	c.args = cmdArgs
	return cmd(c)
}
