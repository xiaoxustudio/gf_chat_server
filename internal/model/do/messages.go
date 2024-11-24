// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Messages is the golang structure of table messages for DAO operations like Where/Data.
type Messages struct {
	g.Meta    `orm:"table:messages, do:true"`
	MessageId interface{} // 消息ID
	SendId    interface{} // 发送方ID
	ReceiveId interface{} // 接收方ID
	Content   interface{} // 内容
	SnedTime  *gtime.Time // 发送时间
}
