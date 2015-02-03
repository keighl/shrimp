package utils

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/cors"
  r "github.com/dancannon/gorethink"
)

func RethinkSession(conf *Configuration) (*r.Session, error) {
  return r.Connect(r.ConnectOpts{
    Address:  conf.RethinkHost,
    Database: conf.RethinkDatabase,
  })
}

func MartiniServer(conf *Configuration) (*martini.ClassicMartini) {
  router := martini.NewRouter()
  server := martini.New()
  if (conf.ServerLoggingEnabled) { server.Use(martini.Logger()) }
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



