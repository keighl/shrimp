package api

import (
  m "shrimp/models"
  r "github.com/dancannon/gorethink"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "net/http"
  "strings"
)

///////////////////////////////

var userFromToken = func(token string) *m.User {
  user := &m.User{}
  res, err := r.Table("users").GetAllByIndex("api_token", token).Run(DB)
  if (err != nil) { return nil }
  err = res.One(user)
  if (err != nil) { return nil }
  return user
}

func Authorize(c martini.Context, r render.Render, req *http.Request) {
  token := req.Header.Get("X-API-TOKEN")
  if (token == "") {
    token = req.URL.Query().Get("api-token")
  }

  user := userFromToken(token)

  if (user == nil) {
    r.JSON(401, ErrorEnvelope("Your token is invalid!", []string{}))
    return
  }

  c.Map(user)
}

/////////////////////////////

var userFromEmail = func(email string) *m.User {
  user := &m.User{}
  res, err := r.Table("users").GetAllByIndex("email", email).Run(DB)
  if (err != nil) { return nil }
  err = res.One(user)
  if (err != nil) { return nil }
  return user
}

func Login(r render.Render, attrs m.UserAttrs) {
  user := userFromEmail(strings.TrimSpace(attrs.Email))
  if (user == nil) {
    r.JSON(401, ErrorEnvelope("Your email or password is invalid!", []string{}))
    return
  }

  success, err := user.CheckPassword(strings.TrimSpace(attrs.Password))

  if (err != nil || !success) {
    r.JSON(401, ErrorEnvelope("Your email or password is invalid!", []string{}))
    return
  }

  data := &Data{APIToken: user.APIToken, User: user}
  r.JSON(200, data)
}
