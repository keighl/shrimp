package models

import (
  "testing"
  "reflect"
)

func NewTestTodo() *Todo {
  return &Todo{
    Title: "cheese",
    Complete: false,
  }
}

func Test_Todo_Table(t *testing.T) {
  x := NewTestTodo()
  expect(t, x.Table(), "todos")
}

func Test_Todo_Title_Presence(t *testing.T) {
  x := NewTestTodo()
  x.Title = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Title"], true)

  x.Title = "TOOOODOOOO"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Title"], false)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Todo_UpdateAttrs(t *testing.T) {
  obj := &Todo{
    Title: "Thing",
    Complete: true,
  }
  attrs := TodoAttrs{
    Title: "Thingz",
    Complete: false,
  }
  obj.UpdateFromAttrs(attrs)
  targetByHand := &Todo{
    Title: attrs.Title,
    Complete: attrs.Complete,
  }

  expect(t, reflect.DeepEqual(targetByHand, obj), true)
}

func Test_TodoAttrs_Todo(t *testing.T) {
  obj := &TodoAttrs{
    Title: "Thing",
    Complete: true,
  }
  targetByMethod := obj.Todo()
  targetByHand := &Todo{
    Title: obj.Title,
    Complete: obj.Complete,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}
