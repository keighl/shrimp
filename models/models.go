package models

import (
  u "shrimp/utils"
  r "github.com/dancannon/gorethink"
  "errors"
  "time"
)

var  (
  Config *u.Configuration
  DB *r.Session
)

type Recorder interface {
  Table() string

  GetId() string
  SetId(string)
  IsNewRecord() bool

  // Validation
  BeforeValidate()
  AfterValidate()
  Validate()
  HasErrors() (bool)
  ErrorOn(attr string, message string)

  // Save
  BeforeSave()
  AfterSave()

  // Create
  BeforeCreate()
  AfterCreate()

  // Update
  BeforeUpdate()
  AfterUpdate()

  // Delete
  BeforeDelete()
  AfterDelete()
}

func Save(x Recorder) error {

  x.BeforeValidate()
  x.Validate()
  x.AfterValidate()

  if (x.HasErrors()) {
    return errors.New("Validation errors")
  }

  x.BeforeSave()

  if (x.IsNewRecord()) {
    x.BeforeCreate()
    res, err := r.Table(x.Table()).Insert(x).RunWrite(DB)
    if (err != nil) { return err }
    if (len(res.GeneratedKeys) > 0) {
      x.SetId(res.GeneratedKeys[0])
    }
    x.AfterCreate()
    x.AfterSave()
    return nil
  }

  x.BeforeUpdate()
  _, err := r.Table(x.Table()).Get(x.GetId()).Replace(x).RunWrite(DB)
  if (err != nil) { return err }
  x.AfterUpdate()
  x.AfterSave()
  return nil
}

func Delete(x Recorder) error {
  x.BeforeDelete()
  _, err := r.Table(x.Table()).Get(x.GetId()).Delete().RunWrite(DB)
  if (err != nil) { return err }
  x.AfterDelete()
  return nil
}

//////////////////////////////
//////////////////////////////

type Record struct {
  Id string `gorethink:"id,omitempty" json:"id"`
  CreatedAt time.Time `gorethink:"created_at" json:"created_at"`
  UpdatedAt time.Time `gorethink:"updated_at" json:"-"`
  Errors []string `gorethink:"-" json:"errors,omitempty"`
  ErrorMap map[string]bool `gorethink:"-" json:"-"`
}

func (x *Record) Table() string {
  return "records"
}

func (x *Record) IsNewRecord() bool {
  return x.GetId() == ""
}

//////////////////////////////
// ID ////////////////////////

func (x *Record) GetId() string {
  return x.Id
}

func (x *Record) SetId(id string) {
  x.Id = id
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Record) BeforeValidate() {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
}

func (x *Record) AfterValidate() {}

func (x *Record) Validate() {}

func (x *Record) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *Record) ErrorOn(attr string, message string) {
  x.ErrorMap[attr] = true
  x.Errors = append(x.Errors, message)
}

//////////////////////////////
// CALLBACK //////////////////

func (x *Record) BeforeSave() {
  x.UpdatedAt = time.Now()
}

func (x *Record) AfterSave() {}

func (x *Record) BeforeCreate() {
  x.CreatedAt = time.Now()
  x.UpdatedAt = x.CreatedAt
}

func (x *Record) AfterCreate() {}

func (x *Record) BeforeUpdate() {}

func (x *Record) AfterUpdate() {}

func (x *Record) BeforeDelete() {}

func (x *Record) AfterDelete() {}
