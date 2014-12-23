package models

import (
  "time"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/dchest/uniuri"
  "bytes"
  "regexp"
  "strings"
  r "github.com/dancannon/gorethink"
  "errors"
)

type User struct {
  Errors []string `gorethink:"-" json:"errors,omitempty"`
  ErrorMap map[string]bool `gorethink:"-" json:"-"`
  Id string `gorethink:"id,omitempty" json:"-"`
  CreatedAt time.Time `gorethink:"created_at" json:"created_at,omitempty"`
  UpdatedAt time.Time `gorethink:"updated_at" json:"updated_at,omitempty"`
  CryptedPassword string `gorethink:"crypted_password" json:"-"`
  Salt string `gorethink:"salt" json:"-"`
  Email string `gorethink:"email" json:"email,omitempty"`
  NameFirst string `gorethink:"name_first" json:"name_first,omitempty"`
  NameLast string `gorethink:"name_last" json:"name_last,omitempty"`
  Password string `gorethink:"-" json:"-" gorethink:"-"`
  PasswordConfirmation string `gorethink:"-" json:"-"`
  IosPushToken string `gorethink:"ios_push_token" json:"-"`
  ApiToken string `gorethink:"api_token" json:"-"`
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
// TRANSACTIONS //////////////

func (x *User) Save() error {

  if (!x.Validate()) {
    return errors.New("Validation errors")
  }

  if (x.Id == "") {
    x.BeforeCreate()
    res, err := r.Table("users").Insert(x).RunWrite(DB)
    if (err != nil) { return err }
    x.Id = res.GeneratedKeys[0]
  }

  x.BeforeUpdate()
  _, err := r.Table("users").Get(x.Id).Replace(x).RunWrite(DB)
  return err
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *User) BeforeCreate() {
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()
  x.ApiToken = uniuri.NewLen(30)
}

func (x *User) BeforeUpdate() {
  x.UpdatedAt = time.Now()
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *User) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.Trimspace()
  x.ValidateName()
  x.ValidateEmail()
  x.ValidateEmailUniqueness()

  // New record check
  if (x.Id == "") {
    x.ValidateRequiredPassword()
  } else {
    x.ValidateOptionalPassword()
  }

  return !x.HasErrors()
}

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
    x.Errors = append(x.Errors, "Email address is invalid.")
    x.ErrorMap["Email"] = true
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
    x.Errors = append(x.Errors, "That email address is already taken.")
    x.ErrorMap["Email"] = true
  }
}

func (x *User) ValidateRequiredPassword() {
  if (x.Password != "") {
    if (x.validatePasswordAndConfirmation()) {
      x.SetPassword(x.Password)
    }
  } else {
    x.Errors = append(x.Errors, "Password can't be blank.")
    x.ErrorMap["Password"] = true
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
  x.IosPushToken         = strings.TrimSpace(x.IosPushToken)
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