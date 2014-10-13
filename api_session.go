package main

import (
  "time"
  "github.com/dchest/uniuri"
  "database/sql"
)

type ApiSession struct {
  Id sql.NullInt64 `json:"-"`
  ApiClient ApiClient `json:"-"`
  ApiClientId sql.NullInt64 `json:"-"`
  UserId sql.NullInt64 `json:"-"`
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