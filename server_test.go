package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetLocalFileForRoute(t *testing.T) {
	// arrange
	unit := newFileServer("testdata/www", "testdata/www/index.html")

	// action
	html, _ := unit.routes["/index.html"]
	rootRedirected, _ := unit.routes["/"]
	// verify
	assert.Equal(t, "testdata/www/index.html", html)
	assert.Equal(t, "testdata/www/index.html", rootRedirected)

	// action
	js, _ := unit.routes["/gopher.png"]
	// verify
	assert.Equal(t, "testdata/www/gopher.png", js)

	// action
	_, ok := unit.routes["/.."]
	// verify
	assert.False(t, ok)
}
