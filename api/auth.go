package main

import (
	"github.com/go-martini/martini"
	m "github.com/keighl/shrimp/models"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
)

///////////////////////////////

var userFromToken = func(token string) *m.User {
	user := &m.User{}
	err := m.FindByIndex(user, "api_token", token)
	if err != nil {
		return nil
	}
	return user
}

func Authorize(c martini.Context, r render.Render, req *http.Request) {
	token := req.Header.Get("X-API-TOKEN")
	if token == "" {
		token = req.URL.Query().Get("api-token")
	}

	user := userFromToken(token)

	if user == nil {
		r.JSON(401, ErrorEnvelope("Your token is invalid!", []string{}))
		return
	}

	c.Map(user)
}

/////////////////////////////

var userFromEmail = func(email string) *m.User {
	user := &m.User{}
	err := m.FindByIndex(user, "email", email)
	if err != nil {
		return nil
	}
	return user
}

func Login(r render.Render, attrs m.UserAttrs) {
	user := userFromEmail(strings.TrimSpace(attrs.Email))
	if user == nil {
		r.JSON(401, ErrorEnvelope("Your email or password is invalid!", []string{}))
		return
	}

	success, err := user.CheckPassword(strings.TrimSpace(attrs.Password))

	if err != nil || !success {
		r.JSON(401, ErrorEnvelope("Your email or password is invalid!", []string{}))
		return
	}

	data := &Data{APIToken: user.APIToken, User: user}
	r.JSON(200, data)
}
