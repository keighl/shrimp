package main

import (
	"github.com/jrallison/go-workers"
	"github.com/keighl/drschollz"
	c "github.com/keighl/shrimp/config"
	m "github.com/keighl/shrimp/models"
	"log"
	"os"
)

var (
	Config *c.Configuration
	Logger = log.New(os.Stdout, "workshop: ", log.Ldate|log.Ltime|log.Lshortfile)
	ds     = &drschollz.Queue{}
)

func init() {
	Config = c.Config(os.Getenv("SHRIMP_ENV"))
	workers.Configure(Config.WorkerConfig)

	workers.Process("do_something", DoSomething, 1)

	drschollz.Conf.MandrillAPIKey = Config.MandrillAPIKey
	drschollz.Conf.EmailsTo = []string{"kyle@seasalt.io"}
	drschollz.Conf.EmailFrom = "errors@seasalt.io"
	drschollz.Conf.AppName = "Shrimp Queue"
}

func main() {

	var err error
	ds, err = drschollz.Start(2)
	if err != nil {
		panic(err)
	}

	workers.Run()
	ds.Stop()
}

func DoSomething(message *workers.Msg) {
	userID, err := message.Args().String()
	if err != nil {
		panic(err)
	}
	user := &m.User{}
	err = m.Find(user, userID)
	if err != nil {
		ds.Error(err)
		panic(err)
	}

	Logger.Println("---> DID SOMETHING", user.ID)
}
