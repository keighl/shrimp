package models

import (
  "testing"
  "reflect"
  "github.com/dchest/uniuri"
)

func NewTestUser() *User {
  return &User{
    NameFirst: "cheese",
    NameLast: "cheese",
    Email: uniuri.NewLen(10) + "cheese@cheese.com",
    Password: "cheesedddd",
    PasswordConfirmation: "cheesedddd",
  }
}

func NewTestUserPersisted() *User {
  user := NewTestUser()
  user.SetId(uniuri.NewLen(10))
  return user
}

func Test_User_Table(t *testing.T) {
  user := NewTestUser()
  expect(t, user.Table(), "users")
}

func Test_User_SetCheckPassword(t *testing.T) {
  x := NewTestUser()
  x.SetPassword("CheesyBread3")
  res, _ := x.CheckPassword("CheesyBread")
  expect(t, res, false)
  res, _ = x.CheckPassword("CheesyBread3")
  expect(t, res, true)
}

func Test_User_Email_Uniqueness_NewTestUser(t *testing.T) {
  x := NewTestUser()
  err := Save(x)
  expect(t, err, nil)

  y := NewTestUser()
  y.Email = x.Email
  err = Save(y)
  refute(t, err, nil)
  expect(t, y.ErrorMap["Email"], true)
}

func Test_User_Email_Uniqueness_ExistingUser(t *testing.T) {
  x := NewTestUser()
  err := Save(x)
  expect(t, err, nil)

  y := NewTestUser()
  err = Save(y)
  expect(t, err, nil)

  y.Email = x.Email
  err = Save(y)
  refute(t, err, nil)
  expect(t, y.ErrorMap["Email"], true)
}

func Test_User_Email_Format(t *testing.T) {
  x := NewTestUser()
  x.Email = "cheese"
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Email"], true)

  x.Email = uniuri.NewLen(30) + "cheese@cheese.com"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Email"], false)
}

func Test_User_Name_Presence(t *testing.T) {
  x := NewTestUser()
  x.NameFirst = ""
  x.NameLast = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["NameFirst"], true)
  expect(t, x.ErrorMap["NameLast"], true)

  x.NameFirst = "Jerry"
  x.NameLast = "Seinfeld"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["NameFirst"], false)
  expect(t, x.ErrorMap["NameLast"], false)
}

func Test_User_Password_Format(t *testing.T) {
  x := NewTestUser()
  x.Password = "pas sw ord"
  x.PasswordConfirmation = "pas sw ord"
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Password"], true)

  x.Password = "password"
  x.PasswordConfirmation = "password"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Password"], false)
}

func Test_User_Password_Confirmed(t *testing.T) {
  x := NewTestUser()
  x.Password = "password"

  // Blank
  x.PasswordConfirmation = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["PasswordConfirmation"], true)

  // Wrong
  x.PasswordConfirmation = "password!!"
  err = Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["PasswordConfirmation"], true)

  // Correct
  x.Password = "password"
  x.PasswordConfirmation = "password"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["PasswordConfirmation"], false)
}

func Test_User_Create_Requires_Password(t *testing.T) {
  x := NewTestUser()
  x.Password = ""
  err := Save(x)

  refute(t, err, nil)
  expect(t, x.ErrorMap["Password"], true)

  x.Password = "password"
  x.PasswordConfirmation = "password"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Password"], false)
}

func Test_User_Update_Optional_Password(t *testing.T) {
  x := NewTestUser()
  x.SetId("XXXXX")
  x.Password = "password"
  x.PasswordConfirmation = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["PasswordConfirmation"], true)

  x.Password = ""
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Password"], false)
}

func Test_User_FullName(t *testing.T) {
  x := NewTestUser()
  expect(t, x.FullName(), x.NameFirst + " " + x.NameLast)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_User_UpdateAttrs(t *testing.T) {
  obj := &User{
    NameFirst: "cheese",
    NameLast: "cheese",
    Email: "cheese",
    Password: "cheese",
    PasswordConfirmation: "cheese",
  }
  attrs := UserAttrs{
    NameFirst: "cheesex",
    NameLast: "cheesex",
    Email: "cheesex",
    Password: "cheesex",
    PasswordConfirmation: "cheesex",
  }
  obj.UpdateFromAttrs(attrs)
  targetByHand := &User{
    NameFirst: attrs.NameFirst,
    NameLast: attrs.NameLast,
    Email: attrs.Email,
    Password: attrs.Password,
    PasswordConfirmation: attrs.PasswordConfirmation,
  }

  expect(t, reflect.DeepEqual(targetByHand, obj), true)
}

func Test_UserAttrs_User(t *testing.T) {
  obj := &UserAttrs{
    NameFirst: "cheese",
    NameLast: "cheese",
    Email: "cheese",
    Password: "cheese",
    PasswordConfirmation: "cheese",
  }
  targetByMethod := obj.User()
  targetByHand := &User{
    NameFirst: obj.NameFirst,
    NameLast: obj.NameLast,
    Email: obj.Email,
    Password: obj.Password,
    PasswordConfirmation: obj.PasswordConfirmation,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}
