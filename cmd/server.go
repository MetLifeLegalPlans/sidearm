package cmd

import (
	"sync"

	"github.com/spf13/cobra"

	"github.com/MetLifeLegalPlans/sidearm/channels"
	"github.com/MetLifeLegalPlans/sidearm/config"
	"github.com/MetLifeLegalPlans/sidearm/server"
)

var serverCmd = &cobra.Command{
	Use:   "server configFile",
	Short: "Starts sidearm in server/distributor mode",
	Long:  `Creates a sidearm server instance that generates tasks for workers.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load(args[0])
		if err != nil {
			panic(err)
		}

		if !conf.QueueConfig.Enabled() {
			panic("No socket configuration available for tasks, exiting")
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			server.Entrypoint(conf)
		}()

		<-channels.Interrupt
		close(channels.Running)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
