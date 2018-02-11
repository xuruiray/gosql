package gosql

import (
	"reflect"
	"testing"
)

func TestGetPreparedStatement(t *testing.T) {
	var tests = []struct {
		sql     string
		datamap map[string]interface{}

		wantSql  string
		wantList []interface{}
		wantErr  error
	}{
		{
			sql: "select #selectElement from #tablename where id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"idset":         []int{123, 124, 125},
			},
			wantSql:  "select id, driver_id, name from driver_info where id in ?,?,?",
			wantList: []interface{}{123, 124, 125},
			wantErr:  nil,
		},
		{
			sql: "select #selectElement from #tablename where id=$id limit #pageSize",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"id":            123,
				"pageSize":      4,
			},
			wantSql:  "select id, driver_id, name from driver_info where id=? limit 4",
			wantList: []interface{}{123},
			wantErr:  nil,
		},
	}

	for _, v := range tests {
		getSql, getList, getErr := GetPreparedStatement(v.sql, v.datamap)
		if getSql != v.wantSql {
			t.Errorf("GetPreparedStatement||want presql=%v||get presql=%v", v.wantSql, getSql)
		}
		if !reflect.DeepEqual(getList, v.wantList) {
			t.Errorf("GetPreparedStatement||want params list=%v||get params list=%v", v.wantList, getList)
		}
		if !reflect.DeepEqual(getErr, v.wantErr) {
			t.Errorf("GetPreparedStatement||want error=%v||get error=%v", v.wantErr, getErr)
		}
	}

	t.Log("test GetPreparedStatement finish")
}
