package gosql

import "testing"

const (
	username = "root"
	password = "1324"
	url      = "0.0.0.0:3307"
	dbname   = "test"
)

func TestGetMySQLConn(t *testing.T) {
	conn, err := GetMySQLConn(username, password, url, dbname)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if conn == nil {
		t.Error("connect is nil")
		t.Fail()
	}

	t.Log("test GetMySQLConn finish")
}
