package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type configUpdateOption struct {
	name  string
	token string
}

func createConfigUpdateCmd() (cmd *cobra.Command) {
	opt := &configUpdateOption{}

	cmd = &cobra.Command{
		Use:               "update",
		Aliases:           []string{"up"},
		Short:             "Update a Jenkins config",
		Example:           "jcli config update --token",
		PreRunE:           opt.preRunE,
		ValidArgsFunction: ValidJenkinsNames,
		RunE:              opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.token, "token", "", "",
		"The token of Jenkins config item")
	return
}

func (o *configUpdateOption) preRunE(_ *cobra.Command, args []string) (err error) {
	if o.token == "" {
		err = fmt.Errorf("no token provided")
	}

	if len(args) > 0 {
		o.name = args[0]
	}
	return
}

func (o *configUpdateOption) runE(_ *cobra.Command, _ []string) (err error) {
	found := false
	for i, cfg := range config.JenkinsServers {
		if cfg.Name == o.name {
			found = true
			config.JenkinsServers[i].Token = o.token
			err = saveConfig()
			break
		}
	}

	if !found {
		err = fmt.Errorf("jenkins '%s' does not exist", o.name)
	}
	return
}
