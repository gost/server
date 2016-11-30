package postgis

import (
	"database/sql"
)

func ExecuteSelect(db *sql.DB, sql string) (interface{}, error) {
	/*rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var params []interface{}

		columns, _ := rows.Columns()
		for _, c := range columns {
			fmt.Println(c)
		}

		err := rows.Scan(params...)
		if err != nil {
			return nil, 0, err
		}

		return nil, nil
	}*/

	return nil, nil
}
