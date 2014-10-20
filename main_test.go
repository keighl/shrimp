package main

import (
  "net/http"
  "testing"
)

func Test_Home_Unauthorized(t *testing.T) {
  setup(t)
  req, _ := http.NewRequest("GET", "/", nil)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Home_Authorized(t *testing.T) {
  setup(t)
  _, apiSession := UserAndSession(t)
  req, _ := http.NewRequest("GET", "/?session_token="+apiSession.SessionToken, nil)
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}