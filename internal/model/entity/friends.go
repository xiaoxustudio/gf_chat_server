// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Friends is the golang structure for table friends.
type Friends struct {
	Id       int         `json:"id"        orm:"id"        ` // ID
	UserId   string      `json:"user_id"   orm:"user_id"   ` // 用户用户名
	FriendId string      `json:"friend_id" orm:"friend_id" ` // 好友用户名
	AddTime  *gtime.Time `json:"add_time"  orm:"add_time"  ` // 添加时间
	FriendData *User `json:"friend_data"  orm:"with:username=friend_id"  `
}
