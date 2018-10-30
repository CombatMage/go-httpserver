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

const dirWWW = "www"
const indexHTML = "www/index.html"
const timeFormat string = "Mon Jan 2 15:04:05 2006"

type fileServer struct {
	routes map[string]string
}

func (server fileServer) configureRoutes() {
	http.HandleFunc("/", log(server.serve))
}

func (server fileServer) listenAndServe(port string) {
	http.ListenAndServe(port, nil)
	fmt.Println("Stop listening")
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
func newFileServer() *fileServer {
	r := make(map[string]string)
	filepath.Walk(dirWWW, func(p string, info os.FileInfo, _ error) error {
		if p != dirWWW {
			p = filepath.ToSlash(p)
			route := strings.TrimPrefix(p, dirWWW)
			r[route] = p
		}
		return nil
	})
	r["/"] = indexHTML

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
