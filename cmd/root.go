// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "singularity",
	Short: "singularity Message Broker CLI",
	Long:  `A command line interface for Singularity`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
