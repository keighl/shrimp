package main

import (
  "github.com/martini-contrib/render"
)

func RouteUserCreate(r render.Render, attrs UserAttrs) {

  var err error

  user := User {
    NameFirst: attrs.NameFirst,
    NameLast: attrs.NameLast,
    Email: attrs.Email,
    Password: attrs.Password,
    PasswordConfirmation: attrs.PasswordConfirmation,
    IosPushToken: attrs.IosPushToken,
  }

  err = db.Create(&user).Error

  if (err != nil) {
    if (user.hasErrors()) {
      r.JSON(400, error400Envelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, error500Envelope())
    }
    return
  }

  apiSession := ApiSession{ UserId: user.Id }
  err = db.Create(&apiSession).Error

  if (err != nil) {
    r.JSON(500, error500Envelope())
    return
  }

  data := &ApiData{User: &user, ApiSession: &apiSession}
  r.JSON(201, ApiEnvelope{data})
}

/////////////////

func RouteUserUpdate(r render.Render, user *User, attrs UserAttrs) {

  var err error
  err = db.Model(user).Updates(attrs).Error

  if (err != nil) {
    if (user.hasErrors()) {
      r.JSON(400, error400Envelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, error500Envelope())
    }
    return
  }

  data := &ApiData{User: user}
  r.JSON(200, ApiEnvelope{data})
}

