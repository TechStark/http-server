# http-server
A simple command-line HTTP server

## Usage:
```sh
cd path/to/web-folder/
http-server start
```

## Help Doc:
```sh
% http-server start --help
Start the HTTP server

Usage:
  http-server start [flags]

Flags:
  -d, --directory string   Path of the web folder
  -f, --file string        Path of a file (for sharing a single file)
  -h, --help               help for start
  -p, --port int           Port number (default 8080)
      --qrcode             Show QR code of server URL
      --upload             Allow uploading files via path /upload
```
