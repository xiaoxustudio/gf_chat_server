// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupMembers is the golang structure of table group_members for DAO operations like Where/Data.
type GroupMembers struct {
	g.Meta  `orm:"table:group_members, do:true"`
	Id      interface{} // ID
	UserId  interface{} // 用户名称
	GroupId interface{} // 加入组ID
	AddTime *gtime.Time // 加入时间
	Auth    interface{} // 权限：0普通，1管理员，2群主
}
