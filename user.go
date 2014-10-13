package main

import (
  "time"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/dchest/uniuri"
  "bytes"
  "errors"
  "regexp"
  "strings"
  "database/sql"
)

type User struct {
  Errors []string `json:"errors,omitempty" sql:"-"`
  Id sql.NullInt64 `json:"-"`
  CreatedAt time.Time `json:"created_at,omitempty"`
  UpdatedAt time.Time `json:"updated_at,omitempty"`
  CryptedPassword string `json:"-"`
  Salt string `json:"-"`
  Email string `json:"email,omitempty"`
  NameFirst string `json:"name_first,omitempty"`
  NameLast string `json:"name_last,omitempty"`
  Password string `json:"-" sql:"-"`
  PasswordConfirmation string `json:"-" sql:"-"`
  IosPushToken string `json:"-"`
}

type UserAttrs struct {
  NameFirst string `form:"user[name_first]"`
  NameLast string `form:"user[name_last]"`
  Email string `form:"user[email]"`
  IosPushToken string `form:"user[ios_push_token]"`
  Password string `form:"user[password]"`
  PasswordConfirmation string `form:"user[password_confirmation]"`
}

type UserLoginAttrs struct {
  Email string `form:"user[email]"`
  Password string `form:"user[password]"`
}

func (x User) TableName() string {
  return "users"
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *User) BeforeCreate() (err error) {

  x.Errors    = []string{}
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()

  x.TrimSpace()
  x.validateName()
  x.validateEmail()


  if (x.Password != "") {
    if (userPasswordValid(x.Password)) {
      if (x.PasswordConfirmation != "") {
        if (x.Password != x.PasswordConfirmation) {
          x.Errors = append(x.Errors, "Password confirmation doesn't match.")
        } else {
          x.SetPassword(x.Password)
        }
      } else {
        x.Errors = append(x.Errors, "Password confirmation can't be blank.")
      }
    } else {
      x.Errors = append(x.Errors, "Password must be at least 8 characters, and have no spaces")
    }
  } else {
    x.Errors = append(x.Errors, "Password can't be blank.")
  }

  if (len(x.Errors) > 0) {
    err = errors.New("There was a problem saving your info.")
  }

  return
}

func (x *User) BeforeUpdate() (err error) {

  x.Errors    = []string{}
  x.UpdatedAt = time.Now()

  x.TrimSpace()
  x.validateName()
  x.validateEmail()

  if (x.Password != "") {
    if (userPasswordValid(x.Password)) {
      if (x.PasswordConfirmation != "") {
        if (x.Password != x.PasswordConfirmation) {
          x.Errors = append(x.Errors, "Password confirmation doesn't match.")
        } else {
          x.SetPassword(x.Password)
        }
      } else {
        x.Errors = append(x.Errors, "Password confirmation can't be blank.")
      }
    } else {
      x.Errors = append(x.Errors, "Password must be at least 8 characters, and have no spaces")
    }
  }

  if (len(x.Errors) > 0) {
    err = errors.New("There was a problem saving your info.")
  }

  return
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *User) hasErrors() (bool) {
  return len(x.Errors) > 0
}

// 8 non-whitespace characters or more
func userPasswordValid(password string) (bool) {
  regex := regexp.MustCompile("\\A\\S{8,}\\z")
  return regex.MatchString(password)
}

func (x *User) validateName() {
  if (x.NameFirst == "") {
    x.Errors = append(x.Errors, "First name can't be blank.")
  }

  if (x.NameLast == "") {
    x.Errors = append(x.Errors, "Last name can't be blank.")
  }
}

func (x *User) validateEmail() {
  regex := regexp.MustCompile("\\A[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\\z")
  emailMatch := regex.MatchString(x.Email)

  if (!emailMatch) {
    x.Errors = append(x.Errors, "Your email address is invalid.")
  }
}

func (x *User) TrimSpace() {
  x.NameFirst            = strings.TrimSpace(x.NameFirst)
  x.NameLast             = strings.TrimSpace(x.NameLast)
  x.Email                = strings.TrimSpace(x.Email)
  x.PasswordConfirmation = strings.TrimSpace(x.PasswordConfirmation)
  x.Password             = strings.TrimSpace(x.Password)
  x.IosPushToken         = strings.TrimSpace(x.IosPushToken)
}

//////////////////////////////
// PASSWORD UTILS ////////////

func (x *User) SetPassword(password string) (err error) {

  var saltedPassword bytes.Buffer
  var cryptedPassword []byte

  x.Salt = uniuri.NewLen(15)

  saltedPassword.WriteString(password)
  saltedPassword.WriteString(x.Salt)

  cryptedPassword, err = bcrypt.GenerateFromPassword(saltedPassword.Bytes(), 10)

  if err != nil { return err }

  x.CryptedPassword = string(cryptedPassword)

  return
}

func (x *User) CheckPassword(password string) (success bool, err error) {

  var saltedPassword bytes.Buffer
  var cryptedPassword []byte

  saltedPassword.WriteString(password)
  saltedPassword.WriteString(x.Salt)

  cryptedPassword, err = bcrypt.GenerateFromPassword(saltedPassword.Bytes(), 10)

  if err != nil { return false, err }

  err = bcrypt.CompareHashAndPassword(cryptedPassword, saltedPassword.Bytes())

  if err != nil { return false, err }

  return true, err
}
