package models

import (
  "time"
  "github.com/dchest/uniuri"
  r "github.com/dancannon/gorethink"
  "errors"
)

type PasswordReset struct {
  Errors []string `gorethink:"-" json:"-"`
  ErrorMap map[string]bool `gorethink:"-" json:"-"`
  Id string `gorethink:"id,omitempty" json:"-"`
  CreatedAt time.Time `gorethink:"created_at" json:"-"`
  UpdatedAt time.Time `gorethink:"updated_at" json:"-"`
  Token string `gorethink:"token" json:"-"`
  UserId string `gorethink:"user_id" json:"-"`
  Active bool `gorethink:"active" json:"-"`
  ExpiresAt time.Time `gorethink:"expires_at" json:"expires_at"`
}

type PasswordResetAttrs struct {
  Email string `json:"email" form:"email"`
}

//////////////////////////////
// TRANSACTIONS //////////////

func (x *PasswordReset) Save() error {

  if (!x.Validate()) {
    return errors.New("Validation errors")
  }

  if (x.Id == "") {
    x.BeforeCreate()
    res, err := r.Table("password_resets").Insert(x).RunWrite(DB)
    if (err != nil) { return err }
    x.Id = res.GeneratedKeys[0]
  }

  x.BeforeUpdate()
  _, err := r.Table("password_resets").Get(x.Id).Replace(x).RunWrite(DB)
  return err
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *PasswordReset) BeforeCreate() {
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()
  x.ExpiresAt = x.CreatedAt.Add(6*time.Hour)
  x.Token     = uniuri.NewLen(30)
  x.Active    = true
}

func (x *PasswordReset) BeforeUpdate() {
  x.UpdatedAt = time.Now()
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *PasswordReset) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.ValidateUserId()
  return !x.HasErrors()
}

func (x *PasswordReset) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *PasswordReset) ValidateUserId() {
  if (x.UserId == "") {
    x.Errors = append(x.Errors, "UserId can't be blank.")
    x.ErrorMap["UserId"] = true
  }
}

