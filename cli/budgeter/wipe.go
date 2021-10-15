package budgeter

import (
	"flag"
	"fmt"
	"os"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
)

const (
	wipeName          = "wipe"
	wipeCancelMessage = "No data deleted."
)

func wipe(c *CLI) int {
	fs := flag.NewFlagSet(wipeName, flag.ContinueOnError)
	fs.SetOutput(c.err.Writer())
	confirmed := fs.Bool("y", false, "Confirms that the user would like to wipe their budgeting information.")
	err := fs.Parse(c.args)
	if err != nil {
		// I'm not logging anything because fs does a good enough job on its own
		// for the wipe command for now.
		return 1
	}
	if len(fs.Args()) > 0 {
		fs.Usage()
		c.err.Println()
		c.err.Printf("%s does not take any arguments", wipeName)
		return 1
	}
	if *confirmed {
		return wipeDB(c)
	}

	fmt.Print("This will delete your budgeting information. Are you sure you want to continue? (y/[n]) ")
	*confirmed, err = inpt.Confirm()
	if err != nil {
		c.err.Println(err)
		return 1
	}
	if !*confirmed {
		return 0
	}
	fmt.Println("Proceeding with deletion...")
	return wipeDB(c)
}

func wipeDB(c *CLI) int {
	if err := os.Remove(c.DBPath); err != nil {
		c.err.Printf("could not wipe database: %v", err)
		return 1
	}
	fmt.Println("Done. All budgeting information deleted.")
	return 0
}
