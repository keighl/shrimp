package main

type ApiError struct {
  Message string `json:"message,omitempty"`
  Details []string `json:"details,omitempty"`
}
