package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "leader <subcommand>",
	Short: "serves leader",
	Run:   nil,
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringP("config-file", "c", "",
		"Path to the config file (eg ./config.yaml) [Optional]")
}
