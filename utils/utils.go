package utils

import (
  "os"
  "encoding/json"
  "github.com/jinzhu/gorm"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/cors"
  _ "github.com/go-sql-driver/mysql"
)

type Configuration struct {
  AppName string
  BaseURL string
  DBDriveSources string
  DBLoggingEnabled bool
  ServerLoggingEnabled bool
  WorkerServer string
  WorkerDatabase string
  WorkerPool string
  WorkerProcess string
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

func DBForConfig(conf *Configuration) gorm.DB {
  d, err := gorm.Open("mysql", conf.DBDriveSources)
  if (err != nil) { panic(err) }
  d.LogMode(conf.DBLoggingEnabled)
  return d
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

