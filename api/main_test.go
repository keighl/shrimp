package main

import (
	"testing"
)

func Test_MartiniServer(t *testing.T) {
	server := MartiniServer()
	refute(t, server, nil)
}

func Test_SetupServerRoutes(t *testing.T) {
	server := MartiniServer()
	SetupServerRoutes(server)
}
