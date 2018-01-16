package gosql

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

//Named 将传入的sql与params参数tablename表名 整理成 prepare statement
func Named(sql string, params map[string]interface{}, tablename string) (string, []interface{}, error) {
	sql = fmt.Sprintf(sql, tablename)
	sql, names, err := getParamsName(sql)
	if err != nil {
		return "", nil, err
	}
	return bindMap(sql, names, params)
}

func bindMap(sql string, names []string, params map[string]interface{}) (string, []interface{}, error) {

	paramsArr := make([]interface{}, 0, len(params)*2)
	var err error
	for i, name := range names {
		v, ok := params[name]
		if !ok {
			return "", nil, errors.New("can`t match param " + name)
		}

		if !isSlice(v) {
			paramsArr = append(paramsArr, v)
			continue
		}

		//若属性名对应值为 切片 或 数组 ，则将数组拆分，并将 prestatement 中对应 占位符 个数增加
		t := reflect.ValueOf(v)
		sql, err = expandPlaceholder(sql, i, t.Len())
		if err != nil {
			return sql, nil, err
		}
		for i := 0; i < t.Len(); i++ {
			paramsArr = append(paramsArr, t.Index(i).Interface())
		}

	}

	return sql, paramsArr, nil
}

//getParamsName 转换 sql 字符串为 prestatement 并将 属性名 提取出来
func getParamsName(sql string) (string, []string, error) {
	names := make([]string, 0, 10)

	begin := 0
	found := false
	sqlResult := make([]rune, 0, len(sql))

	for i, b := range sql {
		if b == ':' {
			//若连续两个：：在转换 prestatement 时 会转义为 ：
			if sql[i-1] == ':' {
				sqlResult = append(sqlResult, ':')
				found = false
				continue
			}
			//若在属性名中发现 ： 则报错
			if found {
				return "", nil, errors.New("can`t insert : in the params")
			}

			//找到属性名的 起始下标
			sqlResult = append(sqlResult, '?')
			found = true
			begin = i + 1

			continue
		}

		if !found {
			sqlResult = append(sqlResult, b)
		}

		//sql 字符串中 属性名 必须为仅包含  '.' '_' 'A-Z' 'a-z' '0-9'
		if found && b != '.' && b != '_' && !unicode.IsLetter(b) && !unicode.IsNumber(b) {
			sqlResult = append(sqlResult, b)
			names = append(names, sql[begin:i])
			found = false
		}

	}
	if found {
		names = append(names, sql[begin:])
	}
	return string(sqlResult), names, nil
}

//expandPlaceholder 替换 sql 中 第 index 个 ‘？’ 为 count 个 ‘？’
func expandPlaceholder(sql string, index int, count int) (string, error) {
	buffer := make([]string, 0, count)
	for count > 0 {
		buffer = append(buffer, "?")
		count--
	}

	replace := strings.Join(buffer, ",")

	time := 0
	for i, b := range sql {
		if b != '?' {
			continue
		}

		time++

		if time != index+1 {
			continue
		}

		buffer := bytes.NewBufferString(sql[:i])
		_, err := buffer.WriteString(replace)
		if err != nil {
			return "", err
		}
		_, err = buffer.WriteString(sql[i+1:])
		return buffer.String(), err
	}

	return sql, errors.New("index is out of range")
}

//isSlice 判断接口是否为切片或数组
func isSlice(params interface{}) bool {
	kind := reflect.TypeOf(params).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		return true
	}
	return false
}
