package common

import "fmt"

// GetJCLIPluginPath returns the path of a jcli plugin
func GetJCLIPluginPath(userHome, name string, binary bool) string {
	suffix := ".yaml"
	if binary {
		suffix = ""
	}
	return fmt.Sprintf("%s/.jenkins-cli/plugins/%s%s", userHome, name, suffix)
}
