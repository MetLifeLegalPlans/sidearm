package cmd

import (
	"fmt"
	"sidearm/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows your configuration after parsing",
	Long:  `Parses your config from YAML and then prints out the resolved version.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load(args[0])
		if err != nil {
			panic(err)
		}

		dumped, _ := yaml.Marshal(conf)
		fmt.Println(string(dumped))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
