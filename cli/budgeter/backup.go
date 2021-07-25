package budgeter

import (
	"io"
	"os"

	_ "embed"
)

const backupName = "backup"

//go:embed backupUsage.txt
var backupUsage string

func backup(c *config) int {
	if len(c.args) != 1 {
		c.log.Printf("%s only takes one argument", backupName)
		c.log.Println()
		c.log.Println(backupUsage)
		return 1
	}

	dbFile, err := os.Open(c.dbPath)
	if err != nil {
		c.log.Printf("error opening \"%s\" to read: %v", c.dbPath, err)
		return 1
	}
	// TODO: make some consideration for the case where a file is already present.
	// Consider writing to a temp file first or something
	targetPath := c.args[0]
	target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.log.Printf("error opening \"%s\" to write: %v", targetPath, err)
		return 1
	}
	_, err = io.Copy(target, dbFile)
	if err != nil {
		c.log.Printf("error writing backup to \"%s\": %v", targetPath, err)
		return 1
	}
	return 0
}
