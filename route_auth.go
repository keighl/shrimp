package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "net/http"
  "strings"
)

/////////////////////////////

func RouteAuthorize(c martini.Context, r render.Render, req *http.Request) {
  var err error
  session_token := req.URL.Query().Get("session_token")
  user := &User{}

  err = db.
    Table("users").
    Select("users.*").
    Joins("INNER JOIN api_sessions x on x.user_id = users.id").
    Where("session_token = ?", strings.TrimSpace(session_token)).
    Limit(1).
    Scan(user).Error

  if (err != nil) {
    r.JSON(401, error400Envelope("Your token is invalid!", []string{}))
    return
  }

  c.Map(user) // Map the user to be used in the route
}

/////////////////////////////

func RouteLogin(r render.Render, attrs UserLoginAttrs) {

  var err error
  var success bool
  user := &User{}

  err = db.Where("email = ?", strings.TrimSpace(attrs.Email)).First(user).Error

  if (err != nil) {
    r.JSON(401, error400Envelope("Your email or password is invalid!", []string{}))
    return
  }

  success, err = user.CheckPassword(strings.TrimSpace(attrs.Password))

  if (err != nil || !success) {
    r.JSON(401, error400Envelope("Your email or password is invalid!", []string{}))
    return
  }

  apiSession := &ApiSession{ UserId: user.Id }
  err = db.Create(apiSession).Error

  if (err != nil) {
    r.JSON(500, error500Envelope())
    return
  }

  data := &ApiData{ApiSession: apiSession}
  r.JSON(200, ApiEnvelope{data})
}
