package data

import (
	"database/sql"
	"testing"
)

func Test_DeparmentCreate(t *testing.T) {
	setup()
	if err := depts[0].Create(); err != nil {
		t.Error(err, "Cannot create department.")
	}
	if depts[0].Id == 0 {
		t.Errorf("No id department.")
	}
	dept, err := DepartmentByID(depts[0].Id)
	if err != nil {
		t.Error(err, "Department not created.")
	}
	if depts[0].Name != dept.Name {
		t.Errorf("Department retrieved is not the same as the one created.")
	}
}
func Test_DepartmentDelete(t *testing.T) {
	setup()
	if err := depts[0].Create(); err != nil {
		t.Error(err, "Cannot create department.")
	}
	if err := depts[0].Delete(); err != nil {
		t.Error(err, "Cannot delete departmnet")
	}
	_, err := DepartmentByID(depts[0].Id)
	if err != sql.ErrNoRows {
		t.Error(err, "Department not deleted.")
	}
}

func Test_DepartmentUpdate(t *testing.T) {
	setup()
	if err := depts[0].Create(); err != nil {
		t.Error(err, "Cannot create department.")
	}
	depts[0].Name = "Random Department"
	if err := depts[0].Update(); err != nil {
		t.Error(err, "Cannot update department")
	}
	dept, err := DepartmentByID(depts[0].Id)
	if err != nil {
		t.Error(err, "Cannot get department")
	}
	if dept.Name != "Random Department" {
		t.Error(err, "- Department not updated")
	}
}

func Test_Departments(t *testing.T) {
	setup()
	for _, dept := range depts {
		if err := dept.Create(); err != nil {
			t.Error(err, "Cannot create department.")
		}
	}
	ds, err := Departments()
	if err != nil {
		t.Error(err, "Cannot retrieve department.")
	}
	if len(ds) != 2 {
		t.Error(err, "Wrong number of department retrieved")
	}
	if ds[0].Name != depts[0].Name {
		t.Error(ds[0], depts[0], "Wrong department retrieved")
	}
}
