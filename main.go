package main

import (
  // "fmt"
  "os"
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/cors"
  "github.com/jinzhu/gorm"
  "github.com/jrallison/go-workers"
  _ "github.com/go-sql-driver/mysql"
)

var db gorm.DB
var conf *Configuration

type Configuration struct {
  DBDriveSources string
  Debug bool
}

/////////////////////////////

func RouteHome(r render.Render, user *User) {
  // Trigger a background job, just for fun
  workers.Enqueue("dummyQueue", "Add", user.Id.Int64)

  data := &ApiData{User: user}
  r.JSON(200, ApiEnvelope{data})
  return
}

func main() {

  var err error

  conffile, err := os.Open("conf/app.json")
  if err != nil { panic(err) }
  decoder := json.NewDecoder(conffile)
  conf = &Configuration{}
  err = decoder.Decode(&conf)
  if err != nil { panic(err) }
  conffile.Close()

  // e.g root:@tcp(localhost:3306)/shrimp_development?charset=utf8&parseTime=True
  db, err = gorm.Open("mysql", conf.DBDriveSources)
  if err != nil { panic(err) }
  defer db.Close()
  if (conf.Debug) {
    db.LogMode(true)
  }

  m := NewMartiniServer()

  ConfigureWorkerServer()
  workers.Start()

  m.Run() // Blocks....
  workers.Quit()
}

////////////////////////////

func NewMartiniServer() *martini.ClassicMartini {
  // Martini
  m := martini.Classic()

  // CORS middleware
  m.Use(cors.Allow(&cors.Options{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
    AllowHeaders: []string{"*,x-requested-with,Content-Type,If-Modified-Since,If-None-Match"},
    ExposeHeaders: []string{"Content-Length"},
  }))

  // Renderer middleware
  m.Use(render.Renderer())

  // Routes!!
  m.Get("/", RouteAuthorize, RouteHome)
  m.Post("/login", binding.Bind(UserLoginAttrs{}), RouteLogin)
  m.Post("/users", binding.Bind(UserAttrs{}), RouteUserCreate)
  m.Put("/me", RouteAuthorize, binding.Bind(UserAttrs{}), RouteUserUpdate)

  return m
}

//////////////////////////////
// ENVELOPE HELPERS //////////

func error500Envelope() (ApiEnvelope) {
  data := new(ApiData)
  data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
  return ApiEnvelope{data}
}

func error400Envelope(message string, details []string) (ApiEnvelope) {
  data := new(ApiData)
  data.ApiError = &ApiError{message, details}
  return ApiEnvelope{data}
}