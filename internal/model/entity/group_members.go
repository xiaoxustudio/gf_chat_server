// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupMembers is the golang structure for table group_members.
type GroupMembers struct {
	Id      int         `json:"id"       orm:"id"       ` // ID
	UserId  int         `json:"user_id"  orm:"user_id"  ` // 用户名称
	GroupId string      `json:"group_id" orm:"group_id" ` // 加入组ID
	AddTime *gtime.Time `json:"add_time" orm:"add_time" ` // 加入时间
	Auth    int         `json:"auth"     orm:"auth"     ` // 权限：0普通，1管理员，2群主
}
