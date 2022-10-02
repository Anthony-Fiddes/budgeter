package budgeter

import (
	"fmt"
	"io"
	"os"

	_ "embed"
)

type backup struct {
	DBPath string
}

func newBackup(c *CLI) *backup {
	result := backup{}
	result.DBPath = c.DBPath
	return &result
}

func (b backup) Name() string {
	return "backup"
}

//go:embed backupUsage.txt
var backupUsage string

func (b backup) Usage() string {
	return backupUsage
}

func (b backup) Run(cmdArgs []string) error {
	if len(cmdArgs) != 1 {
		return fmt.Errorf("%s only takes one argument", b.Name())
	}

	dbFile, err := os.Open(b.DBPath)
	if err != nil {
		return fmt.Errorf("error opening \"%s\" to read: %w", b.DBPath, err)
	}
	// TODO: make some consideration for the case where a file is already present.
	// Consider writing to a temp file first or something
	targetPath := cmdArgs[0]
	target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening \"%s\" to write: %w", targetPath, err)
	}
	_, err = io.Copy(target, dbFile)
	if err != nil {
		return fmt.Errorf("error writing backup to \"%s\": %w", targetPath, err)
	}
	return nil
}
