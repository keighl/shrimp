package models

import (
  "shrimp/utils"
  "github.com/jinzhu/gorm"
  "fmt"
  "github.com/modocache/gory"
)

var  (
  Config *utils.Configuration
  DB gorm.DB
)

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