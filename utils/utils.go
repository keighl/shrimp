package utils

import (
  "os"
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/cors"
  _ "github.com/go-sql-driver/mysql"
  r "github.com/dancannon/gorethink"
)

type Configuration struct {
  AppName string
  BaseURL string
  RethinkHost string
  RethinkDatabase string
  ServerLoggingEnabled bool
  MandrillAPIKey string
}

func ConfigForFile(confFile string) *Configuration {
  file, err := os.Open(confFile)
  if (err != nil) { panic(err) }
  decoder := json.NewDecoder(file)
  c := &Configuration{}
  err = decoder.Decode(c)
  if (err != nil) { panic(err) }
  file.Close()
  return c
}

func RethinkSession(conf *Configuration) *r.Session {
  session, _ := r.Connect(r.ConnectOpts{
    Address:  conf.RethinkHost,
    Database: conf.RethinkDatabase,
  })
  return session
}

func MartiniServer(logginEnabled bool) (*martini.ClassicMartini) {
  router := martini.NewRouter()
  server := martini.New()
  if (logginEnabled) { server.Use(martini.Logger()) }
  server.Use(martini.Recovery())
  server.MapTo(router, (*martini.Routes)(nil))
  server.Action(router.Handle)
  s := &martini.ClassicMartini{server, router}
  s.Use(render.Renderer())
  s.Use(cors.Allow(&cors.Options{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
    AllowHeaders: []string{"*", "x-requested-with", "Content-Type", "If-Modified-Since", "If-None-Match", "X-API-TOKEN"},
    ExposeHeaders: []string{"Content-Length"},
  }))
  return s
}
