package gosql

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPreparedStatement(t *testing.T) {
	var tests = []struct {
		caseName string
		sql      string
		datamap  map[string]interface{}

		wantSql  string
		wantList []interface{}
		wantErr  error
	}{
		{
			caseName: "正常流程",
			sql:      "select #selectElement from #tablename where name=$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"name":          "xurui",
				"idset":         []int{123, 124, 125},
			},
			wantSql:  "select id, driver_id, name from driver_info where name=? and id in ?,?,?",
			wantList: []interface{}{"xurui", 123, 124, 125},
			wantErr:  nil,
		}, {
			caseName: "末尾参数 '$'",
			sql:      "select #selectElement from #tablename where name=$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"name":          "xurui",
				"idset":         []int{123, 124, 125},
			},
			wantSql:  "select id, driver_id, name from driver_info where name=? and id in ?,?,?",
			wantList: []interface{}{"xurui", 123, 124, 125},
			wantErr:  nil,
		}, {
			caseName: "末尾参数 '#'",
			sql:      "select #selectElement from #tablename where name=$name limit #pageSize",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"name":          "xurui",
				"pageSize":      2,
			},
			wantSql:  "select id, driver_id, name from driver_info where name=? limit 2",
			wantList: []interface{}{"xurui"},
			wantErr:  nil,
		}, {
			caseName: "无输入参数",
			sql:      "select * from tablename where name=xurui",
			datamap:  map[string]interface{}{},
			wantSql:  "select * from tablename where name=xurui",
			wantList: []interface{}{},
			wantErr:  nil,
		}, {
			caseName: "'$' 转义测试",
			sql:      "select #selectElement from #tablename where name=$$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"idset":         []int{123, 124, 125},
			},
			wantSql:  "select id, driver_id, name from driver_info where name=$name and id in ?,?,?",
			wantList: []interface{}{123, 124, 125},
			wantErr:  nil,
		}, {
			caseName: "'#' 转义测试",
			sql:      "select ##selectElement from #tablename where name=$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename": "driver_info",
				"name":      "xurui",
				"idset":     []int{123, 124, 125},
			},
			wantSql:  "select #selectElement from driver_info where name=? and id in ?,?,?",
			wantList: []interface{}{"xurui", 123, 124, 125},
			wantErr:  nil,
		}, {
			caseName: "'#' '$' 转义测试",
			sql:      "select ##selectElement from #tablename where name=$$name",
			datamap: map[string]interface{}{
				"tablename": "driver_info",
			},
			wantSql:  "select #selectElement from driver_info where name=$name",
			wantList: []interface{}{},
			wantErr:  nil,
		}, {
			caseName: "参数缺省测试 '$'参数",
			sql:      "select #selectElement from #tablename where name=$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"name":          "xurui",
			},
			wantSql:  "",
			wantList: []interface{}(nil),
			wantErr:  errors.New("lost params idset"),
		}, {
			caseName: "参数缺省测试 '#'参数",
			sql:      "select #selectElement from #tablename where name=$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename": "driver_info",
				"name":      "xurui",
				"idset":     []int{123, 124, 125},
			},
			wantSql:  "",
			wantList: []interface{}(nil),
			wantErr:  errors.New("lost params selectElement"),
		}, {
			caseName: "非法属性名 包含#",
			sql:      "select #select#Element from #tablename where name=$name and id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"name":          "xurui",
				"idset":         []int{123, 124, 125},
			},
			wantSql:  "",
			wantList: []interface{}(nil),
			wantErr:  errors.New("can`t insert # in the params"),
		}, {
			caseName: "非法属性名 包含$",
			sql:      "select #selectElement from #tablename where name=$na$me and id in $idset",
			datamap: map[string]interface{}{
				"tablename":     "driver_info",
				"selectElement": "id, driver_id, name",
				"name":          "xurui",
				"idset":         []int{123, 124, 125},
			},
			wantSql:  "",
			wantList: []interface{}(nil),
			wantErr:  errors.New("can`t insert $ in the params"),
		},
	}

	for _, v := range tests {
		resultSql, resultList, err := GetPreparedStatement(v.sql, v.datamap)
		assert.Equal(t, v.wantErr, err, "caseName:%v", v.caseName)
		assert.Equal(t, v.wantSql, resultSql, "caseName:%v", v.caseName)
		assert.Equal(t, v.wantList, resultList, "caseName:%v", v.caseName)
	}
}
