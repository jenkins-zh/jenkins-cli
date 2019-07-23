package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// OutputOption represent the format of output
type OutputOption struct {
	Format string
}

type FormatOutput interface {
	Output(obj interface{}, format string) (data []byte, err error)
}

const (
	JsonOutputFormat  string = "json"
	YAMLOutputFormat  string = "yaml"
	TableOutputFormat string = "table"
)

func (o *OutputOption) Output(obj interface{}) (data []byte, err error) {
	switch o.Format {
	case JsonOutputFormat:
		return json.MarshalIndent(obj, "", "  ")
	case YAMLOutputFormat:
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", o.Format)
}

func (o *OutputOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Format, "output", "o", "table", "Format the output (default 'json')")
}

func Format(obj interface{}, format string) (data []byte, err error) {
	if format == JsonOutputFormat {
		return json.MarshalIndent(obj, "", "  ")
	} else if format == YAMLOutputFormat {
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", format)
}

// BatchOption represent the options for a batch operation
type BatchOption struct {
	Batch bool
}

// Confirm prompte user if they really want to do this
func (b *BatchOption) Confirm(message string) bool {
	if !b.Batch {
		confirm := false
		prompt := &survey.Confirm{
			Message: message,
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			return false
		}
	}

	return true
}

func (b *BatchOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

// WatchOption for the resources which can be watched
type WatchOption struct {
	Watch    bool
	Interval int
}

// InteractiveOption allow user to choose whether the mode is interactive
type InteractiveOption struct {
	Interactive bool
}

// SetFlag set the option flag to this cmd
func (b *InteractiveOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.Interactive, "interactive", "i", false, "Interactive mode")
}
