package models

import (
  "time"
)

type PasswordReset struct {
  Record
  UserID string `gorethink:"user_id" json:"-"`
  Active bool `gorethink:"active" json:"-"`
  ExpiresAt time.Time `gorethink:"expires_at" json:"expires_at"`
}

type PasswordResetAttrs struct {
  Email string `json:"email" form:"email"`
}

func (x *PasswordReset) Table() string {
  return "password_resets"
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *PasswordReset) BeforeCreate() {
  x.Record.BeforeCreate()
  x.ExpiresAt = x.CreatedAt.Add(6*time.Hour)
  x.Active    = true
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *PasswordReset) Validate() {
  x.Record.Validate()
  x.ValidateUserID()
}

func (x *PasswordReset) ValidateUserID() {
  if (x.UserID == "") {
    x.ErrorOn("UserID", "UserID can't be blank.")
  }
}

