package models

import (
	"testing"
	"time"
)

func NewPasswordReset() *PasswordReset {
	return &PasswordReset{
		UserID: "XXXXXXXXXXXX",
	}
}

func Test_PasswordReset_Table(t *testing.T) {
	x := NewPasswordReset()
	expect(t, x.Table(), "password_resets")
}

func Test_PasswordReset_RequiresUserID(t *testing.T) {
	x := &PasswordReset{}
	err := Save(x, true)
	refute(t, err, nil)
	expect(t, x.ErrorMap["UserID"], true)

	x.UserID = "XXXXXXXXXXXX"
	err = Save(x, true)
	expect(t, err, nil)
	expect(t, x.ErrorMap["UserID"], false)
}

func Test_PasswordReset_BeforeCreate(t *testing.T) {
	x := &PasswordReset{UserID: "XXXXXXXXXXXX"}
	x.BeforeCreate()
	refute(t, x.CreatedAt.Format("RFC3339"), nil)
	expect(t, x.ExpiresAt.Format("RFC3339"), x.CreatedAt.Add(6*time.Hour).Format("RFC3339"))
}

func Test_PasswordReset_BeforeUpdate(t *testing.T) {
	x := &PasswordReset{UserID: "XXXXXXXXXXXX"}
	x.BeforeUpdate()
	refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}
