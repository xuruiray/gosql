package gosql

import "context"

func Do(ctx context.Context, dbname string, tablename string, sqlStr string, params map[string]interface{}, res interface{}) error {
	return nil
}
