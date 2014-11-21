package models

import (
  "shrimp/utils"
  "testing"
  "net/http/httptest"
  "reflect"
  "github.com/go-martini/martini"
  _ "github.com/jrallison/go-workers"
  "github.com/modocache/gory"
)

var (
  alreadySetup bool
  err error
)

// TODO find a library that does setup/tear down instead of this
func setup(t *testing.T) {
  if (alreadySetup) { return }
  Config = utils.ConfigForFile("../conf/test.json")
  DB = utils.DBForConfig(Config)

  // workers.Configure(map[string]string{
  //   "server": config.WorkerServer,
  //   "database": config.WorkerDatabase,
  //   "pool": config.WorkerPool,
  //   "process": config.WorkerProcess,
  // })

  DB.Exec("TRUNCATE TABLE users")
  DB.Exec("TRUNCATE TABLE todos")
  DB.Exec("TRUNCATE TABLE password_resets")
  DefineFactories()
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

var (
  uzer *User // memoized user
)

func Uzer(t *testing.T) (*User) {
  if (uzer != nil) { return uzer }

  uzer = gory.Build("user").(*User)
  err = DB.Create(uzer).Error
  if (err != nil) { t.Error(err) }

  return uzer
}