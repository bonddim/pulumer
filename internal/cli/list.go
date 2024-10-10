package cli

import (
	"fmt"
)

// List lists all stacks in the working directory.
func (c *CLI) List() error {
	if err := c.getStacks(); err != nil {
		return err
	}

	if len(c.stacks) == 0 {
		fmt.Println("No stacks found")
		return nil
	}

	fmt.Println("Stacks:")
	for _, stack := range c.stacks {
		fmt.Printf("  - %s\n", stack.Id)
	}

	return nil
}
