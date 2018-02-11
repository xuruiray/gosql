package gosql

import (
	"fmt"
	"reflect"
	"testing"
	"upper.io/db.v3/lib/sqlbuilder"
)

type DriverInfo struct {
	ID       int    `db:"id"`
	DriverID int    `db:"driver_id"`
	Name     string `db:"name"`
	Age      int    `db:"age"`
}

var conn sqlbuilder.Database

// 初始化 mysql 连接
func init() {
	var err error
	conn, err = GetMySQLConn(username, password, url, dbname)
	if err != nil {
		fmt.Println("mysql connect init failed")
	}
}

// TestQueryOne 测试 QueryOne
func TestQueryOne(t *testing.T) {

	var tests = []struct {
		datamap map[string]interface{}
		want    DriverInfo
		wantErr error
		sql     string
	}{
		{
			datamap: map[string]interface{}{
				"tablename": "driver_info",
				"id":        123,
				"driver_id": 456,
				"name":      "xurui",
				"age":       12,
				"sort":      " ",
			},
			want: DriverInfo{
				ID:       123,
				DriverID: 456,
				Name:     "xurui",
				Age:      12,
			},
			wantErr: nil,
			sql:     "select * from #tablename where id=$id and driver_id=$driver_id and name=$name and age=$age #sort",
		},
	}

	for _, v := range tests {
		var result DriverInfo
		err := QueryOne(conn, v.sql, v.datamap, &result)
		if err != v.wantErr {
			t.Errorf("QueryOne||sql=%v", v.sql)
			t.Errorf("QueryOne||want error=%v||get error=%v", v.wantErr, err)
		}
		if !reflect.DeepEqual(result, v.want) {
			t.Errorf("QueryOne||sql=%v", v.sql)
			t.Errorf("QueryOne||want=%v||get=%v", v.want, result)
		}
	}

	t.Log("test QueryOne finish")

}

// TestQueryList 测试 QueryList
func TestQueryList(t *testing.T) {

	var tests = []struct {
		datamap map[string]interface{}
		want    []DriverInfo
		wantErr error
		sql     string
	}{
		{
			datamap: map[string]interface{}{
				"tablename": "driver_info",
				"limit":     "4",
			},
			want: []DriverInfo{
				{123, 456, "xurui", 12},
				{124, 457, "xurui", 13},
				{125, 458, "xurui", 14},
				{126, 459, "xurui", 15},
			},
			wantErr: nil,
			sql:     "select * from #tablename limit #limit",
		},
	}

	for _, v := range tests {
		var resultList []DriverInfo
		err := QueryList(conn, v.sql, v.datamap, &resultList)
		if err != v.wantErr {
			t.Errorf("QueryOne||sql=%v", v.sql)
			t.Errorf("QueryOne||want error=%v||get error=%v", v.wantErr, err)
		}

		if !reflect.DeepEqual(v.want, resultList) {
			t.Errorf("QueryList||sql=%v", v.sql)
			t.Errorf("QueryList||want=%v||get=%v", v.want, resultList)
		}
	}

	t.Log("test QueryList finish")
}

// TestExecute 测试 Execute
func TestExecute(t *testing.T) {

	var tests = []struct {
		datamap map[string]interface{}
		want    int64
		wantErr error
		sql     string
	}{
		{
			datamap: map[string]interface{}{
				"driver_id": 13579,
				"name":      "xurui",
				"age":       12,
			},
			want:    1,
			wantErr: nil,
			sql:     "insert into driver_info(driver_id, name, age) VALUES ($driver_id,$name,$age)",
		},
	}

	for _, v := range tests {
		affected, err := Execute(conn, v.sql, v.datamap)
		if err != v.wantErr {
			t.Errorf("Execute||sql=%v", v.sql)
			t.Errorf("Execute||want error=%v||get error=%v", v.wantErr, err)
		}
		if affected != v.want {
			t.Errorf("Execute||sql=%v", v.sql)
			t.Errorf("Execute||want=%v||get=%v", v.want, affected)
		}
	}

	t.Log("test Execute finish")
}

// TestIsDuplicatedError 测试 IsDuplicatedError
func TestIsDuplicatedError(t *testing.T) {

	var tests = []struct {
		datamap map[string]interface{}
		want    int64
		wantErr error
		sql     string
	}{
		{
			datamap: map[string]interface{}{
				"id":        123,
				"driver_id": 123,
				"name":      "xurui",
				"age":       12,
			},
			want: 0,
			sql:  "insert into driver_info(id, driver_id, name, age) VALUES ($id, $driver_id, $name, $age)",
		},
	}

	for _, v := range tests {
		affected, err := Execute(conn, v.sql, v.datamap)
		if !IsDuplicatedError(err) {
			t.Errorf("IsDuplicatedError||sql=%v", v.sql)
			t.Errorf("IsDuplicatedError||want error=%v||get error=%v", v.wantErr, err)
		}
		if affected != v.want {
			t.Errorf("IsDuplicatedError||sql=%v", v.sql)
			t.Errorf("IsDuplicatedError||want=%v||get=%v", v.want, affected)
		}
	}

	t.Log("test IsDuplicatedError finish")
}
