// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DocumentTemplate is the golang structure of table document-template for DAO operations like Where/Data.
type DocumentTemplate struct {
	g.Meta  `orm:"table:document-template, do:true"`
	Id      interface{} // ID
	UserId  interface{} // 用户名
	Auth    interface{} // 用户权限：0 可查看 1可编辑 2可管理
	AddTime *gtime.Time // 添加时间
}
