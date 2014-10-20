package main

import (
  "fmt"
  "testing"
  "net/http/httptest"
  "reflect"
  "github.com/go-martini/martini"
  "github.com/jinzhu/gorm"
  "github.com/modocache/gory"
)

var (
  server *martini.ClassicMartini
  recorder *httptest.ResponseRecorder
  alreadySetup bool
  err error
)

// TODO find a library that does setup/tear down instead of this
func setup(t *testing.T) {

  recorder = httptest.NewRecorder()

  if (alreadySetup) { return }

  // TODO use Configuration file/object
  db, err = gorm.Open("mysql", "root:@tcp(localhost:3306)/shrimp_test?charset=utf8&parseTime=True")
  if err != nil { t.Error(err) }
  db.LogMode(false)
  db.Exec("TRUNCATE TABLE users")
  db.Exec("TRUNCATE TABLE api_sessions")
  db.Exec("TRUNCATE TABLE todos")

  DefineFactories()
  server = NewMartiniServer()
  ConfigureWorkerServer(true)

  alreadySetup = true
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

// TODO try our testify for these kinds of helpers
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