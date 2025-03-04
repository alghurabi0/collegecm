package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Subjects   SubjectModel
	Students   StudentModel
	Carryovers CarryoverModel
	Exempteds  ExemptedModel
	Marks      MarkModel
	Customs    CustomModel
	Years      YearModel
	Users      UserModel
	Privileges PrivilegeModel
	Tables     TableModel
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Subjects:   SubjectModel{DB: db},
		Students:   StudentModel{DB: db},
		Carryovers: CarryoverModel{DB: db},
		Exempteds:  ExemptedModel{DB: db},
		Marks:      MarkModel{DB: db},
		Customs:    CustomModel{DB: db},
		Years:      YearModel{DB: db},
		Users:      UserModel{DB: db},
		Privileges: PrivilegeModel{DB: db},
		Tables:     TableModel{DB: db},
	}
}
