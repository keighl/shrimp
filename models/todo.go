package models

import ()

type Todo struct {
	Record
	UserID string `gorethink:"user_id" json:"-"`
	Title  string `gorethink:"title" json:"title"`
}

type TodoAttrs struct {
	Title string `form:"title" json:"title"`
}

func (x *Todo) Table() string {
	return "todos"
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Todo) PerformValidations() {
	x.Record.PerformValidations()
	x.ValidateTitle()
}

func (x *Todo) ValidateTitle() {
	if x.Title == "" {
		x.ErrorOn("Title", "Title can't be blank.")
	}
}

//////////////////////////////
// OTHER /////////////////////

func (x *Todo) UpdateFromAttrs(attrs TodoAttrs) {
	if attrs.Title != "" {
		x.Title = attrs.Title
	}
}

func (x *TodoAttrs) Todo() *Todo {
	return &Todo{
		Title: x.Title,
	}
}
