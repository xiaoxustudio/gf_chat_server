package toArray

import "github.com/gogf/gf/v2/frame/g"

func ToArray(m g.Map) []map[string]interface{} {
	var slice []map[string]interface{}
	for k, v := range m {
		slice = append(slice, map[string]interface{}{k: v})
	}
	return slice
}
