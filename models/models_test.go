package models

import (
	r "github.com/dancannon/gorethink"
	"github.com/dchest/uniuri"
	"reflect"
	"sync"
	"testing"
)

func init() {
	Env("test")

	r.DbDrop(Config.RethinkDatabase).Exec(DB)
	r.DbCreate(Config.RethinkDatabase).Exec(DB)

	tables := []string{
		new(Record).Table(),
		new(User).Table(),
		new(PasswordReset).Table(),
		new(Todo).Table(),
	}

	var wg sync.WaitGroup
	for _, t := range tables {
		wg.Add(1)
		go func(table string) {
			r.Db(Config.RethinkDatabase).TableCreate(table).RunWrite(DB)
			r.Table(table).Delete().RunWrite(DB)
			wg.Done()
		}(t)
	}
	wg.Wait()
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
	err := Save(x, true)
	expect(t, err, nil)
	refute(t, x.ID, "")

	res, err := r.Table(x.Table()).Get(x.ID).Run(DB)
	expect(t, err, nil)
	y := &Record{}
	err = res.One(y)
	expect(t, err, nil)
	expect(t, y.ID, x.ID)
}

func Test_Record_Save_NewRecord_Success_PresetID(t *testing.T) {
	x := &Record{}
	id := uniuri.NewLen(30)
	x.ID = id
	err := Save(x, true)
	expect(t, err, nil)
	expect(t, x.ID, id)
}

type RecordFAILInvalid struct {
	Record
}

func (x *RecordFAILInvalid) PerformValidations() {
	x.Record.PerformValidations()
	x.ErrorOn("Thing", "Isn't good!")
}
func Test_Record_Save_NewRecord_ErrorValidation(t *testing.T) {
	x := &RecordFAILInvalid{}
	err := Save(x, true)
	refute(t, err, nil)
}

type RecordFAILDB struct {
	Record
}

func (x *RecordFAILDB) Table() string {
	return "nonexistent_table"
}
func Test_Record_Save_NewRecord_ErrorDB(t *testing.T) {
	x := &RecordFAILDB{}
	err := Save(x, true)
	refute(t, err, nil)
}

func Test_Record_Save_ExistingRecord_Success(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)

	err = Save(x, true)
	expect(t, err, nil)
}

func Test_Record_Save_ExistingRecord_ErrorValidation(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)

	y := &RecordFAILInvalid{}
	y.ID = y.ID
	err = Save(y, true)
	refute(t, err, nil)
}

func Test_Record_Save_ExistingRecord_ErrorDB(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)

	y := &RecordFAILDB{}
	y.ID = x.ID
	err = Save(y, true)
	refute(t, err, nil)
}

func Test_Record_Delete_Success(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)

	err = Delete(x)
	expect(t, err, nil)

	res, _ := r.Table(x.Table()).Get(x.ID).Run(DB)
	expect(t, res.IsNil(), true)
}

func Test_Record_Delete_FailDB(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)

	y := &RecordFAILDB{}
	x.ID = x.ID
	err = Delete(y)
	refute(t, err, nil)
}

//////////////////////
//////////////////////
//////////////////////

func Test_Record_NewRecord(t *testing.T) {
	x := &Record{}
	expect(t, x.IsNewRecord(), true)

	x.ID = "CHESE"
	expect(t, x.IsNewRecord(), false)
}

func Test_Record_CreatedAt(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)
	refute(t, x.CreatedAt.Format("RFC3339"), nil)
	expect(t, x.UpdatedAt, x.CreatedAt)
}

func Test_Record_UpdatedAt(t *testing.T) {
	x := &Record{}
	err := Save(x, true)
	expect(t, err, nil)
	expect(t, x.UpdatedAt, x.CreatedAt)
	ua := x.UpdatedAt
	ca := x.CreatedAt
	err = Save(x, true)
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

func Test_Record_AssignID(t *testing.T) {
	x := &Record{}
	x.AssignID("CHEESE")
	expect(t, x.ID, "CHEESE")
}

func Test_Record_GetID(t *testing.T) {
	x := &Record{}
	x.ID = "CHEESE"
	expect(t, x.GetID(), "CHEESE")
}

func Test_Find_Fail(t *testing.T) {
	record := &Record{}
	err := Find(record, "666666666")
	refute(t, err, nil)
}

func Test_Find_Success(t *testing.T) {
	record := &Record{}
	err := Save(record, false)
	expect(t, err, nil)
	r := &Record{}
	err = Find(r, record.ID)
	expect(t, err, nil)
}

func Test_FindByIndex_Fail(t *testing.T) {
	record := &Record{}
	err := FindByIndex(record, "id", "666666666")
	refute(t, err, nil)
}

func Test_FindByIndex_Success(t *testing.T) {
	record := &Record{}
	err := Save(record, false)
	expect(t, err, nil)

	r := &Record{}
	err = FindByIndex(r, "id", record.ID)
	expect(t, err, nil)
	expect(t, r.ID, record.ID)
}

func Test_UpdateAttributes(t *testing.T) {
	record := &Record{}
	err := Save(record, false)
	expect(t, err, nil)

	err = UpdateAttributes(record, map[string]interface{}{
		"cheese": "queso",
	})
	expect(t, err, nil)
}

func Test_UpdateAttribute(t *testing.T) {
	record := &Record{}
	err := Save(record, false)
	expect(t, err, nil)

	err = UpdateAttribute(record, "cheese", "queso")
	expect(t, err, nil)
}
