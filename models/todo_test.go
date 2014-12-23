package models

import (
  "testing"
  "reflect"
)

func NewTodo() *Todo {
  return &Todo{
    Title: "cheese",
  }
}

func Test_Todo_Title_Presence(t *testing.T) {
  x := &Todo{}
  x.Title = ""
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Title"], true)

  x.Title = "TOOOODOOOO"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Title"], false)
}

func Test_Todo_BeforeCreate(t *testing.T) {
  x := &Todo{}
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

func Test_Todo_BeforeUpdate(t *testing.T) {
  x := &Todo{}
  x.BeforeUpdate()
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

func Test_Todo_Create_Success(t *testing.T) {
  setup(t)

  x := NewTodo()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")
}

func Test_Todo_Create_Fail(t *testing.T) {
  setup(t)

  x := NewTodo()
  x.Title  = ""
  err := x.Save()
  refute(t, err, nil)
  expect(t, x.Id, "")
}

func Test_Todo_Update_Success(t *testing.T) {
  setup(t)

  x := NewTodo()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Save()
  expect(t, err, nil)
}

func Test_Todo_Update_Fail(t *testing.T) {
  setup(t)

  x := NewTodo()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  x.Title = ""
  err = x.Save()
  refute(t, err, nil)
}

func Test_Todo_Delete(t *testing.T) {
  setup(t)

  x := NewTodo()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Delete()
  expect(t, err, nil)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Todo_TodoAttrs(t *testing.T) {
  obj := &Todo{
    Title: "Thing",
    Complete: true,
  }
  targetByMethod := obj.TodoAttrs()
  targetByHand := &TodoAttrs{
    Title: obj.Title,
    Complete: obj.Complete,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
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
