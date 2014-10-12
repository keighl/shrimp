package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/binding"
  "github.com/jinzhu/gorm"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  _ "fmt"
)

var db gorm.DB

/////////////////////////////

func RouteAuthorize(c martini.Context, res http.ResponseWriter, req *http.Request) {
  var err error
  session_token := req.URL.Query().Get("session_token")
  user := User{}

  err = db.Table("users").Select("users.*").Joins("INNER JOIN api_sessions x on x.user_id = users.id").Where("session_token = ?", session_token).Limit(1).Scan(&user).Error

  if (err != nil) {
    res.WriteHeader(http.StatusUnauthorized)
    return
  }

  c.Map(&user) // Map the user to be used in the route
}

/////////////////////////////

func RouteLogin(r render.Render) {

  var err error
  var success bool
  user := User{}
  err = db.First(&user).Error

  if (err != nil) {
    data := new(ApiData)
    data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
    envelope := ApiEnvelope{data}
    r.JSON(500, envelope)
    return
  }

  success, err = user.CheckPassword("cheese")

  if (err != nil || !success) {
    data := new(ApiData)
    data.ApiError = &ApiError{"Your email or password is invalid!", []string{}}
    envelope := ApiEnvelope{data}
    r.JSON(400, envelope)
    return
  }

  apiSession := ApiSession{ UserId: user.Id }
  err = db.Create(&apiSession).Error

  if (err != nil) {
    data := new(ApiData)
    data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
    envelope := ApiEnvelope{data}
    r.JSON(500, envelope)
    return
  }

  data := new(ApiData)
  data.ApiSession = &apiSession
  envelope := ApiEnvelope{data}
  r.JSON(200, envelope)
}

/////////////////////////////

func RouteHome(r render.Render, user *User) {
  data := new(ApiData)
  data.User = user
  envelope := ApiEnvelope{data}
  r.JSON(200, envelope)
  return
}

/////////////////////////////

func RouteUserUpdate(r render.Render, user *User, attrs UserAttrs) {

  var err error
  err = db.Model(user).Updates(attrs).Error

  if (err != nil) {
    if (user.hasErrors()) {
      data := new(ApiData)
      data.ApiError = &ApiError{ err.Error(), user.Errors }
      envelope := ApiEnvelope{data}
      r.JSON(400, envelope)
      return
    } else {
      data := new(ApiData)
      data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
      envelope := ApiEnvelope{data}
      r.JSON(500, envelope)
      return
    }
  }

  data := new(ApiData)
  data.User = user
  envelope := ApiEnvelope{data}
  r.JSON(200, envelope)
}

/////////////////////////////

func RouteUserCreate(r render.Render, attrs UserAttrs) {

  var err error
  user := User {
    NameFirst: attrs.NameFirst,
    NameLast: attrs.NameLast,
    Email: attrs.Email,
    Password: attrs.Password,
    PasswordConfirmation: attrs.PasswordConfirmation, }

  err = db.Create(&user).Error

  if (err != nil) {
    if (user.hasErrors()) {
      data := new(ApiData)
      data.ApiError = &ApiError{ err.Error(), user.Errors }
      envelope := ApiEnvelope{data}
      r.JSON(400, envelope)
      return
    } else {
      data := new(ApiData)
      data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
      envelope := ApiEnvelope{data}
      r.JSON(500, envelope)
      return
    }
  }

  apiSession := ApiSession{ UserId: user.Id }
  err = db.Create(&apiSession).Error

  if (err != nil) {
    data := new(ApiData)
    data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
    envelope := ApiEnvelope{data}
    r.JSON(500, envelope)
    return
  }

  data := new(ApiData)
  data.User = &user
  data.ApiSession = &apiSession
  envelope := ApiEnvelope{data}
  r.JSON(201, envelope)
}


/////////////////////////////

func main() {

  var err error
  db, err = gorm.Open("mysql", "root:@tcp(localhost:3306)/shrimp_development?charset=utf8&parseTime=True")
  if err != nil { panic(err) }
  defer db.Close()

  db.LogMode(true)

  m := martini.Classic()

  m.Use(render.Renderer())

  m.Get("/", RouteAuthorize, RouteHome)
  m.Get("/login", RouteLogin)
  m.Put("/me", RouteAuthorize, binding.Bind(UserAttrs{}), RouteUserUpdate)
  m.Post("/users", binding.Bind(UserAttrs{}), RouteUserCreate)

  m.Run()
}