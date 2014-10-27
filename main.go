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
  "github.com/keighl/mandrill"
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
  MandrillAPIKey string
}

/////////////////////////////

func main() {
  config = ConfigForFile("conf/app.json")
  db     = DBForConfig(config)

  workers.Configure(map[string]string{
    "server": config.WorkerServer,
    "database": config.WorkerDatabase,
    "pool": config.WorkerPool,
    "process": config.WorkerProcess,
  })
  workers.Start()

  server = martiniServer(config.ServerLoggingEnabled)
  setupServerRoutes(server)
  server.Run() // Blocks....

  workers.Quit()
  db.Close()
}

func ConfigForFile(confFile string) *Configuration {
  file, _ := os.Open(confFile)
  decoder := json.NewDecoder(file)
  c := &Configuration{}
  _ = decoder.Decode(c)
  file.Close()
  return c
}

func DBForConfig(conf *Configuration) gorm.DB {
  d, _ := gorm.Open("mysql", conf.DBDriveSources)
  d.LogMode(conf.DBLoggingEnabled)
  return d
}

func martiniServer(logginEnabled bool) (*martini.ClassicMartini) {
  router := martini.NewRouter()
  mserver := martini.New()
  if (logginEnabled) { mserver.Use(martini.Logger()) }
  mserver.Use(martini.Recovery())
  mserver.MapTo(router, (*martini.Routes)(nil))
  mserver.Action(router.Handle)
  s := &martini.ClassicMartini{mserver, router}
  s.Use(render.Renderer())
  s.Use(cors.Allow(&cors.Options{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
    AllowHeaders: []string{"*,x-requested-with,Content-Type,If-Modified-Since,If-None-Match"},
    ExposeHeaders: []string{"Content-Length"},
  }))
  return s
}

func setupServerRoutes(*martini.ClassicMartini) {
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
  server.Post("/password-reset", binding.Bind(PasswordResetAttrs{}), Mailer, RoutePasswordResetCreate)
}

// MAILER INJECTION /////////////

type SendEmail func(message *mandrill.Message)(bool)

func Mailer(c martini.Context) {
  c.Map(SendEmail(func (message *mandrill.Message)(bool) {
    client := mandrill.ClientWithKey(config.MandrillAPIKey)
    _, apiError, err := client.MessagesSend(message)
    if (apiError != nil || err != nil) {
      return false
    }
    return true
  }))
}
