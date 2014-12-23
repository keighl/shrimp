package models

import (
  u "shrimp/utils"
  r "github.com/dancannon/gorethink"
)

var  (
  Config *u.Configuration
  DB *r.Session
)