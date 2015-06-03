package models

import (
	"errors"
	r "github.com/dancannon/gorethink"
	c "github.com/keighl/shrimp/config"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	Config *c.Configuration
	DB     *r.Session
	Logger = log.New(os.Stdout, "models: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	Env(os.Getenv("SHRIMP_ENV"))
}

func Env(env string) {
	Config = c.Config(env)
	var err error
	DB, err = r.Connect(r.ConnectOpts{
		Address:  Config.RethinkHost,
		Database: Config.RethinkDatabase,
	})
	if err != nil {
		panic(err)
	}
}

type Recorder interface {
	Table() string

	GetID() string
	AssignID(string)
	IsNewRecord() bool

	// Validation
	BeforeValidate()
	PerformValidations()
	AfterValidate()
	HasErrors() bool
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

func Save(x Recorder, validate bool) error {

	if validate {
		Validate(x)
		if x.HasErrors() {
			return errors.New("Validation errors")
		}
	}

	x.BeforeSave()

	if x.IsNewRecord() {
		x.BeforeCreate()
		res, err := r.Table(x.Table()).Insert(x).RunWrite(DB)
		if err != nil {
			return err
		}
		if len(res.GeneratedKeys) > 0 {
			x.AssignID(res.GeneratedKeys[0])
		}
		x.AfterCreate()
		x.AfterSave()
		return nil
	}

	x.BeforeUpdate()
	_, err := r.Table(x.Table()).Get(x.GetID()).Replace(x).RunWrite(DB)
	if err != nil {
		return err
	}
	x.AfterUpdate()
	x.AfterSave()
	return nil
}

func Delete(x Recorder) error {
	x.BeforeDelete()
	_, err := r.Table(x.Table()).Get(x.GetID()).Delete().RunWrite(DB)
	if err != nil {
		return err
	}
	x.AfterDelete()
	return nil
}

func Validate(x Recorder) {
	x.BeforeValidate()
	x.PerformValidations()
	x.AfterValidate()
}

//////////////////////////////
//////////////////////////////

type Record struct {
	ID        string          `gorethink:"id,omitempty" json:"id"`
	CreatedAt time.Time       `gorethink:"created_at" json:"created_at"`
	UpdatedAt time.Time       `gorethink:"updated_at" json:"-"`
	Errors    []string        `gorethink:"-" json:"errors,omitempty"`
	ErrorMap  map[string]bool `gorethink:"-" json:"-"`
}

func (x *Record) Table() string {
	return "records"
}

func (x *Record) IsNewRecord() bool {
	return x.GetID() == ""
}

//////////////////////////////
// ID ////////////////////////

func (x *Record) GetID() string {
	return x.ID
}

func (x *Record) AssignID(id string) {
	x.ID = id
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Record) BeforeValidate() {
	x.Errors = []string{}
	x.ErrorMap = map[string]bool{}
}

func (x *Record) AfterValidate() {}

func (x *Record) PerformValidations() {}

func (x *Record) HasErrors() bool {
	return len(x.Errors) > 0
}

func (x *Record) ErrorOn(attr string, message string) {
	x.ErrorMap[attr] = true
	x.Errors = append(x.Errors, message)
}

//////////////////////////////
// CALLBACKS /////////////////

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

//////////////////////////////
// GETTERS ///////////////////

func Find(result Recorder, id string) error {
	res, err := r.Table(result.Table()).Get(id).Run(DB)
	if err != nil {
		return err
	}
	err = res.One(result)
	return err
}

func FindByIndex(result Recorder, index string, value string) error {
	res, err := r.Table(result.Table()).GetAllByIndex(index, value).Run(DB)
	if err != nil {
		return err
	}
	err = res.One(result)
	return err
}

//////////////////////////////
// SETTERS ///////////////////

func UpdateAttribute(result Recorder, attribute string, value interface{}) error {
	payload := map[string]interface{}{attribute: value}
	return UpdateAttributes(result, payload)
}

func UpdateAttributes(result Recorder, values map[string]interface{}) error {
	_, err := r.Table(result.Table()).Get(result.GetID()).Update(values).RunWrite(DB)
	return err
}
