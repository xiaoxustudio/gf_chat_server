// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Documents is the golang structure of table documents for DAO operations like Where/Data.
type Documents struct {
	g.Meta  `orm:"table:documents, do:true"`
	Id      interface{} // ID
	Block   interface{} // blockID
	UserId  interface{} // 属于用户
	Type    interface{} // 类型 0 文件夹 1 页面
	Content interface{} // 页面内容
	Status  interface{} // block状态 0 锁定 1 可编辑
}
