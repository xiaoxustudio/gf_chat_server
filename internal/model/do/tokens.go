// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tokens is the golang structure of table tokens for DAO operations like Where/Data.
type Tokens struct {
	g.Meta      `orm:"table:tokens, do:true"`
	Id          interface{} // ID
	Token       interface{} // Token
	CreateTime  *gtime.Time // 创建时间
	FailureTime *gtime.Time // 失效时间
	TargetEmail interface{} // 目标邮箱
	Ip          interface{} // IP
}
