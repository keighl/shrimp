package main

import (
	"errors"
	"github.com/go-martini/martini"
	"github.com/keighl/shrimp/comm"
	c "github.com/keighl/shrimp/config"
	m "github.com/keighl/shrimp/models"
	"net/http/httptest"
	"reflect"
	"testing"
)

func init() {
	Config = c.Config("test")
	m.Env("test")
	comm.Env("test")
}

func testTools(t *testing.T) (*martini.ClassicMartini, *httptest.ResponseRecorder) {
	return MartiniServer(), httptest.NewRecorder()
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

//////////////////////////////
// API ENVELOPE //////////////

func Test_ServerErrorEnvelope(t *testing.T) {
	data := Data{}
	data.Error = &Error{"There was an unexpected error!", []string{}}
	targetByHand := data
	targetByMethod := ServerErrorEnvelope(errors.New("XXXX"))
	expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}
