package main

import (
	"flag"
	"fmt"
)

type args struct {
	serveDir    string
	indexHTML   string
	port        int
	interactive bool
}

func readArgs() args {
	serveDir := flag.String("serveDir", "www", "directory to serve")
	indexHTML := flag.String("index", "www/index.html", "file to use as index html")
	port := flag.Int("port", 8080, "port")
	interactive := flag.Bool("interactive", false, "wait for user input to shut down")

	flag.Parse()

	return args{
		serveDir:    *serveDir,
		indexHTML:   *indexHTML,
		port:        *port,
		interactive: *interactive,
	}
}

func main() {
	args := readArgs()

	fmt.Printf("Starting server on %d, interactive %t\n", args.port, args.interactive)
	fmt.Printf("Serving directory: %s\n", args.serveDir)
	fmt.Printf("Serving index html: %s\n", args.indexHTML)

	app := newFileServer(args.serveDir, args.indexHTML)
	app.configureRoutes()

	fmt.Printf("Server listening on port: %d\n", args.port)

	if args.interactive {
		go app.listenAndServe(args.port)

		fmt.Println("Server started, hit Enter-key to close")
		fmt.Scanln()
		fmt.Println("Shuting down...")
	} else {
		app.listenAndServe(args.port)
	}
}
