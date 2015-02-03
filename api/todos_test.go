package api

import (
  "testing"
  m "github.com/keighl/shrimp/models"
  "net/http"
  "errors"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
)

//////////////////////////////////////
// INDEX ///////////////////

func todosIndexRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/todos", AuthorizeOK, TodosIndex)
  req, _ := http.NewRequest("GET", "/v1/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Todos_Index_Error(t *testing.T) {
  loadTodos = func(u *m.User) ([]m.Todo, error) {
    return nil, errors.New("*****")
  }
  todosIndexRunner(t, http.StatusInternalServerError)
}

func Test_Route_Todos_Index_Success(t *testing.T) {
  loadTodos = func(u *m.User) ([]m.Todo, error) {
    return []m.Todo{m.Todo{}}, nil
  }
  todosIndexRunner(t, http.StatusOK)
}

//////////////////////////////////////
// SHOW ///////////////////

func todosShowRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/todos/:todo_id", AuthorizeOK, TodosShow)
  req, _ := http.NewRequest("GET", "/v1/todos/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Todos_Show_NotFound(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }
  todosShowRunner(t, http.StatusNotFound)
}

func Test_Route_Todos_Show_Success(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return &m.Todo{}, nil
  }

  todosShowRunner(t, http.StatusOK)
}

//////////////////////////////////////
// CREATE ///////////////////

func todosCreateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Post("/v1/todos", AuthorizeOK, binding.Bind(m.TodoAttrs{}), TodosCreate)
  body, _ := json.Marshal(m.TodoAttrs{})
  req, _ := http.NewRequest("POST", "/v1/todos", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Todos_Create_Failure_400(t *testing.T) {
  saveTodo = func(todo *m.Todo) (error) {
    todo.Errors = []string{"Something went wrong!"}
    return errors.New("*******")
  }
  todosCreateRunner(t, http.StatusBadRequest)
}

func Test_Route_Todos_Create_Failure_500(t *testing.T) {
  saveTodo = func(todo *m.Todo) (error) {
    return errors.New("*******")
  }
  todosCreateRunner(t, http.StatusInternalServerError)
}

func Test_Route_Todos_Create_Success(t *testing.T) {
  saveTodo = func(todo *m.Todo) (error) {
    return nil
  }
  todosCreateRunner(t, http.StatusCreated)
}

//////////////////////////////////////
// UPDATE ///////////////////

func todosUpdateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Put("/v1/todos/:todo_id", AuthorizeOK, binding.Bind(m.TodoAttrs{}), TodosUpdate)
  body, _ := json.Marshal(m.TodoAttrs{})
  req, _ := http.NewRequest("PUT", "/v1/todos/XXXXXX", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Todos_Update_NotFound(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }

  todosUpdateRunner(t, http.StatusNotFound)
}

func Test_Route_Todos_Update_Failure_400(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return &m.Todo{}, nil
  }

  saveTodo = func(todo *m.Todo) (error) {
    todo.Errors = []string{"Something went wrong!"}
    return errors.New("*********")
  }

  todosUpdateRunner(t, http.StatusBadRequest)
}

func Test_Route_Todos_Update_Failure_500(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return &m.Todo{}, nil
  }

  saveTodo = func(todo *m.Todo) (error) {
    return errors.New("*********")
  }

  todosUpdateRunner(t, http.StatusInternalServerError)
}

func Test_Route_Todos_Update_Success(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return &m.Todo{}, nil
  }

  saveTodo = func(todo *m.Todo) (error) {
    return nil
  }

  todosUpdateRunner(t, http.StatusOK)
}

//////////////////////////////////////
// DELETE ///////////////////

func todosDeleteRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Delete("/v1/todos/:todo_id", AuthorizeOK, TodosDelete)
  req, _ := http.NewRequest("DELETE", "/v1/todos/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Todos_Delete_NotFound(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return nil, nil
  }

  todosDeleteRunner(t, http.StatusNotFound)
}

func Test_Route_Todos_Delete_Success(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return &m.Todo{}, nil
  }

  deleteTodo = func(todo *m.Todo) (error) {
    return nil
  }

  todosDeleteRunner(t, http.StatusOK)
}

func Test_Route_Todos_Delete_400(t *testing.T) {
  loadTodo = func(id string, u *m.User) (*m.Todo, error) {
    expect(t, id, "XXXXXX")
    return &m.Todo{}, nil
  }

  deleteTodo = func(todo *m.Todo) (error) {
    return errors.New("*******")
  }

  todosDeleteRunner(t, http.StatusBadRequest)
}

