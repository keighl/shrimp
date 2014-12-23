package api

import (
  m "shrimp/models"
  "net/http"
  "testing"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
  "errors"
  "github.com/keighl/mandrill"
  "time"
)

func passwordResetCreateRunner(t *testing.T, code int) {
  sendEmail = func(message *mandrill.Message) (bool) { return true }
  server, recorder := testTools(t)
  server.Post("/v1/password-reset", binding.Bind(m.PasswordResetAttrs{}), PasswordResetCreate)
  body, _ := json.Marshal(m.PasswordResetAttrs{Email: "cheese@cheese"})
  req, _ := http.NewRequest("POST", "/v1/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_PasswordReset_Create_UserNotFound(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    return nil
  }
  passwordResetCreateRunner(t, http.StatusBadRequest)
}

func Test_Route_PasswordReset_Create_Success(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    return &m.User{}
  }

  savePasswordReset = func(reset *m.PasswordReset) error {
    return nil
  }
  passwordResetCreateRunner(t, http.StatusCreated)
}

func Test_Route_PasswordReset_Create_400(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    return &m.User{}
  }

  savePasswordReset = func(reset *m.PasswordReset) error {
    reset.Errors = []string{"Something went wrong..."}
    return errors.New("*******")
  }
  passwordResetCreateRunner(t, http.StatusBadRequest)
}

func Test_Route_PasswordReset_Create_500(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    return &m.User{}
  }

  savePasswordReset = func(reset *m.PasswordReset) error {
    return errors.New("*******")
  }
  passwordResetCreateRunner(t, http.StatusInternalServerError)
}

func Test_Route_PasswordReset_Create_EmailFail(t *testing.T) {
  userFromEmail = func(email string) *m.User {
    return &m.User{}
  }

  savePasswordReset = func(reset *m.PasswordReset) error {
    return nil
  }

  saveUser = func(user *m.User) (error) {
    return nil
  }

  sendEmail = func(message *mandrill.Message) (bool) { return false }
  server, recorder := testTools(t)
  server.Post("/v1/password-reset", binding.Bind(m.PasswordResetAttrs{}), PasswordResetCreate)
  body, _ := json.Marshal(m.PasswordResetAttrs{Email: "cheese@cheese"})
  req, _ := http.NewRequest("POST", "/v1/password-reset", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, http.StatusInternalServerError)
}

////////////////

func resetPasswordUpdateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Put("/v1/password-reset/:token", binding.Bind(m.UserAttrs{}), PasswordResetUpdate)
  body, _ := json.Marshal(m.UserAttrs{Password: "cheese", PasswordConfirmation: "cheese"})
  req, _ := http.NewRequest("PUT", "/v1/password-reset/XXXXXX", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_PasswordReset_Update_WrongToken(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    return nil, errors.New("******")
  }
  resetPasswordUpdateRunner(t, http.StatusBadRequest)
}

func Test_Route_PasswordReset_Update_Expired(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    pwr := &m.PasswordReset{ExpiresAt: time.Now().Add(-6*time.Hour)}
    return pwr, nil
  }
  resetPasswordUpdateRunner(t, http.StatusBadRequest)
}

func Test_Route_PasswordReset_Update_NoUser(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    pwr := &m.PasswordReset{ExpiresAt: time.Now().Add(6*time.Hour)}
    return pwr, nil
  }

  loadUser = func(id string) (*m.User, error) {
    return nil, errors.New("**********")
  }

  resetPasswordUpdateRunner(t, http.StatusInternalServerError)
}

func Test_Route_PasswordReset_Update_Success(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    pwr := &m.PasswordReset{ExpiresAt: time.Now().Add(6*time.Hour)}
    return pwr, nil
  }

  loadUser = func(id string) (*m.User, error) {
    return &m.User{}, nil
  }

  saveUser = func(user *m.User) (error) {
    return nil
  }

  resetPasswordUpdateRunner(t, http.StatusOK)

}
func Test_Route_PasswordReset_Update_Fail400(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    pwr := &m.PasswordReset{ExpiresAt: time.Now().Add(6*time.Hour)}
    return pwr, nil
  }

  loadUser = func(id string) (*m.User, error) {
    return &m.User{}, nil
  }

  saveUser = func(user *m.User) (error) {
    user.Errors = []string{"SOmething went wrong!"}
    return errors.New("********")
  }

  resetPasswordUpdateRunner(t, http.StatusBadRequest)

}

func Test_Route_PasswordReset_Update_Fail500(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    pwr := &m.PasswordReset{ExpiresAt: time.Now().Add(6*time.Hour)}
    return pwr, nil
  }

  loadUser = func(id string) (*m.User, error) {
    return &m.User{}, nil
  }

  saveUser = func(user *m.User) (error) {
    return errors.New("********")
  }

  resetPasswordUpdateRunner(t, http.StatusInternalServerError)
}

func Test_Route_PasswordReset_Update_FailMessage(t *testing.T) {
  loadPasswordReset = func(token string) (*m.PasswordReset, error) {
    expect(t, token, "XXXXXX")
    pwr := &m.PasswordReset{ExpiresAt: time.Now().Add(6*time.Hour)}
    return pwr, nil
  }

  loadUser = func(id string) (*m.User, error) {
    return &m.User{}, nil
  }

  saveUser = func(user *m.User) (error) {
    return errors.New("********")
  }

  resetPasswordUpdateRunner(t, http.StatusInternalServerError)
}


