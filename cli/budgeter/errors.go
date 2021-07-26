package budgeter

import (
	"strings"
)

func (c *CLI) logParsingErr(err error) {
	args := strings.Join(c.args, " ")
	c.Log.Printf("could not parse args \"%s\": %v", args, err)
}
