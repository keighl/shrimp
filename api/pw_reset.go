package api

import (
  m "shrimp/models"
  "shrimp/utils"
  "github.com/martini-contrib/render"
  "github.com/keighl/mandrill"
  "github.com/go-martini/martini"
  "text/template"
  "bytes"
  "strings"
  "time"
  r "github.com/dancannon/gorethink"
)

/////////////////////

var savePasswordReset = func(reset *m.PasswordReset) error {
  return reset.Save()
}

func PasswordResetCreate(r render.Render, attrs m.PasswordResetAttrs) {

  user := userFromEmail(strings.TrimSpace(attrs.Email))
  if (user == nil) {
    r.JSON(400, ApiErrorEnvelope("That email isn't in our system!", []string{}))
    return
  }

  reset := &m.PasswordReset{}
  reset.UserId = user.Id
  err := savePasswordReset(reset)

  if (err != nil) {
    if (reset.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), reset.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  message, err := PasswordResetEmailMessage(user, reset)
  if (err != nil) {
    r.JSON(500, Api500Envelope())
    return
  }

  if !sendEmail(message) {
    r.JSON(500, Api500Envelope())
    return
  }

  data := &ApiData{PasswordReset: reset}
  r.JSON(201, data)
}

/////////////////////

var loadPasswordReset = func(token string) (*m.PasswordReset, error) {
  reset := &m.PasswordReset{}
  res, err := r.Table("password_resets").GetAllByIndex("token", token).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(reset)
  if (err != nil) { return nil, err }
  return reset, err
}

//////////////

func PasswordResetUpdate(params martini.Params, r render.Render, attrs m.UserAttrs) {

  reset, err := loadPasswordReset(params["token"])

  if (err != nil) {
    r.JSON(400, ApiErrorEnvelope("Invalid password reset token", nil))
    return
  }

  if (reset.ExpiresAt.Before(time.Now())) {
    r.JSON(400, ApiErrorEnvelope("The reset token has expired", nil))
    return
  }

  user, err := loadUser(reset.UserId)
  if (err != nil) {
    r.JSON(500, Api500Envelope())
    return
  }

  user.Password             = attrs.Password
  user.PasswordConfirmation = attrs.PasswordConfirmation
  err = saveUser(user)

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  reset.Active = false
  err = savePasswordReset(reset)
  if (err != nil) {
    // TODO notify us
  }

  r.JSON(200, ApiMessageEnvelope("Your password was reset"))
}

////////////////////////

type ResetPasswordEmailData struct {
  Config *utils.Configuration
  User *m.User
  PasswordReset *m.PasswordReset
}

func (x *ResetPasswordEmailData) ResetURL() string {
  return Config.BaseURL + "password-reset/" + x.PasswordReset.Token
}

func PasswordResetEmailMessage(user *m.User, reset *m.PasswordReset) (*mandrill.Message, error) {
  var textContent bytes.Buffer
  resetData := &ResetPasswordEmailData{Config, user, reset}

  t, err := template.ParseFiles("./emails/reset_password.text")
  if (err != nil) { return nil, err }

  err = t.Execute(&textContent, resetData)
  if (err != nil) { return nil, err }

  message := &mandrill.Message{}
  message.AddRecipient(user.Email, user.FullName(), "to")
  message.FromEmail = "reset@example.com"
  message.FromName  = Config.AppName
  message.Subject   = "Reset Your Password"
  message.Text      = textContent.String()

  return message, nil
}
