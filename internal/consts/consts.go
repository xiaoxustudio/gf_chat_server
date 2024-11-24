package consts

type ErrorCode int
type WsCode int

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
	// 建立连接
	CreateChannel
	DisConnectChannel
	// 更新消息列表
	UpdateMsgList
)

func (t ErrorCode) String() string {
	names := map[ErrorCode]string{
		TokenExpired: "Expired",
	}
	return names[t]
}
func (t WsCode) String() string {
	names := map[WsCode]string{
		HeartBeat: "心跳链接",
	}
	return names[t]
}
