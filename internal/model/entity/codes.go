// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Codes is the golang structure for table codes.
type Codes struct {
	Id          int         `json:"id"           orm:"id"           ` // ID
	Code        string      `json:"code"         orm:"code"         ` // Token
	CreateTime  *gtime.Time `json:"create_time"  orm:"create_time"  ` // 创建时间
	FailureTime *gtime.Time `json:"failure_time" orm:"failure_time" ` // 失效时间
	TargetEmail string      `json:"target_email" orm:"target_email" ` // 目标邮箱
	Ip          string      `json:"ip"           orm:"ip"           ` // IP
}
