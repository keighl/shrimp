package comm

import (
	"github.com/keighl/mandrill"
	c "github.com/keighl/shrimp/config"
	m "github.com/keighl/shrimp/models"
	"os"
)

var (
	Config *c.Configuration
)

func init() {
	Env(os.Getenv("SHRIMP_ENV"))
}

func Env(env string) {
	Config = c.Config(env)
	m.Env(env)
}

func SendMessage(message *mandrill.Message) error {
	client := mandrill.ClientWithKey(Config.MandrillAPIKey)
	_, err := client.MessagesSend(message)
	return err
}
