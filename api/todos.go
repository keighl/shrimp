package api

import (
  m "shrimp/models"
  r "github.com/dancannon/gorethink"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "errors"
)

type TodoData struct {
  *m.User `json:"user,omitempty"`
  *m.Todo `json:"todo"`
}

type TodosData struct {
  *m.User `json:"user,omitempty"`
  Todos []m.Todo `json:"todos"`
}

////////////////////////

var loadTodos = func(u *m.User) ([]m.Todo, error) {
  todos := []m.Todo{}
  res, err := r.Table("todos").GetAllByIndex(u.ID, "user_id").OrderBy(r.Desc("created_at")).Run(DB)
  if (err != nil) { return nil, err }
  err = res.All(&todos)
  return todos, err
}

func TodosIndex(r render.Render, user *m.User) {
  todos, err := loadTodos(user)

  if (err != nil) {
    r.JSON(500, MessageEnvelope(err.Error()))
    return
  }

  data := TodosData{User: user, Todos: todos}
  r.JSON(200, data)
}

/////////////////////////

var loadTodo = func(id string, u *m.User) (*m.Todo, error) {
  todo := &m.Todo{}
  res, err := r.Table("todos").Get(id).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(todo)
  if (err != nil) { return nil, err }
  if (todo.UserID != u.ID) { return nil, errors.New("Not your todo") }
  return todo, err
}

func TodosShow(params martini.Params, r render.Render, user *m.User) {
  todo, err := loadTodo(params["todo_id"], user)

  if (todo == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  data := TodoData{User: user, Todo: todo}
  r.JSON(200, data)
}

/////////////////////////

var saveTodo = func(todo *m.Todo) (error) {
  return m.Save(todo)
}

func TodosCreate(r render.Render, user *m.User, attrs m.TodoAttrs) {
  todo := attrs.Todo()
  todo.UserID = user.ID
  err := saveTodo(todo)

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := TodoData{User: user, Todo: todo}
  r.JSON(201, data)
}

/////////////////////////

func TodosUpdate(params martini.Params, r render.Render, user *m.User, attrs m.TodoAttrs) {
  todo, err := loadTodo(params["todo_id"], user)

  if (todo == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  todo.UpdateFromAttrs(attrs)
  err = saveTodo(todo)

  if (err != nil) {
    if (todo.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), todo.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := TodoData{User: user, Todo: todo}
  r.JSON(200, data)
}

/////////////////////////

var deleteTodo = func(todo *m.Todo) (error) {
  return m.Delete(todo)
}

func TodosDelete(params martini.Params, r render.Render, user *m.User) {
  todo, err := loadTodo(params["todo_id"], user)

  if (todo == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  err = deleteTodo(todo)

  if (err != nil) {
    r.JSON(400, MessageEnvelope("Couldn't delete the item"))
    return
  }

  data := Data{User: user, Message: &Message{"The todo was deleted"}}
  r.JSON(200, data)
}
