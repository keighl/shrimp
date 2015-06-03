package main

import (
	"github.com/go-martini/martini"
	"github.com/jrallison/go-workers"
	"github.com/keighl/drschollz"
	c "github.com/keighl/shrimp/config"
	m "github.com/keighl/shrimp/models"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"log"
	"os"
)

var (
	Config = c.Config(os.Getenv("SHRIMP_ENV"))
	Logger = log.New(os.Stdout, "api: ", log.Ldate|log.Ltime)
	ds     = &drschollz.Queue{}
)

func init() {
	workers.Configure(Config.WorkerConfig)
	drschollz.Conf.MandrillAPIKey = Config.MandrillAPIKey
	drschollz.Conf.EmailsTo = []string{"devs@example.com"}
	drschollz.Conf.EmailFrom = "errors@example.com"
	drschollz.Conf.AppName = "Shrimp API"
}

func main() {
	server := MartiniServer()
	SetupServerRoutes(server)

	var err error
	ds, err = drschollz.Start(2)
	if err != nil {
		panic(err)
	}
	server.Run()
	ds.Stop()
}

func SetupServerRoutes(server *martini.ClassicMartini) {

	// Default
	server.Get("/v1/", Authorize, Me)

	// Login
	server.Post("/v1/login", binding.Bind(m.UserAttrs{}), Login)

	// Signup
	server.Post("/v1/users", binding.Bind(m.UserAttrs{}), UserCreate)

	// Me
	server.Get("/v1/users/me", Authorize, Me)
	server.Put("/v1/users/me", Authorize, binding.Bind(m.UserAttrs{}), MeUpdate)

	// Todos
	server.Get("/v1/todos", Authorize, TodosIndex)
	server.Post("/v1/todos", Authorize, binding.Bind(m.TodoAttrs{}), TodosCreate)
	server.Get("/v1/todos/:todo_id", Authorize, TodosShow)
	server.Put("/v1/todos/:todo_id", Authorize, binding.Bind(m.TodoAttrs{}), TodosUpdate)
	server.Delete("/v1/todos/:todo_id", Authorize, TodosDelete)

	// Password Reset
	server.Post("/v1/password-reset", binding.Bind(m.PasswordResetAttrs{}), PasswordResetCreate)
	server.Post("/v1/password-reset/:token", binding.Bind(m.UserAttrs{}), PasswordResetUpdate)
}

func MartiniServer() *martini.ClassicMartini {
	router := martini.NewRouter()
	server := martini.New()
	if Config.Env != "test" {
		server.Use(martini.Logger())
	}
	server.Use(martini.Recovery())
	server.MapTo(router, (*martini.Routes)(nil))
	server.Action(router.Handle)
	s := &martini.ClassicMartini{server, router}
	s.Use(render.Renderer())
	s.Use(cors.Allow(&cors.Options{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:  []string{"*", "x-requested-with", "Content-Type", "If-Modified-Since", "If-None-Match", "X-API-TOKEN"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	return s
}
