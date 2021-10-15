package budgeter

import (
	"strings"
)

func (c *CLI) logParsingErr(err error) {
	args := strings.Join(c.args, " ")
	c.Err.Printf("could not parse args `%s`: %v", args, err)
}
