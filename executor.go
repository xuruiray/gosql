package gosql

import (
	"database/sql"
	"strings"
	upperDB "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// Execute 封装sql方法 执行sql语句 返回受影响的行数
func Execute(conn sqlbuilder.Database, sql string, params map[string]interface{}) (int64, error) {
	presql, args, err := GetPreparedStatement(sql, params)
	if err != nil {
		return 0, err
	}

	affected, err := conn.Exec(presql, args...)
	if err != nil {
		return 0, err
	}
	return affected.RowsAffected()
}

// GetList 封装查询方法 查询多行数据 结果封入res中
func QueryList(conn sqlbuilder.Database, sql string, params map[string]interface{}, res interface{}) error {
	rows, err := prepareParams(conn, sql, params)
	if err != nil {
		return err
	}

	if rows != nil {
		iterator := sqlbuilder.NewIterator(rows)
		err = iterator.All(res)
	}

	if err == upperDB.ErrNoMoreRows {
		return nil
	}

	return err
}

// GetOne 封装查询方法 查询单行数据 结果封入res中
func QueryOne(conn sqlbuilder.Database, sqlStr string, params map[string]interface{}, res interface{}) error {
	rows, err := prepareParams(conn, sqlStr, params)
	if err != nil {
		return err
	}

	if rows != nil {
		iterator := sqlbuilder.NewIterator(rows)
		err = iterator.One(res)
	}

	if err == upperDB.ErrNoMoreRows {
		return nil
	}

	return err
}

// IsDuplicatedMySQLError 判断是否为重复插入错误
func IsDuplicatedError(err error) bool {
	return strings.Contains(err.Error(), "Error 1062: Duplicate")
}

func prepareParams(conn sqlbuilder.Database, sql string, params map[string]interface{}) (*sql.Rows, error) {
	presql, args, err := GetPreparedStatement(sql, params)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(presql, args...)
	return rows, err
}
