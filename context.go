package requtil

import (
	"context"
	"fmt"
	"net/http"
)

type keyType string

//Set 将 k-v 存入 context
func Set(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, keyType(key), value)
}

//Get 从 Context 获取 key 对应的 value
func Get(ctx context.Context, key string) (interface{}, error) {
	result := ctx.Value(keyType(key))
	if result == nil {
		return 0, fmt.Errorf("can`t find %s in context", key)
	}
	return result, nil
}

//GetUInt64 从 context 中获取 string 类型的参数
func GetUInt64(ctx context.Context, key string) (uint64, error) {
	value, err := Get(ctx, key)
	if err != nil {
		return 0, nil
	}
	if data, ok := value.(uint64); ok {
		return data, nil
	}

	return 0, fmt.Errorf("%s is not type string", key)
}

//GetString 从 context 中获取 string 类型的参数
func GetString(ctx context.Context, key string) (string, error) {
	value, err := Get(ctx, key)
	if err != nil {
		return "", nil
	}
	if data, ok := value.(string); ok {
		return data, nil
	}

	return "", fmt.Errorf("%s is not type string", key)
}

//GetHTTPHeader 从 context 中获取 http.Header 类型的参数
func GetHTTPHeader(ctx context.Context, key string) (http.Header, error) {
	value, err := Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data, ok := value.(http.Header); ok {
		return data, nil
	}

	return nil, fmt.Errorf("%s is not type Header", key)
}

//SetLogID 将 logID 存入 context 中
func SetLogID(ctx context.Context, logID string) context.Context {
	return Set(ctx, "logID", logID)
}

//GetLogID 从 context 中获取 LogID
func GetLogID(ctx context.Context) (string, error) {
	logID, err := GetString(ctx, "logID")
	if err != nil {
		return "", nil
	}
	return logID, err
}

//SetHeader 将 Header 存入 context 中
func SetHeader(ctx context.Context, header http.Header) context.Context {
	return Set(ctx, "header", header)
}

//GetHeader 从 context 中获取 header
func GetHeader(ctx context.Context) (http.Header, error) {
	h, err := GetHTTPHeader(ctx, "header")
	if err != nil {
		return nil, err
	}
	return h, nil
}
