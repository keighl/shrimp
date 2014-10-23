package main

import (
  "testing"
  "reflect"
)

//////////////////////////////
// API CLIENT ////////////////

func Test_ApiClient_Attrs(t *testing.T) {
  setup(t)
  client := &ApiClient{ Name: "CoolClient" }
  _ = db.Create(client)
  refute(t, client.ClientId, "")
  refute(t, client.ClientSecret, "")
}

//////////////////////////////
// API SESSION ///////////////

func Test_ApiSession_Attrs(t *testing.T) {
  setup(t)
  session := &ApiSession{}
  _ = db.Create(session)
  refute(t, session.SessionToken, "")
}

//////////////////////////////
// API ENVELOPE //////////////

func Test_Api500Envelope(t *testing.T) {
  setup(t)
  data := new(ApiData)
  data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
  targetByHand := ApiEnvelope{data}
  targetByMethod := Api500Envelope()
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}