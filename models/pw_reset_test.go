package models

import (
  "testing"
  "time"
)

func NewPasswordReset() *PasswordReset{
  return &PasswordReset{
    UserId: "XXXXXXXXXXXX",
  }
}

func Test_PasswordReset_RequiresUserId(t *testing.T) {
  x := &PasswordReset{}
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["UserId"], true)
}

func Test_PasswordReset_BeforeCreate(t *testing.T) {
  x := &PasswordReset{UserId: "XXXXXXXXXXXX"}
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
  expect(t, x.ExpiresAt.Format("RFC3339"), x.CreatedAt.Add(6*time.Hour).Format("RFC3339"))
  refute(t, x.Token, "")
}

func Test_PasswordReset_BeforeUpdate(t *testing.T) {
  x := &PasswordReset{UserId: "XXXXXXXXXXXX"}
  x.BeforeUpdate()
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

func Test_PasswordReset_Create_Success(t *testing.T) {
  setup(t)

  x := NewPasswordReset()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")
}

func Test_PasswordReset_Create_Fail(t *testing.T) {
  setup(t)

  x := NewPasswordReset()
  x.UserId  = ""
  err := x.Save()
  refute(t, err, nil)
  expect(t, x.Id, "")
}

func Test_PasswordReset_Update_Success(t *testing.T) {
  setup(t)

  x := NewPasswordReset()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Save()
  expect(t, err, nil)
}

func Test_PasswordReset_Update_Fail(t *testing.T) {
  setup(t)

  x := NewPasswordReset()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  x.UserId = ""
  err = x.Save()
  refute(t, err, nil)
}
