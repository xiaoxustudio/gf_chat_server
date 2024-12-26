// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupTemplate is the golang structure for table group-template.
type GroupTemplate struct {
	Id           int         `json:"id"             orm:"id"             ` // ID
	UserId       string      `json:"user_id"        orm:"user_id"        ` // 用户名
	Auth         int         `json:"auth"           orm:"auth"           ` // 用户权限：0 普通 1管理 2群主
	AddTime      *gtime.Time `json:"add_time"       orm:"add_time"       ` // 加入群聊时间
	LastChatTime *gtime.Time `json:"last_chat_time" orm:"last_chat_time" ` // 最后发言时间
}
