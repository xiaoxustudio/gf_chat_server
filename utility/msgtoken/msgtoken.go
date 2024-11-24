package msgtoken

import "github.com/gogf/gf/v2/frame/g"

type Msg struct {
	code int
	msg  string
	data interface{}
}

func MsgToken(code int, msg string, data interface{}) Msg {
	return Msg{
		code: code,
		msg:  msg,
		data: data,
	}
}
func ToGMap(c Msg) g.Map {
	return g.Map{
		"code": c.code,
		"msg":  c.msg,
		"data": c.data,
	}
}
func ToArray(dataMap map[string]interface{}) []any {
	var dataSlice []any

	// 遍历 map 并将值添加到 slice 中
	for _, value := range dataMap {
		dataSlice = append(dataSlice, value)
	}
	return dataSlice
}
