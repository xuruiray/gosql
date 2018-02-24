package gosql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	username = "root"
	password = "1324"
	url      = "0.0.0.0:3307"
	dbname   = "test"
)

func TestGetMySQLConn(t *testing.T) {
	conn, err := GetMySQLConn(username, password, url, dbname)

	assert.Equal(t, nil, err, "error:%v", err)
	assert.NotEqual(t, nil, conn)
}
