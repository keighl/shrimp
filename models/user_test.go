package models

import (
  "testing"
  "reflect"
  "github.com/dchest/uniuri"
)

func NewUser() *User {
  return &User{
    NameFirst: "cheese",
    NameLast: "cheese",
    Email: uniuri.NewLen(10) + "cheese@cheese.com",
    Password: "cheesedddd",
    PasswordConfirmation: "cheesedddd",
    IosPushToken: "cheese",
  }
}

//////////////////////////////
// TRANSACTIONS //////////////

func Test_User_Create_Success(t *testing.T) {
  setup(t)

  x := NewUser()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")
}

func Test_User_Create_Fail(t *testing.T) {
  setup(t)

  x := NewUser()
  x.NameFirst  = ""
  err := x.Save()
  refute(t, err, nil)
  expect(t, x.Id, "")
}

func Test_User_Update_Success(t *testing.T) {
  setup(t)

  x := NewUser()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Save()
  expect(t, err, nil)
}

func Test_User_Update_Fail(t *testing.T) {
  setup(t)

  x := NewUser()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  x.NameFirst = ""
  err = x.Save()
  refute(t, err, nil)
}

///////////

func Test_User_BeforeCreate(t *testing.T) {
  x := NewUser()
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
  refute(t, x.ApiToken, "")
}

func Test_User_BeforeUpdate(t *testing.T) {
  x := NewUser()
  x.BeforeUpdate()
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

func Test_User_SetCheckPassword(t *testing.T) {
  x := NewUser()
  x.SetPassword("CheesyBread3")
  res, _ := x.CheckPassword("CheesyBread")
  expect(t, res, false)
  res, _ = x.CheckPassword("CheesyBread3")
  expect(t, res, true)
}

func Test_User_Email_Uniqueness_NewUser(t *testing.T) {
  x := NewUser()
  _ = x.Save()

  y := NewUser()
  y.Email = x.Email
  err := y.Save()
  refute(t, err, nil)
  expect(t, y.ErrorMap["Email"], true)
}

func Test_User_Email_Uniqueness_ExistingUser(t *testing.T) {
  x := NewUser()
  err := x.Save()
  expect(t, err, nil)

  y := NewUser()
  err = y.Save()
  expect(t, err, nil)

  y.Email = x.Email
  err = y.Save()
  refute(t, err, nil)
  expect(t, y.ErrorMap["Email"], true)
}

func Test_User_Email_Format(t *testing.T) {
  x := NewUser()
  x.Email = "cheese"
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Email"], true)

  x.Email = uniuri.NewLen(30) + "cheese@cheese.com"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Email"], false)
}

func Test_User_Name_Presence(t *testing.T) {
  x := NewUser()
  x.NameFirst = ""
  x.NameLast = ""
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["NameFirst"], true)
  expect(t, x.ErrorMap["NameLast"], true)
  x.NameFirst = "Jerry"
  x.NameLast = "Seinfeld"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["NameFirst"], false)
  expect(t, x.ErrorMap["NameLast"], false)
}

func Test_User_Password_Format(t *testing.T) {
  x := NewUser()
  x.Password = "pass word"
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Password"], true)
  x.Password = "password"
  x.PasswordConfirmation = "password"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Password"], false)
}

func Test_User_Password_Confirmed(t *testing.T) {
  x := NewUser()
  x.Password = "password"

  // Blank
  x.PasswordConfirmation = ""
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["PasswordConfirmation"], true)

  // Wrong
  x.PasswordConfirmation = "password!!"
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["PasswordConfirmation"], true)

  // Correct
  x.Password = "password"
  x.PasswordConfirmation = "password"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["PasswordConfirmation"], false)
}

func Test_User_Create_Requires_Password(t *testing.T) {
  x := NewUser()
  x.Id = "" // signifies new record
  x.Password = ""
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Password"], true)
  x.Password = "password"
  x.PasswordConfirmation = "password"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Password"], false)
}

func Test_User_Update_Optional_Password(t *testing.T) {
  x := NewUser()
  x.Id = "XXXXX"
  x.Password = "password"
  x.PasswordConfirmation = ""
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["PasswordConfirmation"], true)
  x.Password = ""
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Password"], false)
}

func Test_User_FullName(t *testing.T) {
  x := NewUser()
  expect(t, x.FullName(), x.NameFirst + " " + x.NameLast)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_User_UserAttrs(t *testing.T) {
  obj := NewUser()
  targetByMethod := obj.UserAttrs()
  targetByHand := &UserAttrs{
    NameFirst: obj.NameFirst,
    NameLast: obj.NameLast,
    Email: obj.Email,
    Password: obj.Password,
    PasswordConfirmation: obj.PasswordConfirmation,
    IosPushToken: obj.IosPushToken,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}

func Test_UserAttrs_User(t *testing.T) {
  obj := &UserAttrs{
    NameFirst: "cheese",
    NameLast: "cheese",
    Email: "cheese",
    Password: "cheese",
    PasswordConfirmation: "cheese",
    IosPushToken: "cheese",
  }
  targetByMethod := obj.User()
  targetByHand := &User{
    NameFirst: obj.NameFirst,
    NameLast: obj.NameLast,
    Email: obj.Email,
    Password: obj.Password,
    PasswordConfirmation: obj.PasswordConfirmation,
    IosPushToken: obj.IosPushToken,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}

