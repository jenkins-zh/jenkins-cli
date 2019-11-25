package health

import "github.com/spf13/cobra"

// CommandHealth is the interface for register a command checker
type CommandHealth interface {
	Check() error
}

// CheckRegister is the register container
type CheckRegister struct {
	Member map[*cobra.Command]CommandHealth
}

// Init init the storage
func (c *CheckRegister) Init() {
	c.Member = make(map[*cobra.Command]CommandHealth, 0)
}

// Register can register a command and function
func (c *CheckRegister) Register(cmd *cobra.Command, health CommandHealth) {
	c.Member[cmd] = health
}
