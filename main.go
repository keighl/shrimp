package main

import (
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

var  (
  db gorm.DB
  conf *Configuration
)

type Configuration struct {
  DBDriveSources string
  Debug bool
}

/////////////////////////////

func RouteHome(r render.Render, user *User) {
  data := &ApiData{User: user, CurrentUser: user}
  r.JSON(200, ApiEnvelope{data})
  return
}

/////////////////////////////

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

  ConfigureWorkerServer(false)
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
  m.Post("/login", binding.Bind(UserAttrs{}), RouteLogin)
  m.Post("/users", binding.Bind(UserAttrs{}), RouteUserCreate)
  m.Put("/me", RouteAuthorize, binding.Bind(UserAttrs{}), RouteUserUpdate)

  m.Get("/todos", RouteAuthorize, RouteTodosIndex)
  m.Post("/todos", RouteAuthorize, binding.Bind(TodoAttrs{}), RouteTodosCreate)
  m.Get("/todos/:todo_id", RouteAuthorize, RouteTodosShow)
  m.Put("/todos/:todo_id", RouteAuthorize, binding.Bind(TodoAttrs{}), RouteTodosUpdate)
  // m.Delete("/todos/:todo_id", RouteAuthorize, RouteTodoDestroy)

  return m
}

