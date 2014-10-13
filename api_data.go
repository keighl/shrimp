package main

type ApiData struct {
  *ApiSession `json:"session,omitempty"`
  *ApiError `json:"error,omitempty"`
  *User `json:"user,omitempty"`
}

type ApiEnvelope struct {
  *ApiData `json:"data,omitempty"`
}