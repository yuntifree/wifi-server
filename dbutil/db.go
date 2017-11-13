package dbutil

import "database/sql"

const (
	sqlDsn = "access:^yunti9df3b01c$@tcp(rm-wz9sb2613092ki9xn.mysql.rds.aliyuncs.com:3306)/wifi?charset=utf8"
)

//NewDB return db connection
func NewDB() (*sql.DB, error) {
	return newDB(sqlDsn)
}

func newDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
