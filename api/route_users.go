package api

import (
  "shrimp/models"
  "github.com/martini-contrib/render"
)

func UserCreate(r render.Render, attrs models.UserAttrs) {

  var err error
  user := attrs.User()
  err = DB.Create(user).Error

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  apiSession := models.ApiSession{UserId: user.Id}
  err = DB.Create(&apiSession).Error

  if (err != nil) {
    r.JSON(500, Api500Envelope())
    return
  }

  data := &ApiData{User: user, ApiSession: &apiSession, CurrentUser: user}
  r.JSON(201, ApiEnvelope{data})
}

func UserMe(r render.Render, user *models.User) {
  data := &ApiData{User: user, CurrentUser: user}
  r.JSON(200, ApiEnvelope{data})
}

func UserUpdate(r render.Render, user *models.User, attrs models.UserAttrs) {

  var err error
  err = DB.Model(user).Updates(attrs).Error

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  data := &ApiData{User: user, CurrentUser: user}
  r.JSON(200, ApiEnvelope{data})
}

