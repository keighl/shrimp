package models

import (
  "testing"
  "time"
)

func Test_PasswordReset(t *testing.T) {
  setup(t)
  user := Uzer(t)
  pws := &PasswordReset{}
  pws.UserId = user.Id
  err := DB.Create(pws).Error
  expect(t, err, nil)
}

func Test_PasswordReset_RequiresUserId(t *testing.T) {
  setup(t)
  pws := &PasswordReset{}
  err := DB.Create(pws).Error
  refute(t, err, nil)
  expect(t, pws.ErrorMap["UserId"], true)
}

func Test_PasswordReset_SetsExpiresActive(t *testing.T) {
  setup(t)
  user := Uzer(t)
  pws := &PasswordReset{}
  pws.UserId = user.Id
  err := DB.Create(pws).Error
  expect(t, err, nil)
  expect(t, pws.Active, true)
  expect(t, pws.ExpiresAt.Format("RFC3339"), pws.CreatedAt.Add(6*time.Hour).Format("RFC3339"))
}

func Test_PasswordReset_SetsToken(t *testing.T) {
  setup(t)
  user := Uzer(t)
  pws := &PasswordReset{}
  pws.UserId = user.Id
  err := DB.Create(pws).Error
  expect(t, err, nil)
  refute(t, pws.Token, "")
}
