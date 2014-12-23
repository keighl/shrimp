package api

import (
  m "shrimp/models"
  r "github.com/dancannon/gorethink"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "errors"
)

////////////////////////

var loadTodos = func(u *m.User) ([]m.Todo, error) {
  todos := []m.Todo{}
  res, err := r.Table("todos").OrderBy(r.Desc("created_at")).Run(DB)
  if (err != nil) { return nil, err }
  err = res.All(&todos)
  return todos, err
}

func TodosIndex(r render.Render, user *m.User) {
  todos, err := loadTodos(user)

  if (err != nil) {
    r.JSON(500, ApiMessageEnvelope(err.Error()))
    return
  }

  data := &ApiData{User: user, Todos: todos}
  r.JSON(200, data)
}

/////////////////////////

var loadTodo = func(id string, u *m.User) (*m.Todo, error) {
  todo := &m.Todo{}
  res, err := r.Table("todos").Get(id).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(todo)
  if (err != nil) { return nil, err }
  if (todo.UserId != u.Id) { return nil, errors.New("Not your todo") }
  return todo, err
}

func TodosShow(params martini.Params, r render.Render, user *m.User) {
  todo, err := loadTodo(params["todo_id"], user)

  if (todo == nil || err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, data)
}

/////////////////////////

var saveTodo = func(todo *m.Todo) (error) {
  return todo.Save()
}

func TodosCreate(r render.Render, user *m.User, attrs m.TodoAttrs) {
  todo := attrs.Todo()
  todo.UserId = user.Id
  err := saveTodo(todo)

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(201, data)
}

/////////////////////////

func TodosUpdate(params martini.Params, r render.Render, user *m.User, attrs m.TodoAttrs) {
  todo, err := loadTodo(params["todo_id"], user)

  if (todo == nil || err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  err = saveTodo(todo)

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, ApiErrorEnvelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, Api500Envelope())
    }
    return
  }

  data := &ApiData{CurrentUser: user, Todo: todo}
  r.JSON(200, data)
}

/////////////////////////

var deleteTodo = func(todo *m.Todo) (error) {
  return todo.Delete()
}

func TodosDelete(params martini.Params, r render.Render, user *m.User) {
  todo, err := loadTodo(params["todo_id"], user)

  if (todo == nil || err != nil) {
    r.JSON(404, ApiMessageEnvelope("Record not found"))
    return
  }

  err = deleteTodo(todo)

  if (err != nil) {
    r.JSON(400, ApiMessageEnvelope("Couldn't delete the item"))
    return
  }

  data := &ApiData{CurrentUser: user, ApiMessage: &ApiMessage{"The todo was deleted"}}
  r.JSON(200, data)
}
