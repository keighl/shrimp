package models

import (
	"reflect"
	"testing"
)

func NewTestTodo() *Todo {
	return &Todo{
		Title: "cheese",
	}
}

func Test_Todo_Table(t *testing.T) {
	x := NewTestTodo()
	expect(t, x.Table(), "todos")
}

func Test_Todo_Title_Presence(t *testing.T) {
	x := NewTestTodo()
	x.Title = ""
	err := Save(x, true)
	refute(t, err, nil)
	expect(t, x.ErrorMap["Title"], true)

	x.Title = "TOOOODOOOO"
	err = Save(x, true)
	expect(t, err, nil)
	expect(t, x.ErrorMap["Title"], false)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Todo_UpdateAttrs(t *testing.T) {
	obj := &Todo{
		Title: "Thing",
	}
	attrs := TodoAttrs{
		Title: "Thingz",
	}
	obj.UpdateFromAttrs(attrs)
	targetByHand := &Todo{
		Title: attrs.Title,
	}

	expect(t, reflect.DeepEqual(targetByHand, obj), true)
}

func Test_TodoAttrs_Todo(t *testing.T) {
	obj := &TodoAttrs{
		Title: "Thing",
	}
	targetByMethod := obj.Todo()
	targetByHand := &Todo{
		Title: obj.Title,
	}
	expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}
