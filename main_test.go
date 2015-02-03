package main

import (
  "testing"
  u "github.com/keighl/shrimp/utils"
)

func Test_SetupSever(t *testing.T) {
  config := u.Config("test")
  server := u.MartiniServer(config)
  SetupServerRoutes(server)
}

func Test_SetupSubpackages(t *testing.T) {
  config := u.Config("test")
  DB, _ := u.RethinkSession(config)
  SetupSubpackages(DB, config)
}