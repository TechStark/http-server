package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
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

// upload files
var allowUpload bool

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&folder, "directory", "d", "", "Path of the web folder")
	startCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path of a file (for sharing a single file)")
	startCmd.Flags().BoolVar(&showQrCode, "qrcode", false, "Show QR code of server URL")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number")
	startCmd.Flags().BoolVar(&allowUpload, "upload", false, "Allow uploading files via path /upload")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",

	Args: cobra.MinimumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		if filePath != "" {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Disposition", "attachment;filename=\""+filepath.Base(filePath)+"\"")
				http.ServeFile(w, r, filePath)
			})
		} else {
			fileHandler := http.FileServer(http.Dir(folder))
			http.Handle("/", fileHandler)
		}

		if allowUpload {
			http.HandleFunc("/upload", uploadPage)
			http.HandleFunc("/api/upload", uploadFile)
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

func uploadPage(w http.ResponseWriter, r *http.Request) {
	htmlContent := `<!DOCTYPE html>
<html lang="en">
	<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta http-equiv="X-UA-Compatible" content="ie=edge" />
	<title>Upload Files</title>
	</head>
	<body>
	<form enctype="multipart/form-data" action="/api/upload" method="post">
		<input type="file" name="file01" />
		<input type="submit" value="upload" />
	</form>
	</body>
</html>
`
	w.Write([]byte(htmlContent))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// 1. parse input
	r.ParseMultipartForm(10 << 20) // 10 MB

	// 2. retrieve file
	file, header, err := r.FormFile("file01")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("Error retrieving the file"))
		return
	}
	defer file.Close()

	fileName := header.Filename
	fileSize := header.Size
	fmt.Println("=======================================")
	fmt.Printf("Uploading file: %+v\n", fileName)
	fmt.Printf("File : %.2f MB\n", float64(fileSize)/(1<<20))
	fmt.Printf("MIME type: %+v\n", header.Header["Content-Type"])

	// 3. write temporary file on our server
	tempFile, err := ioutil.TempFile("", "http-upload-*")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Error creating temp file"))
		return
	}
	// fmt.Printf("Temp file %+v\n", tempFile.Name())
	defer tempFile.Close()

	if _, err := io.CopyN(tempFile, file, fileSize); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Error saving the file"))
		return
	}

	// 4. move file
	newFileName := fileName
	extension := filepath.Ext(fileName)
	fileNameNoExt := fileName[0 : len(fileName)-len(extension)]
	i := 2
	for {
		if _, err := os.Stat(filepath.Join(folder, newFileName)); err == nil {
			// file exists
			newFileName = fmt.Sprintf("%s-%d%s", fileNameNoExt, i, extension)
			i++
		} else {
			break
		}
	}
	os.Rename(tempFile.Name(), filepath.Join(folder, newFileName))

	// done
	fmt.Printf("Successfully Uploaded File\n\n\n")
	http.Redirect(w, r, "/", 302)
}
