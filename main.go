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
  config := utils.ConfigForFile("config/app.json")
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


  server.Get("/v1/", api.Authorize, api.UserMe)

  // Login
  server.Post("/v1/login", binding.Bind(models.UserAttrs{}), api.Login)

  // Signup
  server.Post("/v1/users", binding.Bind(models.UserAttrs{}), api.UserCreate)

  server.Options("/v1/users/me", func() {})

  // Me
  server.Get("/v1/users/me", api.Authorize, api.UserMe)



  server.Get("/v1/me", api.Authorize, api.UserMe)
  server.Put("/v1/me", api.Authorize, binding.Bind(models.UserAttrs{}), api.UserUpdate)

  // Todos
  server.Get("/v1/todos", api.Authorize, api.TodosIndex)
  server.Post("/v1/todos", api.Authorize, binding.Bind(models.TodoAttrs{}), api.TodosCreate)
  server.Get("/v1/todos/:todo_id", api.Authorize, api.TodosShow)
  server.Put("/v1/todos/:todo_id", api.Authorize, binding.Bind(models.TodoAttrs{}), api.TodosUpdate)
  server.Delete("/v1/todos/:todo_id", api.Authorize, api.TodosDelete)

  // Password Reset
  server.Post("/v1/password-reset", binding.Bind(models.PasswordResetAttrs{}), api.Mailer, api.PasswordResetCreate)
  server.Post("/v1/password-reset/:token", binding.Bind(models.UserAttrs{}), api.PasswordResetUpdate)
}


