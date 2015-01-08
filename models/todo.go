package models

import (
  "strings"
)

type Todo struct {
  Record
  Title string `gorethink:"title" json:"title"`
  UserId string `gorethink:"user_id" json:"-"`
  Complete bool `gorethink:"complete" json:"complete"`
}

type TodoAttrs struct {
  Title string `json:"title" form:"title"`
  Complete bool `json:"complete" form:"complete"`
}

func (x *Todo) Table() string {
  return "todos"
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Todo) Validate() {
  x.Record.Validate()
  x.Trimspace()
  x.ValidateTitle()
}

func (x *Todo) ValidateTitle() {
  if (x.Title == "") {
    x.ErrorOn("Title", "Title can't be blank.")
  }
}

func (x *Todo) Trimspace() {
  x.Title = strings.TrimSpace(x.Title)
}

//////////////////////////////
// OTHER /////////////////////

func (x *Todo) UpdateFromAttrs(attrs TodoAttrs) {
  if (attrs.Title != "") { x.Title = attrs.Title }
  x.Complete = attrs.Complete
}

func (x *TodoAttrs) Todo() (*Todo) {
  return &Todo{
    Title: x.Title,
    Complete: x.Complete,
  }
}


