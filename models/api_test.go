package models

import (
  "testing"
)

//////////////////////////////
// API CLIENT ////////////////

func Test_ApiClient_Attrs(t *testing.T) {
  setup(t)
  client := &ApiClient{ Name: "CoolClient" }
  _ = DB.Create(client)
  refute(t, client.ClientId, "")
  refute(t, client.ClientSecret, "")
}

//////////////////////////////
// API SESSION ///////////////

func Test_ApiSession_Attrs(t *testing.T) {
  setup(t)
  session := &ApiSession{}
  _ = DB.Create(session)
  refute(t, session.SessionToken, "")
}
