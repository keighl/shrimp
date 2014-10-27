package main

import (
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
)

func Test_Route_PasswordReset_UserNotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/password-reset", binding.Bind(PasswordResetAttrs{}), MockMailerTrue, RoutePasswordResetCreate)
  body, _ := json.Marshal(PasswordResetAttrs{Email: "cheese@cheese"})
  req, _ := http.NewRequest("POST", "/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_PasswordReset_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/password-reset", binding.Bind(PasswordResetAttrs{}), MockMailerTrue, RoutePasswordResetCreate)
  user, _ := UserAndSession(t)
  body, _ := json.Marshal(PasswordResetAttrs{Email: user.Email})
  req, _ := http.NewRequest("POST", "/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusCreated)
  pwr := &PasswordReset{}
  err := db.Where("id = ?", user.Id).First(pwr).Error
  expect(t, err, nil)
}

func Test_Route_PasswordReset_MailFail(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/password-reset", binding.Bind(PasswordResetAttrs{}), MockMailerFalse, RoutePasswordResetCreate)
  user, _ := UserAndSession(t)
  body, _ := json.Marshal(PasswordResetAttrs{Email: user.Email})
  req, _ := http.NewRequest("POST", "/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 500)
}

