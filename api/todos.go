package api

import (
  "shrimp/models"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/jinzhu/gorm"
)

func TodosIndex(r render.Render, user *models.User) {
  todos := []models.Todo{}
  err := DB.Where("user_id = ?", user.Id).Find(&todos).Error

  if (err != nil) {
    if (err != gorm.RecordNotFound) {
      r.JSON(500, ApiMessageEnvelope(err.Error()))
      return
    }
  }

  data := &ApiData{User: user, Todos: todos}
  r.JSON(200, ApiEnvelope{data})
}

func TodosShow(params martini.Params, r render.Render, user *models.User) {
  todo := &models.Todo{}
  err := DB.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, ApiEnvelope{data})
}

func TodosCreate(r render.Render, user *models.User, attrs models.TodoAttrs) {
  todo := attrs.Todo()
  todo.UserId = user.Id
  err := DB.Create(todo).Error

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

func TodosUpdate(params martini.Params, r render.Render, user *models.User, attrs models.TodoAttrs) {
  todo := &models.Todo{}
  err := DB.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  err = DB.Model(todo).Updates(attrs).Error

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

func TodosDelete(params martini.Params, r render.Render, user *models.User) {
  todo := &models.Todo{}
  err := DB.Where("id = ?", params["todo_id"]).Where("user_id = ?", user.Id).First(todo).Error

  if (err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  err = DB.Delete(todo).Error

  if (err != nil) {
    r.JSON(400, ApiMessageEnvelope("Couldn't delete the item"))
    return
  }

  data := &ApiData{CurrentUser: user, ApiMessage: &ApiMessage{"The todo was deleted"}}
  r.JSON(200, ApiEnvelope{data})
}
