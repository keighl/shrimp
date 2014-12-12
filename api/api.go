package api

import (
  "shrimp/models"
  "shrimp/utils"
  "github.com/jinzhu/gorm"
  "github.com/keighl/mandrill"
  "github.com/go-martini/martini"
)

var  (
  Config *utils.Configuration
  DB gorm.DB
)

//////////////////////////////
// API DATA //////////////////

type ApiData struct {
  CurrentUser *models.User `json:"current_user,omitempty"`
  ApiToken string `json:"api_token,omitempty"`
  *ApiError `json:"error,omitempty"`
  *ApiMessage `json:"message,omitempty"`
  *models.User `json:"user,omitempty"`
  *models.Todo `json:"todo,omitempty"`
  *models.PasswordReset `json:"password_reset,omitempty"`
  Todos []models.Todo `json:"todos,omitempty"`
}

//////////////////////////////
// API MESSAGE ///////////////

type ApiMessage struct {
  Message string `json:"message,omitempty"`
}

//////////////////////////////
// API ERROR /////////////////

type ApiError struct {
  Message string `json:"message,omitempty"`
  Details []string `json:"details,omitempty"`
}

func Api500Envelope() (ApiData) {
  data := ApiData{}
  data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
  return data
}

func ApiErrorEnvelope(message string, details []string) (ApiData) {
  data := ApiData{}
  data.ApiError = &ApiError{message, details}
  return data
}

func ApiMessageEnvelope(message string) (ApiData) {
  data := ApiData{}
  data.ApiMessage = &ApiMessage{message}
  return data
}

// MAILER INJECTION /////////////

type SendEmail func(message *mandrill.Message)(bool)

func Mailer(c martini.Context) {
  c.Map(SendEmail(func (message *mandrill.Message)(bool) {
    client := mandrill.ClientWithKey(Config.MandrillAPIKey)
    _, apiError, err := client.MessagesSend(message)
    if (apiError != nil || err != nil) { return false }
    return true
  }))
}


