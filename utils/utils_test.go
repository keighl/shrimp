package utils

import (
  "testing"
  "reflect"
)

func expect(t *testing.T, a interface{}, b interface{}) {
  if a != b {
    t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func refute(t *testing.T, a interface{}, b interface{}) {
  if a == b {
    t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func Test_MartiniServer(t *testing.T) {
  config := Config("test")
  server := MartiniServer(config)
  refute(t, server, nil)
}

func Test_MartiniServer_Logger(t *testing.T) {
  config := Config("test")
  config.ServerLoggingEnabled = true
  server := MartiniServer(config)
  refute(t, server, nil)
}

func Test_RethinkSession(t *testing.T) {
  config := Config("test")
  session, err := RethinkSession(config)
  expect(t, err, nil)
  refute(t, session, nil)
}

func Test_RethinkSession_Error(t *testing.T) {
  config := Config("test")
  config.RethinkHost = "cheese:28015"
  _, err := RethinkSession(config)
  refute(t, err, nil)
}