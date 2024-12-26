// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupTemplate is the golang structure of table group-template for DAO operations like Where/Data.
type GroupTemplate struct {
	g.Meta       `orm:"table:group-template, do:true"`
	Id           interface{} // ID
	UserId       interface{} // 用户名
	Auth         interface{} // 用户权限：0 普通 1管理 2群主
	AddTime      *gtime.Time // 加入群聊时间
	LastChatTime *gtime.Time // 最后发言时间
}
