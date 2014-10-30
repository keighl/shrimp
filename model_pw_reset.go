package main

import (
  "time"
  "github.com/dchest/uniuri"
  "errors"
)

type PasswordReset struct {
  Errors []string `sql:"-" json:"-"`
  ErrorMap map[string]bool `sql:"-" json:"-"`
  Id int64 `json:"-"`
  CreatedAt time.Time `json:"-"`
  UpdatedAt time.Time `json:"-"`
  Token string `json:"-"`
  UserId int64 `json:"-"`
  Active bool `json:"-"`
  ExpiresAt time.Time `json:"expires_at"`
}

type PasswordResetAttrs struct {
  Email string `json:"email" form:"email"`
}

func (x PasswordReset) TableName() string {
  return "password_resets"
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *PasswordReset) BeforeCreate() (err error) {
  x.CreatedAt = time.Now()
  x.ExpiresAt = x.CreatedAt.Add(6*time.Hour)
  x.Token     = uniuri.NewLen(30)
  x.Active    = true
  return
}

func (x *PasswordReset) BeforeSave() (err error) {
  x.Errors    = []string{}
  x.ErrorMap  = map[string]bool{}
  x.UpdatedAt = time.Now()

  x.ValidateUserId()
  if (x.HasErrors()) {
    err = errors.New("There was a problem saving the password reset.")
  }
  return
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *PasswordReset) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *PasswordReset) ValidateUserId() {
  if (x.UserId == 0) {
    x.Errors = append(x.Errors, "UserId can't be blank.")
    x.ErrorMap["UserId"] = true
  }
}

