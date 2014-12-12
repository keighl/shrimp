package api

import (
  "shrimp/models"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "net/http"
  "strings"
)

/////////////////////////////

func Authorize(c martini.Context, r render.Render, req *http.Request) {

  var token string

  token = req.Header.Get("X-API-TOKEN")
  if (token == "") {
    token = req.URL.Query().Get("api-token")
  }

  user := &models.User{}
  err := DB.Where("api_token = ?", token).First(user).Error

  if (err != nil) {
    r.JSON(401, ApiErrorEnvelope("Your token is invalid!", []string{}))
    return
  }
  c.Map(user) // Map the user to be used in the route
}

/////////////////////////////


func Login(r render.Render, attrs models.UserAttrs) {

  var err error
  var success bool
  user := &models.User{}

  err = DB.Where("email = ?", strings.TrimSpace(attrs.Email)).First(user).Error

  if (err != nil) {
    r.JSON(401, ApiErrorEnvelope("Your email or password is invalid!", []string{}))
    return
  }

  success, err = user.CheckPassword(strings.TrimSpace(attrs.Password))

  if (err != nil || !success) {
    r.JSON(401, ApiErrorEnvelope("Your email or password is invalid!", []string{}))
    return
  }

  data := &ApiData{ApiToken: user.ApiToken, CurrentUser: user}
  r.JSON(200, data)
}
