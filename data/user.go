package data

import (
	"time"
)

const CanEdit = 1
const CanBroweAll = 2
const CanEditAll = 4
const CanAdmin = 8

type User struct {
	Id             int
	Name           string
	Password       string
	DepartmentId   int
	Privileges     int
	DepartmentName string
}
type Session struct {
	Id        int
	Uuid      string
	UserId    int
	CreatedAt time.Time
}

func UserDeleteAll() (err error) {
	_, err = Db.Exec("delete from users")
	return
}
func SessionDeleteAll() (err error) {
	_, err = Db.Exec("delete from sessions")
	return
}
func (user *User) Create() (err error) {
	_, err = Db.Exec("insert into users (name,password,department_id,privileges) values($1,$2,$3,$4)",
		user.Name, Encrypt(user.Password), user.DepartmentId, user.Privileges)
	if err != nil {
		return
	}

	err = Db.QueryRow(`select rowid,name,password,department_id,privileges from users where
		rowid = (select last_insert_rowid())`).Scan(&user.Id,
		&user.Name, &user.Password, &user.DepartmentId, &user.Privileges)

	return
}
func (user *User) Delete() (err error) {
	_, err = Db.Exec("delete from users where rowid=$1", user.Id)
	return
}
func (user *User) Update() (err error) {
	if len(user.Password) < 40 {
		user.Password = Encrypt(user.Password)
	}
	_, err = Db.Exec("update users set name=$1, password=$2,department_id=$3,privileges=$4 where rowid=$5",
		user.Name, user.Password, user.DepartmentId, user.Privileges, user.Id)
	return
}
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT rowid, name, password, department_id,privileges FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.DepartmentId, &user.Privileges); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}
func UserByID(id int) (user User, err error) {
	err = Db.QueryRow("select rowid, name, password, department_id,privileges from users where rowid=$1",
		id).Scan(&user.Id, &user.Name, &user.Password, &user.DepartmentId, &user.Privileges)
	return
}
func UserByName(name string) (user User, err error) {
	err = Db.QueryRow("select rowid,name,password,department_id,privileges from users where name=$1",
		name).Scan(&user.Id, &user.Name, &user.Password, &user.DepartmentId, &user.Privileges)
	return
}

func (user *User) CreateSession() (session Session, err error) {
	_, err = Db.Exec("insert into sessions(uuid,user_id,created_at) values($1,$2,$3)",
		createUUID(), user.Id, time.Now().Format(timeformat))
	if err != nil {
		return
	}
	s := ""
	err = Db.QueryRow(`select rowid,uuid,user_id,created_at from sessions where 
		rowid in (select last_insert_rowid())`).Scan(&session.Id,
		&session.Uuid, &session.UserId, &s)
	if err != nil {
		return
	}
	session.CreatedAt, err = time.Parse(timeformat, s)
	return
}
func (user *User) Session() (session Session, err error) {
	s := ""
	err = Db.QueryRow(`select rowid,uuid,user_id,created_at from sessions where 
		user_id = $1`, user.Id).Scan(&session.Id,
		&session.Uuid, &session.UserId, &s)
	if err != nil {
		return
	}
	session.CreatedAt, err = time.Parse(timeformat, s)
	return
}
func (session *Session) DeleteByUUID() (err error) {
	_, err = Db.Exec("delete from sessions where uuid=$1", session.Uuid)
	return
}
func (session *Session) Check() (valid bool, err error) {
	s := ""
	err = Db.QueryRow(`select rowid,uuid,user_id,created_at from sessions where 
		uuid = $1`, session.Uuid).Scan(&session.Id,
		&session.Uuid, &session.UserId, &s)
	if err != nil {
		return
	}
	session.CreatedAt, err = time.Parse(timeformat, s)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}
func (session *Session) User() (user User, err error) {
	err = Db.QueryRow(`select rowid,name,password,department_id,privileges from users  where 
		rowid = $1`, session.UserId).Scan(&user.Id, &user.Name, &user.Password, &user.DepartmentId, &user.Privileges)
	return
}
