package api

import (
  m "shrimp/models"
  "github.com/martini-contrib/render"
  r "github.com/dancannon/gorethink"
)

var loadUser = func(id string) (*m.User, error) {
  user := &m.User{}
  res, err := r.Table("users").Get(id).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(user)
  return user, err
}

//////////////////////////////////////

var saveUser = func(user *m.User) (error) {
  return m.Save(user)
}

func UserCreate(r render.Render, attrs m.UserAttrs) {
  user := attrs.User()
  err := saveUser(user)

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := &Data{User: user, APIToken: user.APIToken}
  r.JSON(201, data)
}

//////////////////////////////////////

func Me(r render.Render, user *m.User) {
  data := &Data{User: user}
  r.JSON(200, data)
}

//////////////////////////////////////

func MeUpdate(r render.Render, user *m.User, attrs m.UserAttrs) {
  user.UpdateFromAttrs(attrs)
  err := saveUser(user)

  if (err != nil) {
    if (user.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), user.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := &Data{User: user}
  r.JSON(200, data)
}

