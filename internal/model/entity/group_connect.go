// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupConnect is the golang structure for table group-connect.
type GroupConnect struct {
	Id      int         `json:"id"       orm:"id"       ` // ID
	UserId  string      `json:"user_id"  orm:"user_id"  ` // 用户名
	GroupId string      `json:"group_id" orm:"group_id" ` // 群组名称
	Auth    int         `json:"auth"     orm:"auth"     ` // 群组权限
	AddTime *gtime.Time `json:"add_time" orm:"add_time" ` // 加入时间
	GroupData *Groups `json:"group_data"  orm:"with:group_id=group_id"  `
}
