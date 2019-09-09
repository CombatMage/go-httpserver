package main

import "testing"

func TestGetLocalFileForRoute(t *testing.T) {
	// arrange
	unit := newFileServer("testdata/www", "testdata/www/index.html")

	// action
	html, _ := unit.routes["/index.html"]
	rootRedirected, _ := unit.routes["/"]

	// verify
	equals(t, "testdata/www/index.html", html)
	equals(t, "testdata/www/index.html", rootRedirected)
}

func TestGetLocalFileForRoute_shouldReturnError(t *testing.T) {
	// arrange
	unit := newFileServer("testdata/www", "testdata/www/index.html")

	// action
	_, ok := unit.routes["/.."]
	
	// verify
	assert(t, ok == false, "Should return not ok")
}
