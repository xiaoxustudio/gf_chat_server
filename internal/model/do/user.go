// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta       `orm:"table:user, do:true"`
	Id           interface{} // ID
	Nickname     interface{} // 用户昵称
	Username     interface{} // 用户用户名ID
	Password     interface{} // 用户密码
	Phone        interface{} // 用户手机号
	Email        interface{} // 用户邮箱
	RegisterTime *gtime.Time // 用户注册时间
	LoginTime    *gtime.Time // 用户最后登录时间
	Token        interface{} // 用户token
	Group        interface{} // 用户分组ID
	Avatar       interface{} // 用户头像
	EmailAuth    interface{} // 邮箱验证 0 未验证 1 验证
}
