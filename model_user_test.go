package main

import (
  "testing"
  "github.com/modocache/gory"
)

func Test_User_SetCheckPassword(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  user.SetPassword("CheesyBread3")
  res, _ := user.CheckPassword("CheesyBread")
  expect(t, res, false)
  res, _ = user.CheckPassword("CheesyBread3")
  expect(t, res, true)
}

func Test_User_Email_Uniqueness(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  err := db.Create(user).Error
  if (err != nil) { t.Error(err) }
  user2 := gory.Build("user").(*User)
  user2.Email = user.Email
  err = db.Create(user2).Error
  refute(t, err, nil)
  expect(t, user2.ErrorMap["Email"], true)
}

func Test_User_Email_Format(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  user.Email = "cheese"
  err := db.Create(user).Error
  refute(t, err, nil)
  expect(t, user.ErrorMap["Email"], true)
  user.Email = "cheese@cheese.com"
  err = db.Create(user).Error
  expect(t, err, nil)
}

func Test_User_Name_Presence(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  user.NameFirst = ""
  user.NameLast = ""
  err := db.Create(user).Error
  refute(t, err, nil)
  expect(t, user.ErrorMap["NameFirst"], true)
  expect(t, user.ErrorMap["NameLast"], true)
  user.NameFirst = "Jerry"
  user.NameLast = "Seinfeld"
  err = db.Create(user).Error
  expect(t, err, nil)
}

func Test_User_Password_Format(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  user.Password = "pass word"
  err := db.Create(user).Error
  refute(t, err, nil)
  expect(t, user.ErrorMap["Password"], true)
  user.Password = "password"
  user.PasswordConfirmation = "password"
  err = db.Create(user).Error
  expect(t, err, nil)
}

func Test_User_Password_Confirmed(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  user.Password = "password"
  user.PasswordConfirmation = "password!!"
  err := db.Create(user).Error
  refute(t, err, nil)
  expect(t, user.ErrorMap["PasswordConfirmation"], true)
  user.Password = "password"
  user.PasswordConfirmation = "password"
  err = db.Create(user).Error
  expect(t, err, nil)
}

func Test_User_Create_Requires_Password(t *testing.T) {
  setup(t)
  user := gory.Build("user").(*User)
  user.Password = ""
  err := db.Create(user).Error
  refute(t, err, nil)
  expect(t, user.ErrorMap["Password"], true)
  user.Password = "password"
  user.PasswordConfirmation = "password"
  err = db.Create(user).Error
  expect(t, err, nil)
}

func Test_User_Update_Optional_Password(t *testing.T) {
  setup(t)
  user, _ := UserAndSession(t)
  user.Password = ""
  user.NameLast = "Cheese"
  err := db.Save(user).Error
  expect(t, err, nil)
}
