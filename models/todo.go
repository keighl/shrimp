package models

import (
  "time"
  "strings"
  r "github.com/dancannon/gorethink"
  "errors"
)

type Todo struct {
  Errors []string `gorethink:"-" json:"errors,omitempty" sql:"-"`
  ErrorMap map[string]bool `gorethink:"-" json:"-" sql:"-"`
  Id string `gorethink:"id,omitempty" json:"id"`
  CreatedAt time.Time `gorethink:"created_at" json:"created_at,omitempty"`
  UpdatedAt time.Time `gorethink:"updated_at" json:"updated_at,omitempty"`
  Title string `gorethink:"title" json:"title"`
  UserId string `gorethink:"user_id" json:"-"`
  Complete bool `gorethink:"complete" json:"complete"`
}

type TodoAttrs struct {
  Title string `json:"title" form:"title"`
  Complete bool `json:"complete" form:"complete"`
}

//////////////////////////////
// TRANSACTIONS //////////////

func (x *Todo) Save() error {

  if (!x.Validate()) {
    return errors.New("Validation errors")
  }

  if (x.Id == "") {
    x.BeforeCreate()
    res, err := r.Table("todos").Insert(x).RunWrite(DB)
    if (err != nil) { return err }
    x.Id = res.GeneratedKeys[0]
  }

  x.BeforeUpdate()
  _, err := r.Table("todos").Get(x.Id).Replace(x).RunWrite(DB)
  return err
}

func (x *Todo) Delete() error {
  _, err := r.Table("todos").Get(x.Id).Delete().RunWrite(DB)
  return err
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *Todo) BeforeCreate() {
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()
}

func (x *Todo) BeforeUpdate() {
  x.UpdatedAt = time.Now()
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Todo) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.Trimspace()
  x.ValidateTitle()
  return !x.HasErrors()
}

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

