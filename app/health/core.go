package health

import "github.com/spf13/cobra"

type CommandHealth interface {
	Check() error
}

type CheckRegister struct {
	Member map[*cobra.Command]CommandHealth
}

func (c *CheckRegister) Register(cmd *cobra.Command, health CommandHealth) {
	c.Member[cmd] = health
}
