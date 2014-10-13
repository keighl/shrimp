package main

import (
  "time"
  "github.com/dchest/uniuri"
  "database/sql"
)

type ApiClient struct {
  Id sql.NullInt64
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