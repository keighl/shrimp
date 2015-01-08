package models

import (
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/dchest/uniuri"
  "bytes"
  "regexp"
  "strings"
  r "github.com/dancannon/gorethink"
)

type User struct {
  Record
  CryptedPassword string `gorethink:"crypted_password" json:"-"`
  Salt string `gorethink:"salt" json:"-"`
  Email string `gorethink:"email" json:"email,omitempty"`
  NameFirst string `gorethink:"name_first" json:"name_first,omitempty"`
  NameLast string `gorethink:"name_last" json:"name_last,omitempty"`
  Password string `gorethink:"-" json:"-" gorethink:"-"`
  PasswordConfirmation string `gorethink:"-" json:"-"`
  APIToken string `gorethink:"api_token" json:"-"`
}

type UserAttrs struct {
  NameFirst string `json:"name_first" form:"name_first"`
  NameLast string `json:"name_last" form:"name_last"`
  Email string `json:"email" form:"email"`
  Password string `json:"password" form:"password"`
  PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation"`
}

func (x *User) Table() string {
  return "users"
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *User) BeforeCreate() {
  x.Record.BeforeCreate()
  x.APIToken = uniuri.NewLen(30)
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *User) Validate() {
  x.Record.Validate()
  x.Trimspace()
  x.ValidateName()
  x.ValidateEmail()
  x.ValidateEmailUniqueness()

  if (x.IsNewRecord()) {
    x.ValidateRequiredPassword()
  } else {
    x.ValidateOptionalPassword()
  }
}

func (x *User) validatePasswordAndConfirmation() (success bool) {
  // 8 non-whitespace characters or more
  regex := regexp.MustCompile("\\A\\S{8,}\\z")

  if (regex.MatchString(x.Password)) {
    if (x.PasswordConfirmation != "") {
      if (x.Password != x.PasswordConfirmation) {
        x.ErrorOn("PasswordConfirmation", "Password confirmation doesn't match.")
      } else {
        return true
      }
    } else {
      x.ErrorOn("PasswordConfirmation", "Password confirmation can't be blank.")
    }
  } else {
    x.ErrorOn("Password", "Password must be at least 8 characters, and have no spaces")
  }

  return false
}

func (x *User) ValidateName() {
  if (x.NameFirst == "") {
    x.ErrorOn("NameFirst", "First name can't be blank.")
  }

  if (x.NameLast == "") {
    x.ErrorOn("NameLast", "Last name can't be blank.")
  }
}

func (x *User) ValidateEmail() {
  regex := regexp.MustCompile("\\A[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\\z")
  emailMatch := regex.MatchString(x.Email)

  if (!emailMatch) {
    x.ErrorOn("Email", "Email address is invalid.")
  }
}

func (x *User) ValidateEmailUniqueness() {

  var count int

  filter := func(row r.Term) r.Term {
    return row.Field("email").Eq(x.Email).And(row.Field("id").Ne(x.Id))
  }

  res, _ := r.Table("users").
    Filter(filter).
    Count().
    Run(DB)

  res.One(&count)

  if (count > 0) {
    x.ErrorOn("Email", "That email address is already taken.")
  }
}

func (x *User) ValidateRequiredPassword() {
  if (x.Password != "") {
    if (x.validatePasswordAndConfirmation()) {
      x.SetPassword(x.Password)
    }
  } else {
    x.ErrorOn("Password", "Password can't be blank.")
  }
}

func (x *User) ValidateOptionalPassword() {
  if (x.Password != "") {
    if (x.validatePasswordAndConfirmation()) {
      x.SetPassword(x.Password)
    }
  }
}

func (x *User) Trimspace() {
  x.NameFirst            = strings.TrimSpace(x.NameFirst)
  x.NameLast             = strings.TrimSpace(x.NameLast)
  x.Email                = strings.TrimSpace(x.Email)
  x.PasswordConfirmation = strings.TrimSpace(x.PasswordConfirmation)
  x.Password             = strings.TrimSpace(x.Password)
}

//////////////////////////////
// PASSWORD UTILS ////////////

func (x *User) SetPassword(password string) {

  var saltedPassword bytes.Buffer
  var cryptedPassword []byte

  x.Salt = uniuri.NewLen(15)

  saltedPassword.WriteString(password)
  saltedPassword.WriteString(x.Salt)

  cryptedPassword, _ = bcrypt.GenerateFromPassword(saltedPassword.Bytes(), 10)

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

func (x *User) UpdateFromAttrs(attrs UserAttrs) {
  if (attrs.NameFirst != "") { x.NameFirst = attrs.NameFirst }
  if (attrs.NameLast != "") { x.NameLast = attrs.NameLast }
  if (attrs.Email != "") { x.Email = attrs.Email }
  if (attrs.Password != "") { x.Password = attrs.Password }
  if (attrs.PasswordConfirmation != "") { x.PasswordConfirmation = attrs.PasswordConfirmation }
}

func (x *UserAttrs) User() (*User) {
  return &User{
    NameFirst: x.NameFirst,
    NameLast: x.NameLast,
    Email: x.Email,
    Password: x.Password,
    PasswordConfirmation: x.PasswordConfirmation,
  }
}

func (x *User) FullName() string {
  return x.NameFirst + " " + x.NameLast
}