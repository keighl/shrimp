package main

import (
  "fmt"
  "testing"
  "net/http/httptest"
  "reflect"
  "github.com/modocache/gory"
  "github.com/keighl/mandrill"
  "github.com/go-martini/martini"
  "github.com/jrallison/go-workers"
)

var (
  alreadySetup bool
  err error
)

// TODO find a library that does setup/tear down instead of this
func setup(t *testing.T) {
  if (alreadySetup) { return }

  config = ConfigForFile("conf/test.json")
  db     = DBForConfig(config)

  workers.Configure(map[string]string{
    "server": config.WorkerServer,
    "database": config.WorkerDatabase,
    "pool": config.WorkerPool,
    "process": config.WorkerProcess,
  })

  db.Exec("TRUNCATE TABLE users")
  db.Exec("TRUNCATE TABLE api_sessions")
  db.Exec("TRUNCATE TABLE todos")
  db.Exec("TRUNCATE TABLE password_resets")
  DefineFactories()
  alreadySetup = true
}

func testTools(t *testing.T) (*martini.ClassicMartini, *httptest.ResponseRecorder) {
  setup(t)
  return martiniServer(false), httptest.NewRecorder()
}

func DefineFactories() {
  gory.Define("user", User{}, func(factory gory.Factory) {
    factory["NameFirst"] = gory.Sequence(
      func(n int) interface{} {
        return fmt.Sprintf("Kyle%d", n)
      })
    factory["NameLast"] = "Truscott"
    factory["Email"] = gory.Sequence(
      func(n int) interface{} {
        return fmt.Sprintf("kyle%d@example.com", n)
      })
    factory["Password"] = "Password1"
    factory["PasswordConfirmation"] = "Password1"
  })

  // TODO better way to handle this?
  gory.Define("userAttrs", UserAttrs{}, func(factory gory.Factory) {
    factory["NameFirst"] = gory.Sequence(
      func(n int) interface{} {
        return fmt.Sprintf("Kyle%d", n)
      })
    factory["NameLast"] = "Truscott"
    factory["Email"] = gory.Sequence(
      func(n int) interface{} {
        return fmt.Sprintf("kyle%d@attrs-example.com", n)
      })
    factory["Password"] = "Password1"
    factory["PasswordConfirmation"] = "Password1"
  })

  gory.Define("todo", Todo{}, func(factory gory.Factory) {
    factory["Title"] = gory.Sequence(
      func(n int) interface{} {
        return fmt.Sprintf("Todo #%d", n)
      })
    factory["Complete"] = false
  })

  gory.Define("todoAttrs", TodoAttrs{}, func(factory gory.Factory) {
    factory["Title"] = gory.Sequence(
      func(n int) interface{} {
        return fmt.Sprintf("Todo #%d", n)
      })
    factory["Complete"] = false
  })
}

func MockMailerTrue(c martini.Context) {
  c.Map(SendEmail(func (message *mandrill.Message)(bool) { return true }))
}

func MockMailerFalse(c martini.Context) {
  c.Map(SendEmail(func (message *mandrill.Message)(bool) { return false }))
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
  apiSezzion *ApiSession // memoized apiSession
)

func UserAndSession(t *testing.T) (*User, *ApiSession) {
  if (uzer != nil) { return uzer, apiSezzion }

  uzer = gory.Build("user").(*User)
  err = db.Create(uzer).Error
  if (err != nil) { t.Error(err) }

  apiSezzion = &ApiSession{ UserId: uzer.Id }
  err = db.Create(apiSezzion).Error
  if (err != nil) { t.Error(err) }

  return uzer, apiSezzion
}