package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/jinzhu/gorm"
)

func RouteTodosIndex(r render.Render, user *User) {
  todos := []Todo{}
  err := db.Where("user_id = ?", user.Id).Find(&todos).Error

  if (err != nil) {
    if (err != gorm.RecordNotFound) {
      r.JSON(500, ApiMessageEnvelope(err.Error()))
      return
    }
  }

  data := &ApiData{User: user, Todos: todos}
  r.JSON(200, ApiEnvelope{data})
}

func RouteTodosShow(params martini.Params, r render.Render, user *User) {
  todo := &Todo{}
  err := db.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, ApiEnvelope{data})
}

func RouteTodosCreate(r render.Render, user *User, attrs TodoAttrs) {
  todo := attrs.Todo()
  todo.UserId = user.Id
  err := db.Create(todo).Error

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(201, ApiEnvelope{data})
}

func RouteTodosUpdate(params martini.Params, r render.Render, user *User, attrs TodoAttrs) {
  todo := &Todo{}
  err := db.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  err = db.Model(todo).Updates(attrs).Error

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, ApiEnvelope{data})
}

func RouteTodosDelete(params martini.Params, r render.Render, user *User) {
  todo := &Todo{}
  err := db.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  err = db.Delete(todo).Error

  if (err != nil) {
    r.JSON(400, ApiMessageEnvelope("Couldn't delete the item"))
    return
  }

  data := &ApiData{CurrentUser: user, ApiMessage: &ApiMessage{"The todo was deleted"}}
  r.JSON(200, ApiEnvelope{data})
}
