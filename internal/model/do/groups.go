// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Groups is the golang structure of table groups for DAO operations like Where/Data.
type Groups struct {
	g.Meta      `orm:"table:groups, do:true"`
	Id          interface{} // ID
	GroupId     interface{} // 群组ID
	GroupStatus interface{} // 群组状态
	GroupName   interface{} // 群聊名称
	GroupDesc   interface{} // 群组简介
	GroupAvatar interface{} // 群组头像
	GroupMaster interface{} // 群主用户名
}
