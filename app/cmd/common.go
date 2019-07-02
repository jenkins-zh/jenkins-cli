package cmd

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// OutputOption represent the format of output
type OutputOption struct {
	Format string
}

// BatchOption represent the options for a batch operation
type BatchOption struct {
	Batch bool
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

func Format(obj interface{}, format string) (data []byte, err error) {
	if format == JsonOutputFormat {
		return json.MarshalIndent(obj, "", "  ")
	} else if format == YAMLOutputFormat {
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", format)
}
