package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

var folder string
var filePath string
var port int
var showQrCode bool

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&folder, "directory", "d", "", "Path of the web folder")
	startCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path of a file (for sharing a single file)")
	startCmd.Flags().BoolVar(&showQrCode, "qrcode", false, "Show QR code of server URL")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",

	Args: cobra.MinimumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		if filePath != "" {
			http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
				res.Header().Add("Content-Disposition", "attachment;filename=\""+filepath.Base(filePath)+"\"")
				http.ServeFile(res, req, filePath)
			})
		} else {
			fileHandler := http.FileServer(http.Dir(folder))
			http.Handle("/", fileHandler)
		}

		ips, err := externalIPs()
		if err != nil {
			panic(err)
		}

		fmt.Println("Http server is running at:")
		fmt.Println("http://localhost:" + strconv.Itoa(port))
		for _, ip := range ips {
			url := "http://" + ip + ":" + strconv.Itoa(port)
			fmt.Println(url)
			if showQrCode {
				qrterminal.Generate(url, qrterminal.L, os.Stdout)
			}
		}

		err = http.ListenAndServe(fmt.Sprint(":", port), nil)
		fmt.Println("bye~")
		if err != nil {
			panic(err)
		}
	},
}
