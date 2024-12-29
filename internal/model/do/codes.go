// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Codes is the golang structure of table codes for DAO operations like Where/Data.
type Codes struct {
	g.Meta      `orm:"table:codes, do:true"`
	Id          interface{} // ID
	Code        interface{} // Token
	CreateTime  *gtime.Time // 创建时间
	FailureTime *gtime.Time // 失效时间
	TargetEmail interface{} // 目标邮箱
	Ip          interface{} // IP
}
