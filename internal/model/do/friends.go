// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Friends is the golang structure of table friends for DAO operations like Where/Data.
type Friends struct {
	g.Meta   `orm:"table:friends, do:true"`
	Id       interface{} // ID
	UserId   interface{} // 用户用户名
	FriendId interface{} // 好友用户名
	AddTime  *gtime.Time // 添加时间
}
