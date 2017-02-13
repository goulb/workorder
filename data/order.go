package data

import (
	"fmt"
	"time"
)

type Order struct {
	Id           int
	DepartmentId int
	DateBegin    string
	DateEnd      string
	ProviderId   int
	CarTypeId    int
	CarNum       string
	UseFor       int
	Submit       bool
	Locked       bool
	CreatedAt    time.Time
}

func (order *Order) Create() (err error) {
	_, err = Db.Exec(`insert into orders(
		department_id,date_begin,date_end,provider_id,cartype_id,usefor,carnum,
		submit,locked,careat_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		order.DepartmentId, order.DateBegin, order.DateEnd, order.ProviderId,
		order.CarTypeId, order.UseFor, order.CarNum, 0, 0,
		time.Now().Format(timeformat))
	if err != nil {
		return
	}
	s := ""
	err = Db.QueryRow(`select rowid, 
		department_id,date_begin,date_end,provider_id,cartype_id,carnum,usefor,
		careat_at from orders where rowid in (select last_insert_rowid())`,
	).Scan(&order.Id, &order.DepartmentId, &order.DateBegin, &order.DateEnd,
		&order.ProviderId, &order.CarTypeId, &order.CarNum, &order.UseFor, &s)
	if err != nil {
		return
	}
	order.CreatedAt, err = time.Parse(timeformat, s)
	return
}
func OrderByID(id int) (order Order, err error) {
	s := ""
	err = Db.QueryRow(`select rowid, 
		department_id,date_begin,date_end,provider_id,cartype_id,carnum,usefor,
		careat_at from orders where rowid =$1`, id,
	).Scan(&order.Id, &order.DepartmentId, &order.DateBegin, &order.DateEnd,
		&order.ProviderId, &order.CarTypeId, &order.CarNum, &order.UseFor, &s)
	if err != nil {
		return
	}
	order.CreatedAt, err = time.Parse(timeformat, s)
	return
}
func Orders() (orders []Order, err error) {

	rows, err := Db.Query(`select 
		rowid, department_id,date_begin,date_end,provider_id,cartype_id,carnum,
		usefor,submit,locked,careat_at from orders order by rowid desc`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		isubmit, ilocked, s := 0, 0, ""
		order := Order{}
		if err = rows.Scan(&order.Id, &order.DepartmentId, &order.DateBegin,
			&order.DateEnd, &order.ProviderId, &order.CarTypeId, &order.CarNum,
			&order.UseFor, &isubmit, &ilocked, &s); err != nil {
			return
		}
		order.Submit = (isubmit == 1)
		order.Locked = (ilocked == 1)
		order.CreatedAt, err = time.Parse(timeformat, s)
		orders = append(orders, order)
	}
	return
}
func OrdersByDepartmentID(id int) (orders []Order, err error) {

	rows, err := Db.Query(`select 
		rowid, department_id,date_begin,date_end,provider_id,cartype_id,carnum,
		usefor,submit,locked,careat_at from orders where department_id=$1 
		order by rowid desc`, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		isubmit, ilocked, s := 0, 0, ""
		order := Order{}
		if err = rows.Scan(&order.Id, &order.DepartmentId, &order.DateBegin,
			&order.DateEnd, &order.ProviderId, &order.CarTypeId, &order.CarNum,
			&order.UseFor, &isubmit, &ilocked, &s); err != nil {
			fmt.Println(err)
			return
		}
		order.Submit = (isubmit == 1)
		order.Locked = (ilocked == 1)
		order.CreatedAt, err = time.Parse(timeformat, s)
		orders = append(orders, order)
	}
	fmt.Println(id, orders)
	return
}