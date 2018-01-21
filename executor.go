package gosql

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3"
	"database/sql"
	"context"
)

// GetList 封装查询方法 查询多行数据 结果封入res中
func GetList(ctx context.Context, dbname string, sql string, params map[string]interface{}, res interface{}) error {

	rows, err := prepareParams(ctx, dbname, sql, params)
	if err != nil {
		return err
	}

	if rows != nil {
		iterator := sqlbuilder.NewIterator(rows)
		err = iterator.All(res)
	}

	return err
}

// GetOne 封装查询方法 查询单行数据 结果封入res中
func GetOne(ctx context.Context, dbname string, sqlStr string, params map[string]interface{}, res interface{}) error {

	rows, err := prepareParams(ctx, dbname, sqlStr, params)
	if err != nil {
		return err
	}

	if rows != nil {
		iterator := sqlbuilder.NewIterator(rows)
		// 实际上和 GetList 就这一行区别
		err = iterator.One(res)
	}

	return err
}

// Execute 封装sql方法 执行sql语句 返回受影响的行数
func Execute(ctx context.Context, dbname string, tablename string, sql string, params map[string]interface{}) (int64, error) {

	sql, args, err := GetPreparedStatement(sql, params)
	if err != nil {
		return 0, err
	}

	affected, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return affected.RowsAffected()
}

// 查询共有逻辑抽离
func prepareParams(ctx context.Context, dbname string, sql string, params map[string]interface{}) (*sql.Rows, error) {

	sql, args, err := GetPreparedStatement(sql, params)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)
	return rows, err
}
