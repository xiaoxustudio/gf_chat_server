package envtool

import (
	"os"
	"strconv"
)

// GetEnvValue 获取环境变量，并尝试将其转换为指定类型T
func GetEnvValue[T any](key string, defaultValue T) any {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		var result any
		var errs error
		switch any(defaultValue).(type) {
		case int:
			intValue, err := strconv.Atoi(value)
			errs = err
			result = intValue // 将int转换为T类型
		case string:
			result = string(value)
		case bool:
			boolValue, err := strconv.ParseBool(value)
			errs = err
			result = boolValue
		case float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			errs = err
			result = floatValue
		default:
			return defaultValue
		}
		if errs != nil {
			return defaultValue
		}
		return result
	}
	return defaultValue
}
