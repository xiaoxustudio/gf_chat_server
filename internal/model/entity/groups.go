// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Groups is the golang structure for table groups.
type Groups struct {
	Id          int    `json:"id"           orm:"id"           ` // ID
	GroupId     string `json:"group_id"     orm:"group_id"     ` // 群组ID
	GroupStatus int    `json:"group_status" orm:"group_status" ` // 群组状态
	GroupName   string `json:"group_name"   orm:"group_name"   ` // 群聊名称
	GroupDesc   string `json:"group_desc"   orm:"group_desc"   ` // 群组简介
	GroupAvatar string `json:"group_avatar" orm:"group_avatar" ` // 群组头像
	GroupMaster string `json:"group_master" orm:"group_master" ` // 群主用户名
}
