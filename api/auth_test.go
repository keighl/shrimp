package api

import (
  "shrimp/models"
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
)

func Test_Route_Auth_Login_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/v1/login", binding.Bind(models.UserAttrs{}), Login)
  body, _ := json.Marshal(models.UserAttrs{Email: "cheese@cheese", Password: "cheese"})
  req, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Login_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/v1/login", binding.Bind(models.UserAttrs{}), Login)
  user := Uzer(t)
  body, _ := json.Marshal(models.UserAttrs{ Email: user.Email, Password: "Password1" })
  req, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

///////////

func Test_Route_Auth_Authorize_Query_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  req, _ := http.NewRequest("GET", "/v1/authTest?api-token=cheese", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Authorize_Query_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  user := Uzer(t)
  req, _ := http.NewRequest("GET", "/v1/authTest?api-token="+user.ApiToken, nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

func Test_Route_Auth_Authorize_Header_Failure(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  req, _ := http.NewRequest("GET", "/v1/authTest", nil)
  req.Header.Set("X-API-TOKEN", "cheese")
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Authorize_Header_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  user := Uzer(t)
  req, _ := http.NewRequest("GET", "/v1/authTest", nil)
  req.Header.Set("X-API-TOKEN", user.ApiToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

