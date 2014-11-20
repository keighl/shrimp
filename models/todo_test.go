package models

import (
  "testing"
  "github.com/modocache/gory"
  "reflect"
)

func Test_Todo_Name_Presence(t *testing.T) {
  setup(t)
  todo := gory.Build("todo").(*Todo)
  todo.Title = ""
  err := DB.Create(todo).Error
  refute(t, err, nil)
  expect(t, todo.ErrorMap["Title"], true)
  todo.Title = "TOOOODOOOO"
  err = DB.Create(todo).Error
  expect(t, err, nil)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Todo_TodoAttrs(t *testing.T) {
  setup(t)
  obj := gory.Build("todo").(*Todo)
  targetByMethod := obj.TodoAttrs()
  targetByHand := &TodoAttrs{
    Title: obj.Title,
    Complete: obj.Complete,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}

func Test_TodoAttrs_Todo(t *testing.T) {
  setup(t)
  obj := gory.Build("todoAttrs").(*TodoAttrs)
  targetByMethod := obj.Todo()
  targetByHand := &Todo{
    Title: obj.Title,
    Complete: obj.Complete,
  }
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}
