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
type WorkItem struct {
	Id       int
	OrderId  int
	Work     string
	Place    string
	Unit     int
	Quantity float32
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
	isubmit, ilocked, s := 0, 0, ""
	err = Db.QueryRow(`select rowid, 
		department_id,date_begin,date_end,provider_id,cartype_id,carnum,usefor,
		submit,locked,careat_at from orders where rowid =$1`, id,
	).Scan(&order.Id, &order.DepartmentId, &order.DateBegin, &order.DateEnd,
		&order.ProviderId, &order.CarTypeId, &order.CarNum, &order.UseFor,
		&isubmit, &ilocked, &s)
	if err != nil {
		return
	}
	order.Submit = (isubmit == 1)
	order.Locked = (ilocked == 1)
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

	return
}
func (order *Order) Update() (err error) {
	boolmap := map[bool]int{false: 0, true: 1}
	_, err = Db.Exec(`update orders
		set department_id=$1,
		date_begin=$2,
		date_end=$3,
		provider_id=$4,
		cartype_id=$5,
		usefor=$6,
		carnum=$7,
		submit=$8,
		locked=$9,
		careat_at=$10 where rowid=$11`,
		order.DepartmentId, order.DateBegin, order.DateEnd, order.ProviderId,
		order.CarTypeId, order.UseFor, order.CarNum, boolmap[order.Submit],
		boolmap[order.Locked], order.CreatedAt.Format(timeformat), order.Id)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (order *Order) Delete() (err error) {
	_, err = Db.Exec("delete from orders where rowid=$1", order.Id)
	return
}
func (order *Order) CreateWorkItem(work string, place string, unit int,
	quantity float32) (wi WorkItem, err error) {
	fmt.Println(order.Id, work, unit, quantity)
	_, err = Db.Exec(`insert into workitems(order_id,work,place,unit,quantity) 
	values($1,$2,$3,$4,$5)`, order.Id, work, place, unit, quantity)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = Db.QueryRow(`select rowid,order_id,work,place,unit,quantity
		from workitems where rowid in (select last_insert_rowid())`,
	).Scan(&wi.Id, &wi.OrderId, &wi.Work, &wi.Place, &wi.Unit, &wi.Quantity)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (order *Order) WorkItems() (workitems []WorkItem, err error) {
	rows, err := Db.Query(`select rowid,order_id,work,place,unit,quantity
		from workitems where order_id=$1`, order.Id)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		wi := WorkItem{}
		rows.Scan(&wi.Id, &wi.OrderId, &wi.Work, &wi.Place, &wi.Unit,
			&wi.Quantity)
		workitems = append(workitems, wi)
	}
	return
}
func (workitem *WorkItem) Delete() (err error) {
	_, err = Db.Exec("delete from workitems where rowid=$1", workitem.Id)
	return
}
func (workitem *WorkItem) Update() (err error) {
	_, err = Db.Exec(`update workitems set work=$1,place=$2,unit=$3,
	quantity=$4 where rowid=$5`, &workitem.Work, &workitem.Place, &workitem.Unit,
		&workitem.Quantity, &workitem.Id)
	return
}
func WorkItemByID(id int) (workitem WorkItem, err error) {
	err = Db.QueryRow(`select rowid,order_id,work,place,unit,quantity
		from workitems where rowid=$1`, id).Scan(&workitem.Id, &workitem.OrderId,
		&workitem.Work, &workitem.Place, &workitem.Unit, &workitem.Quantity)
	if err != nil {
		fmt.Println(err)
	}
	return
}
