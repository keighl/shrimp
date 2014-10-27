package main

import (
  "github.com/martini-contrib/render"
  "github.com/keighl/mandrill"
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

  message := &mandrill.Message{}
  message.AddRecipient(user.Email, user.FullName(), "to")
  message.FromEmail = "reset@example.com"
  message.FromName  = "AppName"
  message.Subject   = "You won the prize!"
  message.HTML      = "<h1>You won!!</h1>"
  message.Text      = "You won!!"

  if !sendEmail(message) {
    r.JSON(500, Api500Envelope())
    return
  }

  data := &ApiData{CurrentUser: user, PasswordReset: pwr}
  r.JSON(201, ApiEnvelope{data})
}