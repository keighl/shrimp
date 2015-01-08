package api

import (
  u "shrimp/utils"
  "testing"
  "net/http/httptest"
  "github.com/go-martini/martini"
  "reflect"
)

func init() {
  Config = u.ConfigForFile("../config/test.json")
  DB = u.RethinkSession(Config)
}

func testTools(t *testing.T) (*martini.ClassicMartini, *httptest.ResponseRecorder) {
  return u.MartiniServer(false), httptest.NewRecorder()
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

//////////////////////////////
// API ENVELOPE //////////////

func Test_ServerErrorEnvelope(t *testing.T) {
  data := Data{}
  data.Error = &Error{"There was an unexpected error!", []string{}}
  targetByHand := data
  targetByMethod := ServerErrorEnvelope()
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}