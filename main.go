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
  config *Configuration
  db gorm.DB
  server *martini.ClassicMartini
)

type Configuration struct {
  DBDriveSources string
  DBLoggingEnabled bool
  ServerLoggingEnabled bool
  WorkerServer string
  WorkerDatabase string
  WorkerPool string
  WorkerProcess string
}

/////////////////////////////

func main() {
  SetupApp("conf/app.json")
  workers.Start()
  server.Run() // Blocks....
  workers.Quit()
  db.Close()
}

func SetupApp(confFile string) {
  // CONFIG
  file, _ := os.Open(confFile)
  decoder := json.NewDecoder(file)
  config = &Configuration{}
  _ = decoder.Decode(config)
  file.Close()

  // GORM
  db, _ = gorm.Open("mysql", config.DBDriveSources)
  db.LogMode(config.DBLoggingEnabled)

  // MARTINI
  router := martini.NewRouter()
  mserver := martini.New()
  if (config.ServerLoggingEnabled) { mserver.Use(martini.Logger()) }
  mserver.Use(martini.Recovery())
  mserver.MapTo(router, (*martini.Routes)(nil))
  mserver.Action(router.Handle)

  server = &martini.ClassicMartini{mserver, router}

  // CORS middleware
  server.Use(cors.Allow(&cors.Options{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
    AllowHeaders: []string{"*,x-requested-with,Content-Type,If-Modified-Since,If-None-Match"},
    ExposeHeaders: []string{"Content-Length"},
  }))

  // Renderer middleware
  server.Use(render.Renderer())

  // Routes!!
  server.Get("/", RouteAuthorize, RouteUserMe)
  server.Post("/login", binding.Bind(UserAttrs{}), RouteLogin)
  server.Post("/users", binding.Bind(UserAttrs{}), RouteUserCreate)

  server.Get("/me", RouteAuthorize, RouteUserMe)
  server.Put("/me", RouteAuthorize, binding.Bind(UserAttrs{}), RouteUserUpdate)

  server.Get("/todos", RouteAuthorize, RouteTodosIndex)
  server.Post("/todos", RouteAuthorize, binding.Bind(TodoAttrs{}), RouteTodosCreate)
  server.Get("/todos/:todo_id", RouteAuthorize, RouteTodosShow)
  server.Put("/todos/:todo_id", RouteAuthorize, binding.Bind(TodoAttrs{}), RouteTodosUpdate)
  server.Delete("/todos/:todo_id", RouteAuthorize, RouteTodosDelete)

  // WORKERS
  workers.Configure(map[string]string{
    "server": config.WorkerServer,
    "database": config.WorkerDatabase,
    "pool": config.WorkerPool,
    "process": config.WorkerProcess,
  })
}

