package main

import (
  "time"
  "errors"
  "strings"
)

type Todo struct {
  Errors []string `json:"errors,omitempty" sql:"-"`
  ErrorMap map[string]bool `json:"-" sql:"-"`
  Id int64 `json:"-"`
  CreatedAt time.Time `json:"created_at,omitempty"`
  UpdatedAt time.Time `json:"updated_at,omitempty"`
  Title string `json:"title"`
  UserId int64 `json:"-"`
  Complete bool `json:"complete"`
}

func (x Todo) TableName() string {
  return "todos"
}

type TodoAttrs struct {
  Title string `json:"title" form:"title"`
  Complete bool `json:"complete" form:"complete"`
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *Todo) BeforeSave() (err error) {
  x.Errors    = []string{}
  x.ErrorMap  = map[string]bool{}
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()

  x.Trimspace()
  x.ValidateTitle()

  if (x.HasErrors()) {
    err = errors.New("There was a problem saving your todo.")
  }

  return
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Todo) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *Todo) ValidateTitle() {
  if (x.Title == "") {
    x.Errors = append(x.Errors, "Title can't be blank.")
    x.ErrorMap["Title"] = true
  }
}

func (x *Todo) Trimspace() {
  x.Title = strings.TrimSpace(x.Title)
}

//////////////////////////////
// OTHER /////////////////////

func (x *TodoAttrs) Todo() (*Todo) {
  return &Todo{
    Title: x.Title,
    Complete: x.Complete,
  }
}

func (x *Todo) TodoAttrs() (*TodoAttrs) {
  return &TodoAttrs{
    Title: x.Title,
    Complete: x.Complete,
  }
}

