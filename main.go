package main

import (
  "shrimp/api"
  m "shrimp/models"
  u "shrimp/utils"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
)

/////////////////////////////

func main() {
  config := u.ConfigForFile("config/app.json")
  DB := u.RethinkSession(config)
  api.DB = DB
  m.DB = DB
  api.Config = u.ConfigForFile("config/app.json")
  server := u.MartiniServer(config.ServerLoggingEnabled)
  SetupServerRoutes(server)
  server.Run() // Blocks....
  // DB.Close()
}

func SetupServerRoutes(server *martini.ClassicMartini) {

  server.Get("/v1/", api.Authorize, api.Me)

  // Login
  server.Post("/v1/login", binding.Bind(m.UserAttrs{}), api.Login)

  // Signup
  server.Post("/v1/users", binding.Bind(m.UserAttrs{}), api.UserCreate)

  // Me
  server.Get("/v1/users/me", api.Authorize, api.Me)
  server.Put("/v1/users/me", api.Authorize, binding.Bind(m.UserAttrs{}), api.MeUpdate)

  // Todos
  server.Get("/v1/todos", api.Authorize, api.TodosIndex)
  server.Post("/v1/todos", api.Authorize, binding.Bind(m.TodoAttrs{}), api.TodosCreate)
  server.Get("/v1/todos/:todo_id", api.Authorize, api.TodosShow)
  server.Put("/v1/todos/:todo_id", api.Authorize, binding.Bind(m.TodoAttrs{}), api.TodosUpdate)
  server.Delete("/v1/todos/:todo_id", api.Authorize, api.TodosDelete)

  // Password Reset
  server.Post("/v1/password-reset", binding.Bind(m.PasswordResetAttrs{}), api.PasswordResetCreate)
  server.Post("/v1/password-reset/:token", binding.Bind(m.UserAttrs{}), api.PasswordResetUpdate)
}


