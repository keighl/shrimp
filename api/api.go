package api

import (
  m "shrimp/models"
  u "shrimp/utils"
  "github.com/keighl/mandrill"
  r "github.com/dancannon/gorethink"
)

var  (
  Config *u.Configuration
  DB *r.Session
)

//////////////////////////////
// API DATA //////////////////

type ApiData struct {
  CurrentUser *m.User `json:"current_user,omitempty"`
  ApiToken string `json:"api_token,omitempty"`
  *ApiError `json:"error,omitempty"`
  *ApiMessage `json:"message,omitempty"`
  *m.User `json:"user,omitempty"`
  *m.Todo `json:"todo,omitempty"`
  *m.PasswordReset `json:"password_reset,omitempty"`
  Todos []m.Todo `json:"todos,omitempty"`
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

//////////////////////////////
// MAILER ////////////////////

var sendEmail = func(message *mandrill.Message) (bool) {
  client := mandrill.ClientWithKey(Config.MandrillAPIKey)
  _, apiError, err := client.MessagesSend(message)
  if (apiError != nil || err != nil) { return false }
  return true
}
