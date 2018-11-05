package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type args struct {
	serveDir    string
	indexHTML   string
	port        int
	interactive bool

	encryption sslArgs
}

type sslArgs struct {
	useEncryption bool
	certFile      string
	keyFile       string
}

func readArgs() args {
	serveDir := flag.String("serveDir", "www", "directory to serve")
	indexHTML := flag.String("index", "www/index.html", "file to use as index html")
	port := flag.Int("port", 8080, "port")
	interactive := flag.Bool("interactive", false, "wait for user input to shut down")

	certFile := flag.String("certFile", "", "certificate file, to use encryption, a key file is also required")
	keyFile := flag.String("keyFile", "", "key file, to use encryption, a cert file is also required")

	flag.Parse()

	return args{
		serveDir:    *serveDir,
		indexHTML:   *indexHTML,
		port:        *port,
		interactive: *interactive,
		encryption: sslArgs{
			useEncryption: len(*certFile) > 0 && len(*keyFile) > 0,
			certFile:      *certFile,
			keyFile:       *keyFile,
		},
	}
}

// checkSSLArgs verifies that the given certificate file and key file do exists.
// It does not check if these files contain the required information.
func checkSSLArgs(args sslArgs) error {
	_, err := os.Stat(args.certFile)
	if os.IsNotExist(err) {
		return errors.New("given certificate file " + args.certFile + " was found found")
	}
	_, err = os.Stat(args.keyFile)
	if os.IsNotExist(err) {
		return errors.New("given key file " + args.keyFile + " was found found")
	}
	return nil
}

func main() {
	args := readArgs()

	fmt.Printf("Starting server on %d, interactive %t\n", args.port, args.interactive)
	fmt.Printf("Serving directory: %s\n", args.serveDir)
	fmt.Printf("Serving index html: %s\n", args.indexHTML)

	if args.encryption.useEncryption {
		fmt.Printf("Using certificate: %s\n", args.encryption.certFile)
		fmt.Printf("Using key: %s\n", args.encryption.keyFile)
		err := checkSSLArgs(args.encryption)
		if err != nil {
			os.Exit(1)
		}
	}

	app := newFileServer(args.serveDir, args.indexHTML)
	app.configureRoutes()

	fmt.Printf("Server listening on port: %d\n", args.port)

	app.startServer(args.interactive, args.port, args.encryption)
	if args.interactive {
		fmt.Println("Server started, hit Enter-key to close")
		fmt.Scanln()
		fmt.Println("Shuting down...")
	} else {
		if args.encryption.useEncryption {
			app.listenAndServeSSL(args.port, args.encryption.certFile, args.encryption.keyFile)
		} else {
			app.listenAndServe(args.port)
		}
	}
}

func (app fileServer) startServer(inBackground bool, port int, encryption sslArgs) {
	if inBackground {
		if encryption.useEncryption {
			go app.listenAndServeSSL(port, encryption.certFile, encryption.keyFile)
		} else {
			go app.listenAndServe(port)
		}
	} else {
		if encryption.useEncryption {
			app.listenAndServeSSL(port, encryption.certFile, encryption.keyFile)
		} else {
			app.listenAndServe(port)
		}
	}
}
