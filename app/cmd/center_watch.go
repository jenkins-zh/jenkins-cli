package cmd

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// CenterWatchOption as the options of watch command
type CenterWatchOption struct {
	WatchOption

	RoundTripper  http.RoundTripper
	CeneterStatus string
}

var centerWatchOption CenterWatchOption

func init() {
	centerCmd.AddCommand(centerWatchCmd)
	centerWatchOption.SetFlag(centerWatchCmd)
}

var centerWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch your update center status",
	Long:  `Watch your update center status`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		printJenkinsStatus(jenkins, centerWatchOption.RoundTripper)

		for ; centerWatchOption.Count >= 0; centerWatchOption.Count-- {
			printUpdateCenter(jenkins, centerOption.RoundTripper)

			time.Sleep(time.Duration(centerOption.Interval) * time.Second)
		}
	},
}
