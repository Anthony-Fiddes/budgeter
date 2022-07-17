package budgeter

import (
	"fmt"

	_ "embed"
)

//go:embed setUsage.txt
var setUsage string

type set struct{}

func (s *set) Name() string {
	return "set"
}

func (s *set) Run(c *CLI) int {
	if len(c.args) > 2 {
		c.err.Printf("%s takes at most 2 arguments", s.Name())
		c.err.Println()
		c.err.Println(setUsage)
		return 1
	}

	var err error
	switch numArgs := len(c.args); numArgs {
	case 0:
		entries, e := c.Config.GetAll()
		err = e
		for key, val := range entries {
			fmt.Fprintf(c.Out, "%s: %s\n", key, val)
		}
	case 1:
		key := c.args[0]
		val, e := c.Config.Get(key)
		err = e
		if err == nil {
			fmt.Fprintf(c.Out, "%s: %s\n", key, val)
		}
	case 2:
		err = c.Config.Put(c.args[0], c.args[1])
	}
	if err != nil {
		c.err.Printf("configuration error: %v", err)
		return 1
	}
	return 0
}
