package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const timeFormat string = "Mon Jan 2 15:04:05 2006"

type fileServer struct {
	routes map[string]string
}

func (server fileServer) configureRoutes() {
	http.HandleFunc("/", log(server.serve))
}

func (server fileServer) listenAndServe(port int) {
	normPort := fmt.Sprintf(":%d", port)
	http.ListenAndServe(normPort, nil)
}

func (server fileServer) listenAndServeSSL(port int, certFile string, keyFile string) {
	normPort := fmt.Sprintf(":%d", port)
	http.ListenAndServeTLS(normPort, certFile, keyFile, nil)
}

func (server fileServer) serve(w http.ResponseWriter, r *http.Request) {
	requested := r.RequestURI
	local, ok := server.routes[requested]
	if !ok {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	html, err := ioutil.ReadFile(local)
	if err != nil {
		logError(r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(html)
}

// newFileServer new instance of server for static files.
// It constructs all valid routes from walking the www-dir.
func newFileServer(serveDir string, index string) *fileServer {
	r := make(map[string]string)
	filepath.Walk(serveDir, func(p string, info os.FileInfo, _ error) error {
		if p != serveDir {
			p = filepath.ToSlash(p)
			route := strings.TrimPrefix(p, serveDir)
			r[route] = p
		}
		return nil
	})
	r["/"] = index

	server := fileServer{routes: r}
	return &server
}

func log(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		fmt.Printf("%s - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		handler(w, r)
	}
}

func logError(r *http.Request, err string) {
	timestamp := time.Now()
	fmt.Printf("%s - %s - %s\n", timestamp.Format(timeFormat), r.URL.Path, err)
}
