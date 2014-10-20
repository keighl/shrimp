package main

import (
  "github.com/martini-contrib/render"
)

func RouteUserCreate(r render.Render, attrs UserAttrs) {

  var err error
  user := attrs.User()
  err = db.Create(user).Error

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, Error400Envelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, Error500Envelope())
    }
    return
  }

  apiSession := ApiSession{UserId: user.Id}
  err = db.Create(&apiSession).Error

  if (err != nil) {
    r.JSON(500, Error500Envelope())
    return
  }

  data := &ApiData{User: user, ApiSession: &apiSession, CurrentUser: user}
  r.JSON(201, ApiEnvelope{data})
}

/////////////////

func RouteUserUpdate(r render.Render, user *User, attrs UserAttrs) {

  var err error
  err = db.Model(user).Updates(attrs).Error

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, Error400Envelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, Error500Envelope())
    }
    return
  }

  data := &ApiData{User: user, CurrentUser: user}
  r.JSON(200, ApiEnvelope{data})
}

