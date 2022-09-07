package cmd

import (
	"github.com/spf13/cobra"
	"sidearm/channels"
	"sidearm/client"
	"sidearm/config"

	"sync"
)

var clientCmd = &cobra.Command{
	Use:   "client configFile",
	Short: "Starts sidearm in client/worker mode",
	Long:  `Connects to a sidearm server instance and executes tasks received from it.`,
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
			client.Entrypoint(conf, quiet)
		}()

		<-channels.Interrupt
		close(channels.Running)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
