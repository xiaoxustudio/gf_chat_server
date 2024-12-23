// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id           int         `json:"id"            orm:"id"            ` // ID
	Nickname     string      `json:"nickname"      orm:"nickname"      ` // 用户昵称
	Username     string      `json:"username"      orm:"username"      ` // 用户用户名ID
	Password     string      `json:"password"      orm:"password"      ` // 用户密码
	Phone        int         `json:"phone"         orm:"phone"         ` // 用户手机号
	Email        string      `json:"email"         orm:"email"         ` // 用户邮箱
	RegisterTime *gtime.Time `json:"register_time" orm:"register_time" ` // 用户注册时间
	LoginTime    *gtime.Time `json:"login_time"    orm:"login_time"    ` // 用户最后登录时间
	Token        string      `json:"token"         orm:"token"         ` // 用户token
	Group        int         `json:"group"         orm:"group"         ` // 用户分组ID
	Avatar       string      `json:"avatar"        orm:"avatar"        ` // 用户头像
}
