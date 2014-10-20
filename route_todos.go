package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
)

func RouteTodosIndex(r render.Render, user *User) {
  data := &ApiData{User: user}
  r.JSON(200, ApiEnvelope{data})
}

func RouteTodosShow(params martini.Params, r render.Render, user *User) {
  var err error
  todo := &Todo{}
  err = db.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, Error500Envelope())
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, ApiEnvelope{data})
}

func RouteTodosCreate(r render.Render, user *User, attrs TodoAttrs) {
  var err error
  todo := attrs.Todo()
  todo.UserId = user.Id
  err = db.Create(todo).Error

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, Error400Envelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, Error500Envelope())
    }
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(201, ApiEnvelope{data})
}

func RouteTodosUpdate(params martini.Params, r render.Render, user *User, attrs TodoAttrs) {
  var err error
  todo := &Todo{}
  err = db.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, Error500Envelope())
    return
  }

  err = db.Model(todo).Updates(attrs).Error

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, Error400Envelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, Error500Envelope())
    }
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, ApiEnvelope{data})
}

