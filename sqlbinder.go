package gosql

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

//GetPreparedStatement 将传入的sql与params参数 整理成 prepare statement
func GetPreparedStatement(sql string, params map[string]interface{}) (string, []interface{}, error) {

	//拼接参数
	sql, params, err := getStatement([]byte(sql), params)
	if err != nil {
		return "", nil, err
	}

	//占位符替换参数
	sql, names, err := getPrepared(sql)
	if err != nil {
		return "", nil, err
	}

	return bindMap(sql, names, params)
}

// bindMap 绑定 map
func bindMap(sql string, names []string, params map[string]interface{}) (string, []interface{}, error) {

	paramsArr := make([]interface{}, 0, len(params)*2)
	var err error
	for i, name := range names {
		v, ok := params[name]
		if !ok {
			return "", nil, errors.New("lost params " + string(name))
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

// getStatement 获取 sql statement 处理字符串拼接
func getStatement(sqlStr []byte, datamap map[string]interface{}) (string, map[string]interface{}, error) {

	//保证最后一个参数也能被处理
	sqlStr = append(sqlStr, ' ')

	found := false
	name := make([]byte, 0, 12)
	sqlResult := make([]byte, 0, len(sqlStr))

	for i, b := range sqlStr {
		if b == '#' {
			//若存在 '##' 则转译为 #
			if sqlStr[i-1] == '#' {
				sqlResult = append(sqlResult, '#')
				found = false
				continue
			}
			if found {
				return "", nil, errors.New("can`t insert # in the params")
			}
			found = true
			continue
		}

		if !found {
			sqlResult = append(sqlResult, b)
		}

		//sql 字符串中 属性名 必须为仅包含  '.' '_' 'A-Z' 'a-z' '0-9'
		if found && b != '.' && b != '_' && !unicode.IsLetter(rune(b)) && !unicode.IsNumber(rune(b)) {
			value, ok := datamap[string(name)]
			if !ok {
				return "", nil, errors.New("lost params " + string(name))
			}
			paramValue := transToString(value)
			delete(datamap, string(name))
			sqlResult = bytes.Join([][]byte{sqlResult, []byte(paramValue)}, nil)
			sqlResult = append(sqlResult, b)
			found = false
			name = make([]byte, 0, 12)
		}

		if found {
			name = append(name, b)
		}

	}

	return string(sqlResult[:len(sqlResult)-1]), datamap, nil
}

//getPreparedStatement 转换 sql 字符串为 Prepared Statement 并将 属性名 提取出来
func getPrepared(sqlStr string) (string, []string, error) {

	begin := 0
	found := false
	names := make([]string, 0, 10)
	sqlResult := make([]rune, 0, len(sqlStr))

	for i, b := range sqlStr {
		if b == '$' {
			//若存在 '$$' 则转译为 $
			if sqlStr[i-1] == '$' {
				sqlResult[len(sqlResult)-1] = '$'
				found = false
				continue
			}
			//若在属性名中发现 $ 则报错
			if found {
				return "", nil, errors.New("can`t insert $ in the params")
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
			names = append(names, sqlStr[begin:i])
			found = false
		}

	}
	if found {
		names = append(names, sqlStr[begin:])
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

// TransToString 强制类型转换为 string
func transToString(data interface{}) (res string) {
	switch v := data.(type) {
	case bool:
		res = strconv.FormatBool(v)
	case float32:
		res = strconv.FormatFloat(float64(v), 'f', 6, 32)
	case float64:
		res = strconv.FormatFloat(v, 'f', 6, 64)
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(data)
		res = strconv.FormatInt(val.Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(data)
		res = strconv.FormatUint(val.Uint(), 10)
	case string:
		res = v
	case []byte:
		res = string(v)
	default:
		res = fmt.Sprintf("%v", v)
	}
	return
}
