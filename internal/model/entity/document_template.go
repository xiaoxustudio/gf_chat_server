// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DocumentTemplate is the golang structure for table document-template.
type DocumentTemplate struct {
	Id      int         `json:"id"       orm:"id"       ` // ID
	UserId  string      `json:"user_id"  orm:"user_id"  ` // 用户名
	Auth    int         `json:"auth"     orm:"auth"     ` // 用户权限：0 可查看 1可编辑 2可管理
	AddTime *gtime.Time `json:"add_time" orm:"add_time" ` // 添加时间
}
