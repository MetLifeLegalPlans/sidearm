package cmd

import (
	"sidearm/channels"

	"github.com/spf13/cobra"
)

// standaloneCmd represents the standalone command
var standaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Starts a standalone sidearm instance",
	Long: `Creates and connects a sidearm server and client instance internally.
Designed to be used on a single machine.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		go clientCmd.Run(cmd, args)
		go serverCmd.Run(cmd, args)
		go dashboardCmd.Run(cmd, args)
		<-channels.Running
	},
}

func init() {
	rootCmd.AddCommand(standaloneCmd)
}
