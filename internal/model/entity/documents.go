// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity


// Documents is the golang structure for table documents.
type Documents struct {
	Id        int    `json:"id"         orm:"id"         ` // ID
	Block     string `json:"block"      orm:"block"      ` // blockID
	BlockDesc string `json:"block_desc" orm:"block_desc" ` // 页面简介
	BlockName string `json:"block_name" orm:"block_name" ` // 页面名称
	UserId    string `json:"user_id"    orm:"user_id"    ` // 属于用户
	Type      int    `json:"type"       orm:"type"       ` // 类型 0 文件夹 1 页面
	Content   string `json:"content"    orm:"content"    ` // 页面内容
	Status    int    `json:"status"     orm:"status"     ` // block状态 0 锁定 1 可编辑
	DocData any `json:"doc_data"         ` 
}