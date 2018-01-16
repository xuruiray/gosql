package gosql

import (
	"testing"
)

func TestGetParamsName(t *testing.T) {
	sqls := []string{
		"select event_id,order_id from %s where driver_id in (:ids) or driver_id=:driver_id",

		"select event_id,order_id from %s where driver_id in (:ids) or" +
			" driver_id=:driver_id and event_id>:Ece limit :pageStart,:PageCount order by :pagesize desc",

		"select event_id,order_id from %s where driver_id in (:id:s)",
	}

	resultMap := map[string][]string{
		"select event_id,order_id from %s where driver_id in (?) or driver_id=?": {"ids", "driver_id"},

		"select event_id,order_id from %s where driver_id in (?) or driver_id=? and event_id>? limit ?,? order by ? desc": {"ids",
			"driver_id", "Ece", "pageStart", "PageCount", "pagesize"},

		"select event_id,order_id from %s where driver_id in (?)": {"ids"},
	}

	output := []bool{
		true,
		true,
		false,
	}

	for index, sql := range sqls {
		presql, result, err := getParamsName(sql)
		if err != nil {
			if output[index] {
				t.Log(err)
				t.Fail()
			} else {
				continue
			}
		}

		if params, ok := resultMap[presql]; ok {
			for i, param := range params {
				if param != result[i] {
					if output[index] {
						t.Log(param, result[i])
						t.Fail()
					} else {
						continue
					}
				}
			}
		} else {
			if output[index] {
				t.Log("can find pre SQL")
				t.Fail()
			} else {
				continue
			}
		}

		if !output[index] {
			t.Logf("sql index %d should be faild", index)
			t.Fail()
		}

	}
}
