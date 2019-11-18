package helper

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

// CheckErr print a friendly error message
func CheckErr(cmd *cobra.Command, err error) {
	switch {
	case err == nil:
		return
	default:
		switch err := err.(type) {
		default: // for any other error type
			msg, ok := StandardErrorMessage(err)
			if !ok {
				msg = err.Error()
				if !strings.HasPrefix(msg, "error: ") {
					msg = fmt.Sprintf("error: %s", msg)
				}
			}
			cmd.PrintErr(msg)
		}
	}
}

// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func StandardErrorMessage(err error) (string, bool) {
	switch t := err.(type) {
	case *url.Error:
		switch {
		case strings.Contains(t.Err.Error(), "connection refused"):
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			return fmt.Sprintf("The connection to the server %s was refused - did you specify the right host or port?", host), true
		}
		return fmt.Sprintf("Unable to connect to the server: %v", t.Err), true
	case *os.PathError:
		return fmt.Sprintf("error: %s %s: %s", t.Op, t.Path, t.Err), true
	}
	return "", false
}
