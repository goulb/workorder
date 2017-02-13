package data

type CarType struct {
	Id       int
	Weight   int
	TypeName string
}

func (cartype *CarType) Create() (err error) {
	_, err = Db.Exec("insert into cartypes(weight,type_name) values($1,$2)",
		cartype.Weight, cartype.TypeName)
	if err != nil {
		return
	}

	err = Db.QueryRow(`select rowid,weight,type_name from cartypes where 
		rowid in (select last_insert_rowid())`).Scan(
		&cartype.Id, &cartype.Weight, &cartype.TypeName)

	return
}

func CarTypes() (cartypes []CarType, err error) {
	rows, err := Db.Query(`select rowid,weight,type_name 
		from cartypes order by rowid`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		cartype := CarType{}
		if err = rows.Scan(
			&cartype.Id, &cartype.Weight, &cartype.TypeName); err != nil {
			return
		}
		cartypes = append(cartypes, cartype)
	}
	return
}
