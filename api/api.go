package main

import (
	m "github.com/keighl/shrimp/models"
)

//////////////////////////////
// API DATA //////////////////

type Data struct {
	*m.User          `json:"current_user,omitempty"`
	APIToken         string `json:"api_token,omitempty"`
	*Error           `json:"error,omitempty"`
	*Message         `json:"message,omitempty"`
	*m.PasswordReset `json:"password_reset,omitempty"`
}

//////////////////////////////
// API MESSAGE ///////////////

type Message struct {
	Message string `json:"message,omitempty"`
}

//////////////////////////////
// API ERROR /////////////////

type Error struct {
	Message string   `json:"message,omitempty"`
	Details []string `json:"details,omitempty"`
}

func ServerErrorEnvelope(err error) Data {
	data := Data{}
	data.Error = &Error{"There was an unexpected error!", []string{}}
	ds.Error(err)
	return data
}

func ErrorEnvelope(message string, details []string) Data {
	data := Data{}
	data.Error = &Error{message, details}
	return data
}

func MessageEnvelope(message string) Data {
	data := Data{}
	data.Message = &Message{message}
	return data
}
