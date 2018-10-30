package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetLocalFileForRoute(t *testing.T) {
	// arrange
	unit := newFileServer()

	// action
	html, _ := unit.routes["/view-stream.html"]
	rootRedirected, _ := unit.routes["/"]
	// verify
	assert.Equal(t, "www/view-stream.html", html)
	assert.Equal(t, "www/view-stream.html", rootRedirected)

	// action
	js, _ := unit.routes["/jsmpeg.min.js"]
	// verify
	assert.Equal(t, "www/jsmpeg.min.js", js)

	// action
	_, ok := unit.routes["/.."]
	// verify
	assert.False(t, ok)
}
