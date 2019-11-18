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

// StandardErrorMessage is generic to the command in use
func StandardErrorMessage(err error) (msg string, ok bool) {
	ok = true
	switch t := err.(type) {
	case url.InvalidHostError:
		msg = t.Error()
	case *url.Error:
		switch {
		case strings.Contains(t.Err.Error(), "connection refused"):
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			msg = fmt.Sprintf("The connection to the server %s was refused - did you specify the right host or port?", host)
		default:
			msg = fmt.Sprintf("Unable to connect to the server: %v", t.Err)
		}
	case *os.PathError:
		msg = fmt.Sprintf("error: %s %s: %s", t.Op, t.Path, t.Err)
	default:
		ok = false
	}
	return
}
