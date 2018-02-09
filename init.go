package gosql

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

func GetMySQLConn(username, password, url, dbname string) (sqlbuilder.Database, error) {
	var settings mysql.ConnectionURL
	settings, err := mysql.ParseURL(username + ":" + password + "@tcp(" + url + ")/" + dbname)

	if err != nil {
		return nil, err
	}

	var db sqlbuilder.Database
	db, err = mysql.Open(settings)

	return db, err
}
