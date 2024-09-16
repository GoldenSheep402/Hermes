package cmd

import (
	"github.com/GoldenSheep402/Hermes/cmd/config"
	"github.com/GoldenSheep402/Hermes/cmd/create"
	"github.com/GoldenSheep402/Hermes/cmd/server"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "jframe",
	SilenceUsage: true,
	Short:        "jframe is a Golang framework with unlimited creativity",
	Example:      "jframe server -c ./config.yaml",
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(create.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
