package api

import (
  m "github.com/keighl/shrimp/models"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
)

// STUBS ///////////////////////////

func AuthorizeOK(c martini.Context) {
  user := &m.User{}
  c.Map(user)
}

func AuthorizeFAIL(r render.Render) {
  r.JSON(401, ErrorEnvelope("Your token is invalid!", []string{}))
}

// LOGIN ///////////////////////////

func loginRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Post("/v1/login", binding.Bind(m.UserAttrs{}), Login)
  body, _ := json.Marshal(m.UserAttrs{Email: "cheese@cheese", Password: "cheese"})
  req, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Auth_Login_Email_Failure(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    return nil
  }
  loginRunner(t, http.StatusUnauthorized)
}

func Test_Route_Auth_Login_Password_Failure(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    user := &m.User{}
    user.SetPassword("queso")
    return user
  }

  loginRunner(t, http.StatusUnauthorized)
}

func Test_Route_Auth_Login_Success(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    user := &m.User{}
    user.SetPassword("cheese")
    return user
  }

  loginRunner(t, http.StatusOK)
}

// AUTHORIZE ///////////////////////////

func Test_Route_Auth_Authorize_Query_Failure(t *testing.T) {
  userFromToken = func(tok string) *m.User {
    expect(t, tok, "cheese")
    return nil
  }

  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  req, _ := http.NewRequest("GET", "/v1/authTest?api-token=cheese", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Authorize_Query_Success(t *testing.T) {
  userFromToken = func(tok string) *m.User {
    expect(t, tok, "cheese")
    return &m.User{}
  }

  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  req, _ := http.NewRequest("GET", "/v1/authTest?api-token=cheese", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

func Test_Route_Auth_Authorize_Header_Failure(t *testing.T) {
  userFromToken = func(tok string) *m.User {
    expect(t, tok, "cheese")
    return nil
  }

  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  req, _ := http.NewRequest("GET", "/v1/authTest", nil)
  req.Header.Set("X-API-TOKEN", "cheese")
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusUnauthorized)
}

func Test_Route_Auth_Authorize_Header_Success(t *testing.T) {
  userFromToken = func(tok string) *m.User {
    expect(t, tok, "cheese")
    return &m.User{}
  }

  server, recorder := testTools(t)
  server.Get("/v1/authTest", Authorize)
  req, _ := http.NewRequest("GET", "/v1/authTest", nil)
  req.Header.Set("X-API-TOKEN", "cheese")
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusOK)
}

