package gosql

import (
	"fmt"
	"testing"
)

const (
	username = "root"
	password = "1324"
	url      = "0.0.0.0:3307"
	dbname   = "test"
)

type DriverInfo struct {
	ID       int    `db:"id"`
	DriverID int    `db:"driver_id"`
	Name     string `db:"name"`
	Age      int    `db:"age"`
}

func TestExecute(t *testing.T) {
	conn, err := GetMySQLConn(username, password, url, dbname)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	affected, err := Execute(conn, "insert into driver_info(driver_id, name, age) VALUES (123,'xurui',12)", nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if affected == -1 {
		t.Fail()
	}
}

func TestGetOne(t *testing.T) {
	conn, err := GetMySQLConn(username, password, url, dbname)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	datamap := map[string]interface{}{
		"tablename": "driver_info",
		"driver_id": 123,
		"sort":      " ",
	}

	var result DriverInfo
	err = GetOne(conn, "select * from #tablename where driver_id = $driver_id #sort", datamap, &result)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	fmt.Println(result)
}
