package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github.com/MetLifeLegalPlans/sidearm",
	Short: "Distributed HTTP load generator",
	Long: `A distributed HTTP load generator for the rest of us. Inspired by https://github.com/rakyll/hey

Sidearm is a tool that allows you to use an arbitrary number of
arbitrarily located hosts to generate nuanced HTTP load in a flexible
but opinionated way.

It can also be used in standalone mode on a single host.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable response status output")
	// rootCmd.PersistentFlags().String("configFile", "config.yml", "The spec file for your usecase")

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sidearm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
