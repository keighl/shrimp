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
  server.Get("/todos", Authorize, TodosIndex)
  req, _ := http.NewRequest("GET", "/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Index_Empty_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/todos", Authorize, TodosIndex)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("GET", "/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

func Test_Route_Todos_Index_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/todos", Authorize, TodosIndex)
  user, apiSession := UserAndSession(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  req, _ := http.NewRequest("GET", "/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// SHOW ///////////////////

func Test_Route_Todos_Show_NotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/todos/:todo_id", Authorize, TodosShow)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("GET", "/todos/pppp", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusNotFound)
}

func Test_Route_Todos_Show_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/todos/:todo_id", Authorize, TodosShow)
  req, _ := http.NewRequest("GET", "/todos/12", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Show_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/todos/:todo_id", Authorize, TodosShow)
  user, apiSession := UserAndSession(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  req, _ := http.NewRequest("GET", fmt.Sprintf("/todos/%v", todo.Id), nil)
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// CREATE ///////////////////

func Test_Route_Todos_Create_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/todos", Authorize, binding.Bind(models.TodoAttrs{}), TodosCreate)
  req, _ := http.NewRequest("POST", "/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Create_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/todos", Authorize, binding.Bind(models.TodoAttrs{}), TodosCreate)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("POST", "/todos", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_Todos_Create_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/todos", Authorize, binding.Bind(models.TodoAttrs{}), TodosCreate)
  _, apiSession := UserAndSession(t)
  body, _ := json.Marshal(gory.Build("todoAttrs"))
  req, _ := http.NewRequest("POST", "/todos", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusCreated)
}

// UPDATE ///////////////////

func Test_Route_Todos_Update_NotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("PUT", "/todos/pppp", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusNotFound)
}

func Test_Route_Todos_Update_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  req, _ := http.NewRequest("PUT", "/todos/12", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Update_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  user, apiSession := UserAndSession(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  body, _ := json.Marshal(models.TodoAttrs{Title: "   "})
  req, _ := http.NewRequest("PUT", fmt.Sprintf("/todos/%v", todo.Id), bytes.NewReader(body))
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_Todos_Update_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/todos/:todo_id", Authorize, binding.Bind(models.TodoAttrs{}), TodosUpdate)
  user, apiSession := UserAndSession(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  body, _ := json.Marshal(models.TodoAttrs{Title: "HOTTTTTNESS"})
  req, _ := http.NewRequest("PUT", fmt.Sprintf("/todos/%d", todo.Id), bytes.NewReader(body))
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// DELETE ///////////////////

func Test_Route_Todos_Delete_NotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Delete("/todos/:todo_id", Authorize, TodosDelete)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("DELETE", "/todos/pppp", nil)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusNotFound)
}

func Test_Route_Todos_Delete_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Delete("/todos/:todo_id", Authorize, TodosDelete)
  req, _ := http.NewRequest("DELETE", "/todos/12", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Todos_Delete_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Delete("/todos/:todo_id", Authorize, TodosDelete)
  user, apiSession := UserAndSession(t)
  todo := gory.Build("todo").(*models.Todo)
  todo.UserId = user.Id
  _ = DB.Create(todo).Error
  req, _ := http.NewRequest("DELETE", fmt.Sprintf("/todos/%d", todo.Id), nil)
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
  err := DB.Where("id = ?", todo.Id).First(todo).Error
  expect(t, err, gorm.RecordNotFound)
}