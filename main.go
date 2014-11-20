package main

import (
  "shrimp/api"
  "shrimp/models"
  "shrimp/utils"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
  _ "github.com/jrallison/go-workers"
  _ "github.com/go-sql-driver/mysql"
)

/////////////////////////////

func main() {
  config := utils.ConfigForFile("conf/app.json")
  DB     := utils.DBForConfig(config)

  api.DB        = DB
  models.DB     = DB
  api.Config    = config
  models.Config = config

  server := utils.MartiniServer(config.ServerLoggingEnabled)

  SetupServerRoutes(server)

  // workers.Configure(map[string]string{
  //   "server": config.WorkerServer,
  //   "database": config.WorkerDatabase,
  //   "pool": config.WorkerPool,
  //   "process": config.WorkerProcess,
  // })
  // workers.Start()

  server.Run() // Blocks....

  // workers.Quit()
  DB.Close()
}

func SetupServerRoutes(server *martini.ClassicMartini) {
  server.Get("/", api.Authorize, api.UserMe)

  server.Post("/login", binding.Bind(models.UserAttrs{}), api.Login)
  server.Post("/users", binding.Bind(models.UserAttrs{}), api.UserCreate)
  server.Get("/me", api.Authorize, api.UserMe)
  server.Put("/me", api.Authorize, binding.Bind(models.UserAttrs{}), api.UserUpdate)
  server.Get("/todos", api.Authorize, api.TodosIndex)
  server.Post("/todos", api.Authorize, binding.Bind(models.TodoAttrs{}), api.TodosCreate)
  server.Get("/todos/:todo_id", api.Authorize, api.TodosShow)
  server.Put("/todos/:todo_id", api.Authorize, binding.Bind(models.TodoAttrs{}), api.TodosUpdate)
  server.Delete("/todos/:todo_id", api.Authorize, api.TodosDelete)
  server.Post("/password-reset", binding.Bind(models.PasswordResetAttrs{}), api.Mailer, api.PasswordResetCreate)
  server.Post("/password-reset/:token", binding.Bind(models.UserAttrs{}), api.PasswordResetUpdate)
}


