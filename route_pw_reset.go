package main

import (
  "github.com/martini-contrib/render"
  "github.com/keighl/mandrill"
  "github.com/go-martini/martini"
  "text/template"
  "bytes"
)


func RoutePasswordResetCreate(r render.Render, attrs PasswordResetAttrs, sendEmail SendEmail) {

  user := &User{}
  err := db.Where("email = ?", attrs.Email).First(user).Error

  if (err != nil) {
    r.JSON(400, ApiErrorEnvelope("That email isn't in our system", nil))
    return
  }

  pwr := &PasswordReset{}
  pwr.UserId = user.Id
  err = db.Create(pwr).Error

  if (err != nil) {
    r.JSON(500, Api500Envelope())
    return
  }

  message, err := PasswordResetEmailMessage(user, pwr)
  if (err != nil) {
    r.JSON(500, Api500Envelope())
    return
  }

  if !sendEmail(message) {
    r.JSON(500, Api500Envelope())
    return
  }

  data := &ApiData{PasswordReset: pwr}
  r.JSON(201, ApiEnvelope{data})
}

func RoutePasswordResetUpdate(params martini.Params, r render.Render, attrs UserAttrs) {
  user := &User{}
  pwr := &PasswordReset{}
  err := db.Where("token = ?", params["token"]).First(pwr).Error

  if (err != nil || !pwr.Active) {
    r.JSON(400, ApiErrorEnvelope("Invalid reset token", nil))
    return
  }

  err = db.Where("id = ?", pwr.UserId).First(user).Error
  if (err != nil) {
    r.JSON(400, ApiErrorEnvelope("Invalid reset token", nil))
    return
  }

  user.Password             = attrs.Password
  user.PasswordConfirmation = attrs.PasswordConfirmation
  err = db.Save(user).Error

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  r.JSON(200, ApiMessageEnvelope("Your password was reset"))
}

///////

type ResetPasswordEmailData struct {
  Config *Configuration
  User *User
  PasswordReset *PasswordReset
}

func (x *ResetPasswordEmailData) ResetURL() string {
  return config.BaseURL + "password-reset/" + x.PasswordReset.Token
}

func PasswordResetEmailMessage(user *User, pwr *PasswordReset) (*mandrill.Message, error) {
  var textContent bytes.Buffer
  resetData := &ResetPasswordEmailData{config, user, pwr}

  t, err := template.ParseFiles("emails/reset_password.text")
  if (err != nil) { return nil, err }

  err = t.Execute(&textContent, resetData)
  if (err != nil) { return nil, err }

  message := &mandrill.Message{}
  message.AddRecipient(user.Email, user.FullName(), "to")
  message.FromEmail = "reset@example.com"
  message.FromName  = config.AppName
  message.Subject   = "Reset Your Password"
  message.Text      = textContent.String()

  return message, nil
}
