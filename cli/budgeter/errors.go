package budgeter

import (
	"strings"
)

func (c *config) logParsingErr(err error) {
	args := strings.Join(c.args, " ")
	c.log.Printf("could not parse args \"%s\": %v", args, err)
}
