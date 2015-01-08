package models

import (
  "testing"
  "reflect"
  u "shrimp/utils"
  r "github.com/dancannon/gorethink"
  "github.com/dchest/uniuri"
)

func init() {
  Config = u.ConfigForFile("../config/test.json")
  DB = u.RethinkSession(Config)

  _, _ = r.Table("users").Delete().RunWrite(DB)
  _, _ = r.Table("password_resets").Delete().RunWrite(DB)
  _, _ = r.Table("todos").Delete().RunWrite(DB)
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

//////////////////////
//////////////////////
//////////////////////

func Test_Record_Save_NewRecord_Success(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)
  refute(t, x.Id, "")

  res, err := r.Table(x.Table()).Get(x.Id).Run(DB)
  expect(t, err, nil)
  y := &Record{}
  err = res.One(y)
  expect(t, err, nil)
  expect(t, y.Id, x.Id)
}

func Test_Record_Save_NewRecord_Success_PresetId(t *testing.T) {
  x := &Record{}
  id := uniuri.NewLen(30)
  x.Id = id
  err := Save(x)
  expect(t, err, nil)
  expect(t, x.Id, id)
}

type RecordFAILInvalid struct {
  Record
}
func (x *RecordFAILInvalid) Validate() {
  x.Record.Validate()
  x.ErrorOn("Thing", "Isn't good!")
}
func Test_Record_Save_NewRecord_ErrorValidation(t *testing.T) {
  x := &RecordFAILInvalid{}
  err := Save(x)
  refute(t, err, nil)
}

type RecordFAILDB struct {
  Record
}
func (x *RecordFAILDB) Table() string  {
  return "nonexistent_table"
}
func Test_Record_Save_NewRecord_ErrorDB(t *testing.T) {
  x := &RecordFAILDB{}
  err := Save(x)
  refute(t, err, nil)
}

func Test_Record_Save_ExistingRecord_Success(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)

  err = Save(x)
  expect(t, err, nil)
}

func Test_Record_Save_ExistingRecord_ErrorValidation(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)

  y := &RecordFAILInvalid{}
  y.Id = y.Id
  err = Save(y)
  refute(t, err, nil)
}

func Test_Record_Save_ExistingRecord_ErrorDB(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)

  y := &RecordFAILDB{}
  y.Id = x.Id
  err = Save(y)
  refute(t, err, nil)
}

func Test_Record_Delete_Success(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)

  err = Delete(x)
  expect(t, err, nil)

  res, _ := r.Table(x.Table()).Get(x.Id).Run(DB)
  expect(t, res.IsNil(), true)
}

func Test_Record_Delete_FailDB(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)

  y := &RecordFAILDB{}
  x.Id = x.Id
  err = Delete(y)
  refute(t, err, nil)
}

//////////////////////
//////////////////////
//////////////////////

func Test_Record_NewRecord(t *testing.T) {
  x := &Record{}
  expect(t, x.IsNewRecord(), true)

  x.Id = "CHESE"
  expect(t, x.IsNewRecord(), false)
}

func Test_Record_CreatedAt(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
  expect(t, x.UpdatedAt, x.CreatedAt)
}

func Test_Record_UpdatedAt(t *testing.T) {
  x := &Record{}
  err := Save(x)
  expect(t, err, nil)
  expect(t, x.UpdatedAt, x.CreatedAt)
  ua := x.UpdatedAt
  ca := x.CreatedAt
  err = Save(x)
  expect(t, err, nil)
  refute(t, x.UpdatedAt, ua)
  expect(t, x.CreatedAt, ca)
}

func Test_Record_BeforeValidate(t *testing.T) {
  x := &Record{}
  x.BeforeValidate()
  expect(t, len(x.Errors), 0)
  expect(t, len(x.ErrorMap), 0)
}

func Test_Record_HasErrors(t *testing.T) {
  x := &Record{}
  x.BeforeValidate()
  expect(t, x.HasErrors(), false)
  x.ErrorOn("Cheese", "ipsun")
  expect(t, x.HasErrors(), true)
}

func Test_Record_SetId(t *testing.T) {
  x := &Record{}
  x.SetId("CHEESE")
  expect(t, x.Id, "CHEESE")
}

func Test_Record_GetId(t *testing.T) {
  x := &Record{}
  x.Id = "CHEESE"
  expect(t, x.GetId(), "CHEESE")
}


