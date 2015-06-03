package main

import (
	"github.com/go-martini/martini"
	"github.com/keighl/shrimp/comm"
	m "github.com/keighl/shrimp/models"
	"github.com/martini-contrib/render"
	"strings"
	"time"
)

/////////////////////

var savePasswordReset = func(reset *m.PasswordReset) error {
	return m.Save(reset, true)
}

func PasswordResetCreate(r render.Render, attrs m.PasswordResetAttrs) {

	user := userFromEmail(strings.TrimSpace(attrs.Email))
	if user == nil {
		r.JSON(400, ErrorEnvelope("That email isn't in our system!", []string{}))
		return
	}

	reset := &m.PasswordReset{}
	reset.UserID = user.ID
	err := savePasswordReset(reset)

	if err != nil {
		if reset.HasErrors() {
			r.JSON(400, ErrorEnvelope(err.Error(), reset.Errors))
		} else {
			r.JSON(500, ServerErrorEnvelope(err))
		}
		return
	}

	err = comm.DeliverPasswordReset(reset, user)

	if err != nil {
		r.JSON(500, ServerErrorEnvelope(err))
		return
	}

	data := &Data{PasswordReset: reset}
	r.JSON(201, data)
}

/////////////////////

var loadPasswordReset = func(id string) (*m.PasswordReset, error) {
	reset := &m.PasswordReset{}
	err := m.Find(reset, id)
	if err != nil {
		return nil, err
	}
	return reset, err
}

//////////////

func PasswordResetUpdate(params martini.Params, r render.Render, attrs m.UserAttrs) {

	reset, err := loadPasswordReset(params["token"])

	if err != nil {
		r.JSON(400, ErrorEnvelope("Invalid password reset token", nil))
		return
	}

	if reset.ExpiresAt.Before(time.Now()) {
		r.JSON(400, ErrorEnvelope("The reset token has expired", nil))
		return
	}

	if !reset.Active {
		r.JSON(400, ErrorEnvelope("The reset token has been used", nil))
		return
	}

	user, err := loadUser(reset.UserID)
	if err != nil {
		r.JSON(500, ServerErrorEnvelope(err))
		return
	}

	user.Password = attrs.Password
	user.PasswordConfirmation = attrs.PasswordConfirmation
	err = saveUser(user)

	if err != nil {
		if user.HasErrors() {
			r.JSON(400, ErrorEnvelope(err.Error(), user.Errors))
		} else {
			r.JSON(500, ServerErrorEnvelope(err))
		}
		return
	}

	reset.Active = false
	err = savePasswordReset(reset)
	if err != nil {
		r.JSON(500, ServerErrorEnvelope(err))
	}

	r.JSON(200, MessageEnvelope("Your password was reset"))
}
