package gosql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
		caseName string
		datamap  map[string]interface{}
		want     DriverInfo
		wantErr  error
		sql      string
	}{
		{
			caseName: "正常流程",
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
		assert.Equal(t, v.wantErr, err, "caseName: %v", v.caseName)
		assert.Equal(t, v.want, result, "caseName: %v", v.caseName)
	}
}

// TestQueryList 测试 QueryList
func TestQueryList(t *testing.T) {

	var tests = []struct {
		caseName string
		datamap  map[string]interface{}
		want     []DriverInfo
		wantErr  error
		sql      string
	}{
		{
			caseName: "正常流程",
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
		assert.Equal(t, v.wantErr, err, "caseName: %v", v.caseName)
		assert.Equal(t, v.want, resultList, "caseName: %v", v.caseName)
	}
}

// TestExecute 测试 Execute
func TestExecute(t *testing.T) {

	var tests = []struct {
		caseName string
		datamap  map[string]interface{}
		want     int64
		wantErr  error
		sql      string
	}{
		{
			caseName: "正常流程",
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
		assert.Equal(t, v.wantErr, err, "caseName: %v", v.caseName)
		assert.Equal(t, v.want, affected, "caseName: %v", v.caseName)
	}
}

// TestIsDuplicatedError 测试 IsDuplicatedError
func TestIsDuplicatedError(t *testing.T) {

	var tests = []struct {
		caseName   string
		datamap    map[string]interface{}
		want       bool
		wantAffect int64
		sql        string
	}{
		{
			caseName: "正常流程",
			datamap: map[string]interface{}{
				"id":        123,
				"driver_id": 123,
				"name":      "xurui",
				"age":       12,
			},
			want:       true,
			wantAffect: 0,
			sql:        "insert into driver_info(id, driver_id, name, age) VALUES ($id, $driver_id, $name, $age)",
		},
	}

	for _, v := range tests {
		affected, err := Execute(conn, v.sql, v.datamap)
		isDup := IsDuplicatedError(err)
		assert.Equal(t, v.wantAffect, affected, "caseName: %v", v.caseName)
		assert.Equal(t, v.want, isDup, "caseName: %v", v.caseName)
	}
}
