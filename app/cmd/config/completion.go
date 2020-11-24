package config

import (
	"github.com/spf13/cobra"
	"strings"
)

// ValidPluginNames returns the valid plugin name list
func ValidPluginNames(cmd *cobra.Command, args []string, prefix string) (pluginNames []string, directive cobra.ShellCompDirective) {
	directive = cobra.ShellCompDirectiveNoFileComp
	if plugins, err := findPlugins(); err == nil {
		pluginNames = make([]string, 0)
		for i := range plugins {
			plugin := plugins[i]
			name := plugin.Use

			switch cmd.Use {
			case "install":
				if plugin.Installed {
					continue
				}
			case "uninstall":
				if !plugin.Installed {
					continue
				}
			}

			duplicated := false
			for j := range args {
				if name == args[j] {
					duplicated = true
					break
				}
			}

			if !duplicated && strings.HasPrefix(name, prefix) {
				pluginNames = append(pluginNames, name)
			}
		}
	}
	return
}
