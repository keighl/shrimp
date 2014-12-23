package models

import (
  "testing"
  "reflect"
  u "shrimp/utils"
)

var (
  alreadySetup bool
)

func setup(t *testing.T) {
  if (alreadySetup) { return }
  Config = u.ConfigForFile("../config/test.json")
  DB = u.RethinkSession(Config)
  alreadySetup = true
}

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
