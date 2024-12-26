// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupConnect is the golang structure of table group-connect for DAO operations like Where/Data.
type GroupConnect struct {
	g.Meta  `orm:"table:group-connect, do:true"`
	Id      interface{} // ID
	UserId  interface{} // 用户名
	GroupId interface{} // 群组名称
	Auth    interface{} // 群组权限
	AddTime *gtime.Time // 加入时间
}
