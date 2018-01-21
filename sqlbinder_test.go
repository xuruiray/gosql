package gosql

import (
	"reflect"
	"testing"
)

func Test_getStatement(t *testing.T) {
	type args struct {
		sqlStr  string
		datamap map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    string
			want1   map[string]interface{}
			wantErr bool
		}{name: "test01", args: args{
			sqlStr: "select * from driver_info where id=#id,driver_id=#driver_id #sort",
			datamap: map[string]interface{}{
				"id":        123,
				"driver_id": "456",
				"sort":      "order by id",
			},
		}, want: "select * from driver_info where id=123,driver_id=456 order by id ", want1: map[string]interface{}{}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getStatement([]byte(tt.args.sqlStr), tt.args.datamap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStatement() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getStatement() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetPreparedStatement(t *testing.T) {
	type args struct {
		sql    string
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    string
			want1   []interface{}
			wantErr bool
		}{name: "test01", args: args{
			sql: "select * from driver_info where id=#id,driver_id=$driver_id #sort",
			params: map[string]interface{}{
				"id":        123,
				"driver_id": "456",
				"sort":      "order by id",
			},
		}, want: "select * from driver_info where id=123,driver_id=? order by id ",want1:[]interface{}{"456"}, wantErr: false},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetPreparedStatement(tt.args.sql, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPreparedStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPreparedStatement() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPreparedStatement() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
