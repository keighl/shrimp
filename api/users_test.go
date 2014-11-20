package api

import (
  "shrimp/models"
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/modocache/gory"
  "github.com/martini-contrib/binding"
)

// CREATE ///////////////////

func Test_Route_Users_Create_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/users", binding.Bind(models.UserAttrs{}), UserCreate)
  req, _ := http.NewRequest("POST", "/users", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_Users_Create_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/users", binding.Bind(models.UserAttrs{}), UserCreate)
  body, _ := json.Marshal(gory.Build("userAttrs"))
  req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusCreated)
}

// SHOW ///////////////////

func Test_Route_Users_Me_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/me", Authorize, UserMe)
  req, _ := http.NewRequest("GET", "/me", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Users_Me_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/me", Authorize, UserMe)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("GET", "/me", nil)
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

// UPDATE ///////////////////

func Test_Route_Users_Update_Unauthorized(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/me", Authorize, binding.Bind(models.UserAttrs{}), UserUpdate)
  req, _ := http.NewRequest("PUT", "/me", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Users_Update_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/me", Authorize, binding.Bind(models.UserAttrs{}), UserUpdate)
  _, apiSession := UserAndSession(t)
  body, _ := json.Marshal(models.UserAttrs{Email: "cheese"})
  req, _ := http.NewRequest("PUT", "/me", bytes.NewReader(body))
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_Users_Update_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/me", Authorize, binding.Bind(models.UserAttrs{}), UserUpdate)
  user, apiSession := UserAndSession(t)
  body, _ := json.Marshal(user.UserAttrs())
  req, _ := http.NewRequest("PUT", "/me", bytes.NewReader(body))
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

