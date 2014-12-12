package api

import (
  "shrimp/models"
  "shrimp/utils"
  "github.com/martini-contrib/render"
  "github.com/keighl/mandrill"
  "github.com/go-martini/martini"
  "text/template"
  "bytes"
)

func PasswordResetCreate(r render.Render, attrs models.PasswordResetAttrs, sendEmail SendEmail) {

  user := &models.User{}
  err := DB.Where("email = ?", attrs.Email).First(user).Error

  if (err != nil) {
    r.JSON(400, ApiErrorEnvelope("That email isn't in our system", nil))
    return
  }

  pwr := &models.PasswordReset{}
  pwr.UserId = user.Id
  err = DB.Create(pwr).Error

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
  r.JSON(201, data)
}

func PasswordResetUpdate(params martini.Params, r render.Render, attrs models.UserAttrs) {
  user := &models.User{}
  pwr := &models.PasswordReset{}
  err := DB.Where("token = ?", params["token"]).First(pwr).Error

  if (err != nil || !pwr.Active) {
    r.JSON(400, ApiErrorEnvelope("Invalid reset token", nil))
    return
  }

  err = DB.Where("id = ?", pwr.UserId).First(user).Error
  if (err != nil) {
    r.JSON(400, ApiErrorEnvelope("Invalid reset token", nil))
    return
  }

  user.Password             = attrs.Password
  user.PasswordConfirmation = attrs.PasswordConfirmation
  err = DB.Save(user).Error

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
  Config *utils.Configuration
  User *models.User
  PasswordReset *models.PasswordReset
}

func (x *ResetPasswordEmailData) ResetURL() string {
  return Config.BaseURL + "password-reset/" + x.PasswordReset.Token
}

func PasswordResetEmailMessage(user *models.User, pwr *models.PasswordReset) (*mandrill.Message, error) {
  var textContent bytes.Buffer
  resetData := &ResetPasswordEmailData{Config, user, pwr}

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
