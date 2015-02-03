package main

import (
  "os"
  "github.com/keighl/shrimp/api"
  m "github.com/keighl/shrimp/models"
  u "github.com/keighl/shrimp/utils"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
  r "github.com/dancannon/gorethink"
)

/////////////////////////////

func main() {
  env := os.Getenv("MARTINI_ENV")
  config := u.Config(env)
  DB, _ := u.RethinkSession(config)
  SetupSubpackages(DB, config)
  server := u.MartiniServer(config)
  SetupServerRoutes(server)
  server.Run()
  DB.Close()
}

func SetupSubpackages(DB *r.Session, config *u.Configuration) {
  api.DB = DB
  m.DB = DB
  api.Config = config
  m.Config = config
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


