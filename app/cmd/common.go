package cmd

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

type OutputOption struct {
	Format string
}

const (
	JsonOutputFormat string = "json"
	YAMLOutputFormat string = "yaml"
)

func Format(obj interface{}, format string) (data []byte, err error) {
	if format == JsonOutputFormat {
		return json.MarshalIndent(obj, "", "  ")
	} else if format == YAMLOutputFormat {
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", format)
}
