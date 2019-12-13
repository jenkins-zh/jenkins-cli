package health

// CommandHealth is the interface for register a command checker
type CommandHealth interface {
	Check() error
}

// CheckRegister is the register container
type CheckRegister struct {
	Member map[string]CommandHealth
}

// Register can register a command and function
func (c *CheckRegister) Register(path string, health CommandHealth) {
	c.Member[path] = health
}
