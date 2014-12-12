package api

import (
  "shrimp/models"
  "fmt"
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/modocache/gory"
  "github.com/jinzhu/gorm"
  "github.com/martini-contrib/binding"
)

// INDEX ///////////////////

func Test_Route_Todos_Index_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/todos", Authorize, TodosIndex)
  req, _ := http.NewRequest("GET", "/v1/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Index_Empty_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/todos", Authorize, TodosIndex)
  user := Uzer(t)
  req, _ := http.NewRequest("GET", "/v1/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

func Test_Route_Todos_Index_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/todos", Authorize, TodosIndex)
  user := Uzer(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  req, _ := http.NewRequest("GET", "/v1/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// SHOW ///////////////////

func Test_Route_Todos_Show_NotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/todos/:todo_id", Authorize, TodosShow)
  user := Uzer(t)
  req, _ := http.NewRequest("GET", "/v1/todos/pppp", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusNotFound)
}

func Test_Route_Todos_Show_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/todos/:todo_id", Authorize, TodosShow)
  req, _ := http.NewRequest("GET", "/v1/todos/12", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Show_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/todos/:todo_id", Authorize, TodosShow)
  user := Uzer(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/todos/%v", todo.Id), nil)
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// CREATE ///////////////////

func Test_Route_Todos_Create_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/v1/todos", Authorize, binding.Bind(models.TodoAttrs{}), TodosCreate)
  req, _ := http.NewRequest("POST", "/v1/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Create_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/v1/todos", Authorize, binding.Bind(models.TodoAttrs{}), TodosCreate)
  user := Uzer(t)
  req, _ := http.NewRequest("POST", "/v1/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_Todos_Create_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/v1/todos", Authorize, binding.Bind(models.TodoAttrs{}), TodosCreate)
  user := Uzer(t)
  body, _ := json.Marshal(gory.Build("todoAttrs"))
  req, _ := http.NewRequest("POST", "/v1/todos", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusCreated)
}

// UPDATE ///////////////////

func Test_Route_Todos_Update_NotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/v1/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  user := Uzer(t)
  req, _ := http.NewRequest("PUT", "/v1/todos/pppp", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusNotFound)
}

func Test_Route_Todos_Update_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/v1/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  req, _ := http.NewRequest("PUT", "/v1/todos/12", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Update_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/v1/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  user := Uzer(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  body, _ := json.Marshal(models.TodoAttrs{Title: "   "})
  req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/todos/%v", todo.Id), bytes.NewReader(body))
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_Todos_Update_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/v1/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  user := Uzer(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  body, _ := json.Marshal(models.TodoAttrs{Title: "HOTTTTTNESS"})
  req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/todos/%d", todo.Id), bytes.NewReader(body))
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// DELETE ///////////////////

func Test_Route_Todos_Delete_NotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Delete("/v1/todos/:todo_id", Authorize, TodosDelete)
  user := Uzer(t)
  req, _ := http.NewRequest("DELETE", "/v1/todos/pppp", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusNotFound)
}

func Test_Route_Todos_Delete_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Delete("/v1/todos/:todo_id", Authorize, TodosDelete)
  req, _ := http.NewRequest("DELETE", "/v1/todos/12", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Delete_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Delete("/v1/todos/:todo_id", Authorize, TodosDelete)
  user := Uzer(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/todos/%d", todo.Id), nil)
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
  err := DB.Where("id = ?", todo.Id).First(todo).Error
  expect(t, err, gorm.RecordNotFound)
}