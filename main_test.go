package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckSSLArgs(t *testing.T) {
	// arrange
	args := sslArgs{
		useEncryption: false,
		certFile:      "testdata/test_cert.perm",
		keyFile:       "testdata/test_key.perm",
	}

	// action
	err := checkSSLArgs(args)

	// verify
	assert.NoError(t, err)
}
