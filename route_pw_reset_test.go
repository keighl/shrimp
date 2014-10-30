package main

import (
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
)

func Test_Route_PasswordReset_Create_UserNotFound(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/password-reset", binding.Bind(PasswordResetAttrs{}), MockMailerTrue, RoutePasswordResetCreate)
  body, _ := json.Marshal(PasswordResetAttrs{Email: "cheese@cheese"})
  req, _ := http.NewRequest("POST", "/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusBadRequest)
}

func Test_Route_PasswordReset_Create_Success(t *testing.T) {
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

func Test_Route_PasswordReset_Create_MailFail(t *testing.T) {
  server, recorder := testTools(t)
  server.Post("/password-reset", binding.Bind(PasswordResetAttrs{}), MockMailerFalse, RoutePasswordResetCreate)
  user, _ := UserAndSession(t)
  body, _ := json.Marshal(PasswordResetAttrs{Email: user.Email})
  req, _ := http.NewRequest("POST", "/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 500)
}

//////////

func Test_Route_PasswordReset_Update_WrongToken(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/password-reset/:token", binding.Bind(UserAttrs{}), RoutePasswordResetUpdate)
  body, _ := json.Marshal(UserAttrs{Password: "cheesed", PasswordConfirmation: "cheesed"})
  req, _ := http.NewRequest("PUT", "/password-reset/cheese", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 400)
}

func Test_Route_PasswordReset_Update_Inactive(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/password-reset/:token", binding.Bind(UserAttrs{}), RoutePasswordResetUpdate)
  user, _ := UserAndSession(t)
  pwr := &PasswordReset{UserId: user.Id}
  db.Create(pwr)
  pwr.Active = false
  db.Save(pwr)
  body, _ := json.Marshal(UserAttrs{Password: "cheesedddd", PasswordConfirmation: "cheesedddd"})
  req, _ := http.NewRequest("PUT", "/password-reset/"+string(pwr.Token), bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 400)
}

func Test_Route_PasswordReset_Update_NoUser(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/password-reset/:token", binding.Bind(UserAttrs{}), RoutePasswordResetUpdate)
  pwr := &PasswordReset{UserId: 123098}
  db.Create(pwr)
  body, _ := json.Marshal(UserAttrs{Password: "cheesedddd", PasswordConfirmation: "cheesedddd"})
  req, _ := http.NewRequest("PUT", "/password-reset/"+string(pwr.Token), bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 400)
}

func Test_Route_PasswordReset_Update_BadPassword(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/password-reset/:token", binding.Bind(UserAttrs{}), RoutePasswordResetUpdate)
  user, _ := UserAndSession(t)
  pwr := &PasswordReset{UserId: user.Id}
  db.Create(pwr)
  body, _ := json.Marshal(UserAttrs{Password: "chees", PasswordConfirmation: "cheesedddd"})
  req, _ := http.NewRequest("PUT", "/password-reset/"+string(pwr.Token), bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 400)
}

func Test_Route_PasswordReset_Update_Success(t *testing.T) {
  server, recorder := testTools(t)
  server.Put("/password-reset/:token", binding.Bind(UserAttrs{}), RoutePasswordResetUpdate)
  user, _ := UserAndSession(t)
  pwr := &PasswordReset{UserId: user.Id}
  db.Create(pwr)
  body, _ := json.Marshal(UserAttrs{Password: "cheesedddd", PasswordConfirmation: "cheesedddd"})
  req, _ := http.NewRequest("PUT", "/password-reset/"+string(pwr.Token), bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, 200)
}