package consts

type ErrorCode int    // 错误码
type WsCode int       // websocket 码
type ChatItemType int // 聊天数据类型
type DocType int      // 文档类型

const (
	Status ErrorCode = iota + 1
	// 表单错误
	DataEmpty = 100
	// token错误
	TokenEmpty   = 200
	TokenInValid = 201
	TokenExpired = 202

	// 成功
	Success = 1
)
const (
	WSStatus WsCode = 800 + iota
	// 心跳
	HeartBeat
	HeartBeatServer
	HeartBeatClient
	// 消息
	Receive
	Send
	WithDraw // 撤回
	// 建立连接
	CreateChannel
	DisConnectChannel
	// 更新消息列表
	UpdateMsgList

	// 修改文档内容
	ChangeContent
	// 修改文档标题
	ChangeTitle
)

const (
	chatType ChatItemType = 900 + iota
	Common                // 普通消息类型
	System                // 系统消息类型
)

func (t ErrorCode) String() string {
	names := map[ErrorCode]string{
		TokenExpired: "Expired",
	}
	return names[t]
}
func (t WsCode) String() string {
	names := map[WsCode]string{
		HeartBeat:       "心跳链接",
		HeartBeatServer: "心跳链接（服务器）",
		HeartBeatClient: "心跳链接（客户端）",
		// 消息
		Receive:           "接收消息",
		Send:              "发送消息",
		WithDraw:          "撤回消息",
		CreateChannel:     "建立连接",
		DisConnectChannel: "断开连接",
		UpdateMsgList:     "更新消息列表",
	}
	return names[t]
}
func (t ChatItemType) String() string {
	names := map[ChatItemType]string{
		Common: "普通消息",
		System: "系统消息",
	}
	return names[t]
}
