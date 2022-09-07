package cmd

import (
	"github.com/spf13/cobra"
	"sidearm/channels"
	"sidearm/config"
	"sidearm/dashboard"

	"sync"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Gives you live stats about your sidearm session",
	Long: `Creates a task ventilator socket that receives reports from your workers
and calculates statistics from them. This is then displayed in a TUI.

NOT RECOMMENDED FOR USAGE IN DOCKER.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load(args[0])
		if err != nil {
			panic(err)
		}

		if !conf.SinkConfig.Enabled() {
			panic("No socket configuration available for the dashboard, exiting")
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			dashboard.Entrypoint(conf)
		}()

		<-channels.Interrupt
		close(channels.Running)
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
