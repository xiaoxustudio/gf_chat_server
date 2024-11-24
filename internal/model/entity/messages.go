// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Messages is the golang structure for table messages.
type Messages struct {
	MessageId int         `json:"message_id" orm:"message_id" ` // 消息ID
	SendId    int         `json:"send_id"    orm:"send_id"    ` // 发送方ID
	ReceiveId int         `json:"receive_id" orm:"receive_id" ` // 接收方ID
	Content   string      `json:"content"    orm:"content"    ` // 内容
	SnedTime  *gtime.Time `json:"sned_time"  orm:"sned_time"  ` // 发送时间
}
