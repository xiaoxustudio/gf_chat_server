package websocketunit

import (
	"context"
	"encoding/json"
	"gf_chat_server/internal/consts"
	array "gf_chat_server/utility/array"
	scmsg "gf_chat_server/utility/scMsg"
	"gf_chat_server/utility/token"
	"gf_chat_server/utility/tw"
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
)

// WebSocket连接结构体
type WebSocketConnection struct {
	Conn             *websocket.Conn
	LastPing         time.Time
	UserName         string
	HearbeatLastTime int64
	Target           string
}

// 连接池
type WebSocketPool struct {
	Connections map[*WebSocketConnection]bool
	Lock        sync.Mutex
}

type ChatItem struct {
	NickName  string `json:"nickname"`
	SendID    string `json:"send_id"`
	ReceiveID string `json:"receive_id"`
	Content   string `json:"content"`
	Time      string `json:"time"`
}

// 聊天单元
type ChatToken struct {
	Users    []string
	ChatList []ChatItem
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const HeartBeatDelay = 30 * time.Second

// websocket
type WebSocketUnit struct {
	WebSocketPool *WebSocketPool
	ChatListData  []ChatToken
}

func New() WebSocketUnit {
	pool := &WebSocketPool{
		Connections: make(map[*WebSocketConnection]bool),
	}
	return WebSocketUnit{WebSocketPool: pool, ChatListData: []ChatToken{}}
}

func (r *WebSocketUnit) ValidToken(username string, Request *http.Request) bool {
	md := g.Model("user")
	tok := Request.URL.Query().Get("token")
	if len(tok) == 0 {
		return false
	}
	res, err := token.ValidToken(tok)
	if err == nil {
		// 进行数据库比对
		r, err := md.Where("token", tok).Where("username", username).All()
		if err == nil && len(r) > 0 {
			return res
		} else {
			return false
		}
	} else {
		return false
	}
}

func (r *WebSocketUnit) HandleWebSocket(ResponseWriter http.ResponseWriter, Request *http.Request) {
	var username = Request.URL.Query().Get("user")
	vaild := r.ValidToken(username, Request)
	if !vaild {
		return
	}
	conn, err := upgrader.Upgrade(ResponseWriter, Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	var pool = r.WebSocketPool
	// 创建新的WebSocket连接实例
	wsConn := &WebSocketConnection{
		Conn:     conn,
		LastPing: time.Now(),
		UserName: username,
		Target:   "",
	}
	// 将连接添加到连接池
	pool.Lock.Lock()
	pool.Connections[wsConn] = true
	pool.Lock.Unlock()

	// 处理连接消息
	go func() {
		defer func() {
			// 连接关闭时从连接池中移除
			pool.Lock.Lock()
			delete(pool.Connections, wsConn)
			r.RemoveChatSlot(wsConn.UserName)
			pool.Lock.Unlock()
			conn.Close()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// 更新心跳状态
			wsConn.LastPing = time.Now()
			r.HandleWebSocketMessage(conn, message)
		}
	}()
	r.sendHeartbeat()
}

// 处理message
func (r *WebSocketUnit) HandleWebSocketMessage(conn *websocket.Conn, msg []byte) {
	var data scmsg.SCMsgO
	err := json.Unmarshal(msg, &data)
	res := r.GetWsForConn(conn)

	if err != nil || res == nil {
		return
	} else if data.Type == consts.HeartBeat { // 判断是否是心跳
		// tw.Tw(context.Background(), "接收心跳消息：%s", data.Message)
		res.HearbeatLastTime = gtime.Now().Unix()
	} else if data.Type == consts.CreateChannel { // 是否是建立连接
		tw.Tw(context.Background(), "消息通道创建：%s", data.Data.(g.Map)["target"])
		targetName := data.Data.(g.Map)["target"].(string)
		res.Target = targetName
		// 同步历史消息（如果有的话）
		r.SyncChat(res.UserName)
	} else if data.Type == consts.Send { // 是否是发送聊天消息
		if res.Target != "" {
			// tw.Tw(context.Background(), "接收消息通道：%s", data.Data.(g.Map)["target"])
			r.OnMessage(res, data)
		}
	} else if data.Type == consts.WithDraw { // 是否是撤回聊天消息
		if res.Target != "" && res.UserName != "" {
			index := data.Data.(g.Map)["index"].(float64)
			r.RemoveChat(res.UserName, res.Target, int(index))
		}
	}
}

// 使用Conn获取Ws连接
func (r *WebSocketUnit) GetWsForConn(conn *websocket.Conn) *WebSocketConnection {
	var pool = r.WebSocketPool
	for wsConn := range pool.Connections {
		if wsConn.Conn == conn {
			return wsConn
		}
	}
	return nil
}

// 使用Username获取Ws连接
func (r *WebSocketUnit) GetWsForUsername(un string) *WebSocketConnection {
	var pool = r.WebSocketPool
	for wsConn := range pool.Connections {
		if wsConn.UserName == un {
			return wsConn
		}
	}
	return nil
}

func (r *WebSocketUnit) sendHeartbeat() {
	var pool = r.WebSocketPool
	for {
		time.Sleep(HeartBeatDelay) // 心跳间隔时间

		pool.Lock.Lock()
		for wsConn := range pool.Connections {
			// 发送心跳消息
			var NewData = scmsg.SCMsgO{Type: consts.HeartBeatServer, Message: "server hearbeat"}
			NewDataBytes, err := json.Marshal(NewData)
			if err != nil {
				return
			}
			err = wsConn.Conn.WriteMessage(websocket.TextMessage, NewDataBytes)
			if err != nil {
				delete(pool.Connections, wsConn)
				wsConn.Conn.Close()
			} else {
				// 更新心跳状态
				wsConn.LastPing = time.Now()
			}
		}
		pool.Lock.Unlock()
	}
}

// 聊天消息处理

// 获取消息索引
func (r *WebSocketUnit) getList(names ...string) int {
	userName := names[0]
	if len(names) == 1 {
		for i, it := range r.ChatListData {
			if it.Users[0] == userName || it.Users[1] == userName {
				return i
			}
		}
		return -1
	}
	targetName := names[1]
	for i, it := range r.ChatListData {
		if it.Users[0] == userName || it.Users[1] == targetName {
			return i
		}
	}
	return -1
}

// 聊天消息包裹
func (r *WebSocketUnit) MsgWrapper(tp consts.WsCode, msg string, da any) g.Map {
	return g.Map{
		"type":    tp,
		"message": msg,
		"data":    da,
	}
}

// 接收到消息
func (r *WebSocketUnit) OnMessage(connect *WebSocketConnection, data scmsg.SCMsgO) {
	// 将聊天信息存储到自身进程
	// 聊天记录将在双方的一方断开时自动写入
	tw.Tw(context.Background(), "接收聊天消息 to %s：%s", connect.Target, data.Message)
	// 判断是否有记录
	index := r.getList(connect.UserName)
	if index == -1 {
		// 创建新记录
		var listTk = ChatToken{
			Users:    []string{connect.UserName, connect.Target},
			ChatList: []ChatItem{},
		}
		r.ChatListData = append(r.ChatListData, listTk)
		index = len(r.ChatListData) - 1
	}
	md := g.Model("user")
	res, err := md.Where("username", connect.UserName).One()
	if err != nil {
		return
	}
	ListToken := &r.ChatListData[index]
	ListToken.ChatList = append(ListToken.ChatList, ChatItem{
		NickName:  res.GMap().Get("nickname").(string),
		SendID:    connect.UserName,
		ReceiveID: connect.Target,
		Content:   data.Message,
		Time:      gtime.Datetime()})
	r.SyncChat(connect.UserName)
	r.SyncChat(connect.Target)
}

// 同步聊天消息列表
func (r *WebSocketUnit) SyncChat(un string) {
	index := r.getList(un)
	if index == -1 {
		return
	}
	ListToken := r.ChatListData[index]
	NewDataBytes, err := json.Marshal(r.MsgWrapper(consts.UpdateMsgList, "ok", ListToken.ChatList))
	if err != nil {
		return
	}
	wsRes := r.GetWsForUsername(un)
	if wsRes == nil {
		return
	}
	wsRes.Conn.WriteMessage(websocket.TextMessage, NewDataBytes)
}

// 移除指定用户的聊天消息列表
func (r *WebSocketUnit) RemoveChatSlot(un string) {
	index := r.getList(un)
	if index == -1 {
		return
	}
	ListToken := r.ChatListData[index]
	UserName := ListToken.Users[0]
	TargetName := ListToken.Users[1]
	link1 := r.GetWsForUsername(UserName)
	link2 := r.GetWsForUsername(TargetName)
	if link1 == nil && link2 == nil {
		// 双方都断开连接，保存记录到数据库，并删除内存聊天数据
		// r.ChatListData = append(r.ChatListData[:index-1], r.ChatListData[index:]...)
	}
}

// 移除指定用户的指定聊天消息
func (r *WebSocketUnit) RemoveChat(un string, tg string, chatIndex int) {
	index := r.getList(un, tg)
	if index == -1 {
		return
	}
	ListToken := &r.ChatListData[index]
	list := ListToken.ChatList
	if chatIndex >= 0 && chatIndex < len(list) {
		// 移除指定索引的聊天消息
		ListToken.ChatList = array.RemoveItemFromArray(list, chatIndex)
		// 同步聊天
		UserName := ListToken.Users[0]
		TargetName := ListToken.Users[1]
		r.SyncChat(UserName)
		r.SyncChat(TargetName)
	}
}
