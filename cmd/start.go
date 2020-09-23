package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var folder string
var filePath string

var port int

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&folder, "directory", "d", "", "Path of the web folder")
	startCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path of a file (for hosting a single file)")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",

	Args: cobra.MinimumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) {

		if filePath != "" {
		}

		fileHandler := http.FileServer(http.Dir(folder))

		ips, err := externalIPs()
		if err != nil {
			panic(err)
		}

		fmt.Println("Http server is running at: ")
		if len(ips) == 0 {
			url := "http://localhost:" + strconv.Itoa(port)
			fmt.Println(url)
		}
		for _, ip := range ips {
			url := "http://" + ip + ":" + strconv.Itoa(port)
			fmt.Println(url)
		}

		err = http.ListenAndServe(fmt.Sprint(":", port), fileHandler)
		fmt.Println("bye~")
		if err != nil {
			panic(err)
		}
	},
}
