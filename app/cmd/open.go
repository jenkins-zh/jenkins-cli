package cmd

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open your Jenkins in the browse",
	Long:  `Open your Jenkins in the browse`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		if jenkins.URL != "" {
			open(jenkins.URL)
		} else {
			log.Fatalf("No URL found with Jenkins %s", jenkins.Name)
		}
	},
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
