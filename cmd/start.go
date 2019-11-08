package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var folder string
var port int

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&folder, "directory", "d", "", "Path of the web folder")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",

	Args: cobra.MinimumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listening on port: ", port)
		fileHandler := http.FileServer(http.Dir(folder))
		err := http.ListenAndServe(fmt.Sprint(":", port), fileHandler)
		if err != nil {
			panic(err)
		}
	},
}
