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

func Test_PasswordReset_Table(t *testing.T) {
  x := NewPasswordReset()
  expect(t, x.Table(), "password_resets")
}

func Test_PasswordReset_RequiresUserId(t *testing.T) {
  x := &PasswordReset{}
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["UserId"], true)

  x.UserId = "XXXXXXXXXXXX"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["UserId"], false)
}

func Test_PasswordReset_BeforeCreate(t *testing.T) {
  x := &PasswordReset{UserId: "XXXXXXXXXXXX"}
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
  expect(t, x.ExpiresAt.Format("RFC3339"), x.CreatedAt.Add(6*time.Hour).Format("RFC3339"))
}

func Test_PasswordReset_BeforeUpdate(t *testing.T) {
  x := &PasswordReset{UserId: "XXXXXXXXXXXX"}
  x.BeforeUpdate()
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

