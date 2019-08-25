package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// CenterWatchOption as the options of watch command
type CenterWatchOption struct {
	WatchOption

	CeneterStatus string
}

var centerWatchOption CenterWatchOption

func init() {
	centerCmd.AddCommand(centerWatchCmd)
	centerWatchCmd.Flags().IntVarP(&centerWatchOption.Interval, "interval", "i", 1, "Interval of watch")
}

var centerWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch your update center status",
	Long:  `Watch your update center status`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		printJenkinsStatus(jenkins)

		for {
			printUpdateCenter(jenkins)

			time.Sleep(time.Duration(centerOption.Interval) * time.Second)
		}
	},
}
