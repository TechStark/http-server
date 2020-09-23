package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "http-server",
	Short: "A simple command-line HTTP server for static files",
}

// Execute is the command line entry
func Execute() error {
	return rootCmd.Execute()
}
