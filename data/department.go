package data

import (
	//"database/sql"
	"log"
)

type Department struct {
	Id   int
	Name string
}

func (dept *Department) Create() (err error) {

	stmt, err := Db.Prepare("insert into departments values(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(dept.Name)
	if err != nil {
		log.Fatal(err)
	}
	err = Db.QueryRow(`select rowid,name from departments where 
		rowid in (select max(rowid) from departments LIMIT 1)`).Scan(&dept.Id, &dept.Name)
	if err != nil {
		log.Fatal(err)
	}

	return
}
func Departments() (depts []Department, err error) {
	rows, err := Db.Query("select rowid,name from departments order by rowid")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		dept := Department{}
		if err = rows.Scan(&dept.Id, &dept.Name); err != nil {
			return
		}
		depts = append(depts, dept)
	}
	return
}
func (dept Department) Delete() (err error) {

	_, err = Db.Exec("delete from departments where rowid = $1", dept.Id)
	return
}
func (dept Department) Update() (err error) {
	_, err = Db.Exec("update departments set name = $1 where rowid = $2", dept.Name, dept.Id)
	return
}
func DepartmentByID(id int) (dept Department, err error) {
	dept = Department{}
	err = Db.QueryRow(`select rowid,name from departments where 
		rowid = $1 `, id).Scan(&dept.Id, &dept.Name)

	return
}

func DepartmentDeleteAll() (err error) {
	statement := "delete from departments"
	_, err = Db.Exec(statement)
	return
}
