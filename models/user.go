package models

import (
  "time"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/dchest/uniuri"
  "bytes"
  "errors"
  "regexp"
  "strings"
)

type User struct {
  Errors []string `json:"errors,omitempty" sql:"-"`
  ErrorMap map[string]bool `json:"-" sql:"-"`
  Id int64 `json:"-"`
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
  ApiToken string `json:"-"`
}

func (x User) TableName() string {
  return "users"
}

type UserAttrs struct {
  NameFirst string `json:"name_first" form:"name_first"`
  NameLast string `json:"name_last" form:"name_last"`
  Email string `json:"email" form:"email"`
  IosPushToken string `json:"ios_push_token" form:"ios_push_token"`
  Password string `json:"password" form:"password"`
  PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation"`
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *User) BeforeSave() (err error) {
  x.Errors    = []string{}
  x.ErrorMap  = map[string]bool{}
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()

  x.Trimspace()
  x.ValidateName()
  x.ValidateEmail()
  x.ValidateEmailUniqueness()

  return
}

func (x *User) BeforeCreate() (err error) {

  x.ApiToken = uniuri.NewLen(30)

  if (x.Password != "") {
    if (x.validatePasswordAndConfirmation()) {
      x.SetPassword(x.Password)
    }
  } else {
    x.Errors = append(x.Errors, "Password can't be blank.")
    x.ErrorMap["Password"] = true
  }

  if (x.HasErrors()) {
    err = errors.New("There was a problem saving your info.")
  }

  return
}

func (x *User) BeforeUpdate() (err error) {

  if (x.Password != "") {
    if (x.validatePasswordAndConfirmation()) {
      x.SetPassword(x.Password)
    }
  }

  if (x.HasErrors()) {
    err = errors.New("There was a problem saving your info.")
  }

  return
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *User) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *User) validatePasswordAndConfirmation() (success bool) {
  // 8 non-whitespace characters or more
  regex := regexp.MustCompile("\\A\\S{8,}\\z")

  if (regex.MatchString(x.Password)) {
    if (x.PasswordConfirmation != "") {
      if (x.Password != x.PasswordConfirmation) {
        x.Errors = append(x.Errors, "Password confirmation doesn't match.")
        x.ErrorMap["PasswordConfirmation"] = true
      } else {
        return true
      }
    } else {
      x.Errors = append(x.Errors, "Password confirmation can't be blank.")
      x.ErrorMap["PasswordConfirmation"] = true
    }
  } else {
    x.Errors = append(x.Errors, "Password must be at least 8 characters, and have no spaces")
    x.ErrorMap["Password"] = true
  }

  return false
}

func (x *User) ValidateName() {
  if (x.NameFirst == "") {
    x.Errors = append(x.Errors, "First name can't be blank.")
    x.ErrorMap["NameFirst"] = true
  }

  if (x.NameLast == "") {
    x.Errors = append(x.Errors, "Last name can't be blank.")
    x.ErrorMap["NameLast"] = true
  }
}

func (x *User) ValidateEmail() {
  regex := regexp.MustCompile("\\A[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\\z")
  emailMatch := regex.MatchString(x.Email)

  if (!emailMatch) {
    x.Errors = append(x.Errors, "Your email address is invalid.")
    x.ErrorMap["Email"] = true
  }
}

func (x *User) ValidateEmailUniqueness() {
  count := 0
  if (x.Id == 0) {
    DB.
      Model(&User{}).
      Where("email = ?", strings.TrimSpace(x.Email)).
      Count(&count)
  } else {
    DB.
      Model(&User{}).
      Where("email = ?", strings.TrimSpace(x.Email)).
      Not([]int64{x.Id}).
      Count(&count)
  }

  if (count > 0) {
    x.Errors = append(x.Errors, "That email address is already taken.")
    x.ErrorMap["Email"] = true
  }
}

func (x *User) Trimspace() {
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

  saltedPassword.WriteString(password)
  saltedPassword.WriteString(x.Salt)

  err = bcrypt.CompareHashAndPassword([]byte(x.CryptedPassword), saltedPassword.Bytes())

  if err != nil { return false, err }

  return true, err
}

//////////////////////////////
// OTHER /////////////////////

func (x *UserAttrs) User() (*User) {
  return &User{
    NameFirst: x.NameFirst,
    NameLast: x.NameLast,
    Email: x.Email,
    Password: x.Password,
    PasswordConfirmation: x.PasswordConfirmation,
    IosPushToken: x.IosPushToken,
  }
}

func (x *User) UserAttrs() (*UserAttrs) {
  return &UserAttrs{
    NameFirst: x.NameFirst,
    NameLast: x.NameLast,
    Email: x.Email,
    Password: x.Password,
    PasswordConfirmation: x.PasswordConfirmation,
    IosPushToken: x.IosPushToken,
  }
}

func (x *User) FullName() string {
  return x.NameFirst + " " + x.NameLast
}