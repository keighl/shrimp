package main

import (
  "time"
  "github.com/dchest/uniuri"
  "database/sql"
)

//////////////////////////////
// API CLIENT ////////////////

type ApiClient struct {
  Id int64
  ClientId string
  ClientSecret string
  Name string
  CreatedAt time.Time
  UpdatedAt time.Time
}

func (x ApiClient) TableName() string {
  return "api_clients"
}

func (x *ApiClient) BeforeCreate() (err error) {
  x.CreatedAt    = time.Now()
  x.UpdatedAt    = time.Now()
  x.ClientId     = uniuri.NewLen(50)
  x.ClientSecret = uniuri.NewLen(50)
  return
}

//////////////////////////////
// API SESSION ///////////////

type ApiSession struct {
  Id int64 `json:"-"`
  ApiClient ApiClient `json:"-"`
  ApiClientId sql.NullInt64 `json:"-"`
  UserId int64 `json:"-"`
  SessionToken string `json:"token"`
  CreatedAt time.Time `json:"-"`
  UpdatedAt time.Time `json:"-"`
}

func (x ApiSession) TableName() string {
  return "api_sessions"
}

func (x *ApiSession) BeforeCreate() (err error) {
  x.CreatedAt    = time.Now()
  x.UpdatedAt    = time.Now()
  x.SessionToken = uniuri.NewLen(50)
  return
}

//////////////////////////////
// API DATA //////////////////

type ApiData struct {
  CurrentUser *User `json:"current_user,omitempty"`
  *ApiSession `json:"session,omitempty"`
  *ApiError `json:"error,omitempty"`
  *User `json:"user,omitempty"`
  *Todo `json:"todo,omitempty"`
}

//////////////////////////////
// API ERROR /////////////////

type ApiError struct {
  Message string `json:"message,omitempty"`
  Details []string `json:"details,omitempty"`
}

//////////////////////////////
// API ENVELOPE //////////////

type ApiEnvelope struct {
  *ApiData `json:"data,omitempty"`
}

func Error500Envelope() (ApiEnvelope) {
  data := new(ApiData)
  data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
  return ApiEnvelope{data}
}

func Error400Envelope(message string, details []string) (ApiEnvelope) {
  data := new(ApiData)
  data.ApiError = &ApiError{message, details}
  return ApiEnvelope{data}
}

