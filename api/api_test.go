package api

import (
  "shrimp/utils"
  "testing"
  "net/http/httptest"
  "github.com/go-martini/martini"
  "reflect"
)

var (
  alreadySetup bool
)

func setup(t *testing.T) {
  if (alreadySetup) { return }
  Config = utils.ConfigForFile("../config/test.json")
  DB = utils.RethinkSession(Config)
  alreadySetup = true
}

func testTools(t *testing.T) (*martini.ClassicMartini, *httptest.ResponseRecorder) {
  setup(t)
  return utils.MartiniServer(false), httptest.NewRecorder()
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

func Test_Api500Envelope(t *testing.T) {
  setup(t)
  data := ApiData{}
  data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
  targetByHand := data
  targetByMethod := Api500Envelope()
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}