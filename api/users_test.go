package api

import (
  m "shrimp/models"
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "errors"
  "github.com/martini-contrib/binding"
)

// CREATE ///////////////////

func createUserRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Post("/v1/users", binding.Bind(m.UserAttrs{}), UserCreate)
  body, _ := json.Marshal(m.UserAttrs{})
  req, _ := http.NewRequest("POST", "/v1/users", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Users_Create_400(t *testing.T) {
  saveUser = func(user *m.User) (error) {
    user.Errors = []string{"Something went wrong!"}
    return errors.New("*****")
  }

  createUserRunner(t, http.StatusBadRequest)
}

func Test_Route_Users_Create_500(t *testing.T) {
  saveUser = func(user *m.User) (error) {
    return errors.New("*****")
  }

  createUserRunner(t, http.StatusInternalServerError)
}

func Test_Route_Users_Create_Success(t *testing.T) {
  saveUser = func(user *m.User) (error) {
    return nil
  }
  createUserRunner(t, http.StatusCreated)
}

// SHOW ///////////////////

func meRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/users/me", AuthorizeOK, Me)
  body, _ := json.Marshal(m.UserAttrs{})
  req, _ := http.NewRequest("GET", "/v1/users/me", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Me_Success(t *testing.T) {
  meRunner(t, http.StatusOK)
}

// UPDATE ///////////////////

func updateUserRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Put("/v1/users/me", Authorize, binding.Bind(m.UserAttrs{}), MeUpdate)
  body, _ := json.Marshal(m.UserAttrs{})
  req, _ := http.NewRequest("PUT", "/v1/users/me", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Users_Update_400(t *testing.T) {
  saveUser = func(user *m.User) (error) {
    user.Errors = []string{"Something went wrong!"}
    return errors.New("*****")
  }

  updateUserRunner(t, http.StatusBadRequest)
}

func Test_Route_Users_Update_500(t *testing.T) {
  saveUser = func(user *m.User) (error) {
    return errors.New("*****")
  }

  updateUserRunner(t, http.StatusInternalServerError)
}

func Test_Route_Users_Update_Success(t *testing.T) {
  saveUser = func(user *m.User) (error) {
    return nil
  }

  updateUserRunner(t, http.StatusOK)
}

