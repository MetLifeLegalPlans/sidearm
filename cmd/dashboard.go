package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Gives you live stats about your sidearm session",
	Long: `Creates a task ventilator socket that receives reports from your workers
and calculates statistics from them. This is then displayed in a TUI.

NOT RECOMMENDED FOR USAGE IN DOCKER.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Coming soon!")
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
