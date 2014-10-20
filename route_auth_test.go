package main

import (
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
)

func Test_Route_Auth_Login_Failure(t *testing.T) {
  setup(t)
  body, _ := json.Marshal(UserAttrs{Email: "cheese@cheese", Password: "cheese"})
  req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Login_Success(t *testing.T) {
  setup(t)
  user, _ := UserAndSession(t)
  body, _ := json.Marshal(UserAttrs{ Email: user.Email, Password: "Password1" })
  req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

///////////

func Test_Route_Auth_Authorize_Query_Failure(t *testing.T) {
  setup(t)
  server.Get("/authTest", RouteAuthorize)
  req, _ := http.NewRequest("GET", "/authTest?session_token=cheese", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Authorize_Query_Success(t *testing.T) {
  setup(t)
  server.Get("/authTest", RouteAuthorize)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("GET", "/authTest?session_token="+apiSession.SessionToken, nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

func Test_Route_Auth_Authorize_Header_Failure(t *testing.T) {
  setup(t)
  server.Get("/authTest", RouteAuthorize)
  req, _ := http.NewRequest("GET", "/authTest", nil)
  req.Header.Set("X-SESSION-TOKEN", "cheese")
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Authorize_Header_Success(t *testing.T) {
  setup(t)
  server.Get("/authTest", RouteAuthorize)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("GET", "/authTest", nil)
  req.Header.Set("X-SESSION-TOKEN", apiSession.SessionToken)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

