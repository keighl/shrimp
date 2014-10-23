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
  var sessionToken string

  sessionToken = req.Header.Get("X-SESSION-TOKEN")
  if (sessionToken == "") {
    sessionToken = req.URL.Query().Get("session_token")
  }

  user := &User{}

  err = db.
    Table("users").
    Select("users.*").
    Joins("INNER JOIN api_sessions x on x.user_id = users.id").
    Where("session_token = ?", strings.TrimSpace(sessionToken)).
    Limit(1).
    Scan(user).Error

  if (err != nil) {
    r.JSON(401, ApiErrorEnvelope("Your token is invalid!", []string{}))
    return
  }

  c.Map(user) // Map the user to be used in the route
}

/////////////////////////////

func RouteLogin(r render.Render, attrs UserAttrs) {

  var err error
  var success bool
  user := &User{}

  err = db.Where("email = ?", strings.TrimSpace(attrs.Email)).First(user).Error

  if (err != nil) {
    r.JSON(401, ApiErrorEnvelope("Your email or password is invalid!", []string{}))
    return
  }

  success, err = user.CheckPassword(strings.TrimSpace(attrs.Password))

  if (err != nil || !success) {
    r.JSON(401, ApiErrorEnvelope("Your email or password is invalid!", []string{}))
    return
  }

  apiSession := &ApiSession{ UserId: user.Id}
  err = db.Create(apiSession).Error

  if (err != nil) {
    r.JSON(500, Api500Envelope())
    return
  }

  data := &ApiData{ApiSession: apiSession}
  r.JSON(200, ApiEnvelope{data})
}
