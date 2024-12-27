package websocketUnitGroup

import (
	"context"
	"encoding/json"
	"fmt"
	"gf_chat_server/internal/consts"
	"gf_chat_server/internal/model/entity"
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
	GroupID          string // 群组ID
}

// 连接池
type WebSocketPool struct {
	Connections map[*WebSocketConnection]bool
	Lock        sync.Mutex
}

type ChatItem struct {
	NickName  string              `json:"nickname"`
	Avatar    string              `json:"avatar"`
	SendID    string              `json:"send_id"`
	ReceiveID string              `json:"receive_id"`
	Content   string              `json:"content"`
	Time      string              `json:"time"`
	Files     []string            `json:"files"`
	Type      consts.ChatItemType `json:"type"`
}

// 群组聊天单元
type ChatToken struct {
	Users        []string // 连接的成员
	UsersAvatars []string
	ChatList     []ChatItem
	GroupInfo    *entity.Groups // 群组信息
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const HeartBeatDelay = time.Duration(10) * time.Second // 5s

// websocket
type WebSocketUnitGroup struct {
	WebSocketPool *WebSocketPool
	ChatListData  []ChatToken
}

func New() WebSocketUnitGroup {
	pool := &WebSocketPool{
		Connections: make(map[*WebSocketConnection]bool),
	}
	return WebSocketUnitGroup{WebSocketPool: pool, ChatListData: []ChatToken{}}
}

func (r *WebSocketUnitGroup) ValidToken(username string, Request *http.Request) bool {
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

func (r *WebSocketUnitGroup) HandleWebSocket(ResponseWriter http.ResponseWriter, Request *http.Request) {
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
		GroupID:  "",
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
func (r *WebSocketUnitGroup) HandleWebSocketMessage(conn *websocket.Conn, msg []byte) {
	var data scmsg.SCMsgO
	err := json.Unmarshal(msg, &data)
	res := r.GetWsForConn(conn)

	if err != nil || res == nil {
		return
	} else if data.Type == consts.HeartBeat { // 判断是否是心跳
		// tw.Tw(context.Background(), "接收心跳消息：%s", data.Message)
		res.HearbeatLastTime = gtime.Now().Unix()
	} else if data.Type == consts.CreateChannel { // 是否是建立连接
		tw.Tw(context.Background(), "消息通道创建：%s", data.Data.(g.Map)["group"])
		group_id := data.Data.(g.Map)["group"].(string)
		// 判断是否加入群聊
		md := g.Model(fmt.Sprintf("group-%s", group_id))
		resSelf, err := md.Clone().Where("user_id", res.UserName).One()
		if err != nil {
			// "失败"+err.Error()
			res.Conn.Close()
			return
		} else if len(resSelf) == 0 {
			// "你还未加入群聊"
			res.Conn.Close()
			return
		}
		res.GroupID = group_id
		index := r.getGroupList(res.GroupID)
		if index != -1 {
			ct := &r.ChatListData[index]
			ct.Users = append(ct.Users, res.UserName)
		}
		// 同步历史消息（如果有的话）
		r.SyncChat(res.GroupID)
	} else if data.Type == consts.Send { // 是否是发送聊天消息
		if res.GroupID != "" {
			// tw.Tw(context.Background(), "接收消息通道：%s", data.Data.(g.Map)["GroupID"])
			r.OnMessage(res, data)
		}
	} else if data.Type == consts.WithDraw { // 是否是撤回聊天消息
		if res.GroupID != "" && res.UserName != "" {
			index := data.Data.(g.Map)["index"].(float64)
			// tw.Tw(context.Background(), "撤回数据：%s %s %d", res.UserName, res.GroupID, index)
			r.RemoveChat(res.UserName, res.GroupID, int(index))
		}
	}
}

// 使用Conn获取Ws连接
func (r *WebSocketUnitGroup) GetWsForConn(conn *websocket.Conn) *WebSocketConnection {
	var pool = r.WebSocketPool
	for wsConn := range pool.Connections {
		if wsConn.Conn == conn {
			return wsConn
		}
	}
	return nil
}

// 使用Username获取Ws连接
func (r *WebSocketUnitGroup) GetWsForUsername(un string) *WebSocketConnection {
	var pool = r.WebSocketPool
	for wsConn := range pool.Connections {
		if wsConn.UserName == un {
			return wsConn
		}
	}
	return nil
}

func (r *WebSocketUnitGroup) sendHeartbeat() {
	var pool = r.WebSocketPool
	for {
		time.Sleep(HeartBeatDelay) // 心跳间隔时间

		// 使用读写锁来提高并发性能
		pool.Lock.Lock()
		connections := make(map[*WebSocketConnection]bool)
		for wsConn := range pool.Connections {
			connections[wsConn] = true
		}
		pool.Lock.Unlock()

		for wsConn := range connections {
			// 发送心跳消息
			var newData = scmsg.SCMsgO{Type: consts.HeartBeatServer, Message: "server heartbeat"}
			newDataBytes, err := json.Marshal(newData)
			if err != nil {
				// 记录错误，而不是直接返回
				tw.Tw(context.Background(), "Error marshalling heartbeat message: %v", err)
				continue
			}

			err = wsConn.Conn.WriteMessage(websocket.TextMessage, newDataBytes)
			if err != nil {
				// 错误处理：删除连接并关闭
				r.handleConnectionError(wsConn)
			} else {
				// 更新心跳状态
				wsConn.LastPing = time.Now()
			}
		}
	}
}

func (r *WebSocketUnitGroup) handleConnectionError(wsConn *WebSocketConnection) {
	// 锁定写操作
	r.WebSocketPool.Lock.Lock()
	defer r.WebSocketPool.Lock.Unlock()

	// 删除连接并关闭
	if _, ok := r.WebSocketPool.Connections[wsConn]; ok {
		delete(r.WebSocketPool.Connections, wsConn)
		wsConn.Conn.Close()
	}

	// 记录错误信息
	tw.Tw(context.Background(), "Error sending heartbeat to connection: %v", wsConn.Conn.RemoteAddr())
}

// 聊天消息处理

// 获取指定 群组id 的群组单元索引
func (r *WebSocketUnitGroup) getGroupList(groupId string) int {
	for i, it := range r.ChatListData {
		if it.GroupInfo.GroupId == groupId {
			return i
		}
	}
	return -1
}

// 聊天消息包裹
func (r *WebSocketUnitGroup) MsgWrapper(tp consts.WsCode, msg string, da any) g.Map {
	return g.Map{
		"type":    tp,
		"message": msg,
		"data":    da,
	}
}

// 接收到消息
func (r *WebSocketUnitGroup) OnMessage(connect *WebSocketConnection, data scmsg.SCMsgO) {
	// 将聊天信息存储到自身进程
	// 聊天记录将在双方的一方断开时自动写入
	tw.Tw(context.Background(), "接收群组聊天消息 to %s：%s", connect.GroupID, data.Message)
	res, err := g.Model("user").Where("username", connect.UserName).One()
	if err != nil {
		return
	}
	md := g.Model("groups")
	gres, err := md.Where("group_id", connect.GroupID).One()
	if err != nil {
		return
	}
	// 判断是否有群组记录
	index := r.getGroupList(connect.GroupID)
	if index == -1 {
		// 创建新记录并更新头像
		var listTk = ChatToken{
			Users:        []string{connect.UserName},
			UsersAvatars: []string{res.GMap().Get("avatar").(string)},
			ChatList:     []ChatItem{},
			GroupInfo: &entity.Groups{
				GroupId:     gres.GMap().Get("group_id").(string),
				GroupStatus: int(gres.GMap().Get("group_status").(int32)),
				GroupDesc:   gres.GMap().Get("group_desc").(string),
				GroupName:   gres.GMap().Get("group_name").(string),
				GroupAvatar: gres.GMap().Get("group_avatar").(string),
				GroupMaster: gres.GMap().Get("group_master").(string),
			},
		}
		r.ChatListData = append(r.ChatListData, listTk)
		index = len(r.ChatListData) - 1
	}

	// 将传入的[]interface{}转换为p[]string
	cacheInterface := data.Data.(g.Map)["files"].([]interface{})
	ImageStrings := make([]string, len(cacheInterface))
	for i, val := range cacheInterface {
		if str, ok := val.(string); ok {
			ImageStrings[i] = str
		}
	}

	ListToken := &r.ChatListData[index]
	ListToken.ChatList = append(ListToken.ChatList, ChatItem{
		Avatar:    res.GMap().Get("avatar").(string),
		NickName:  res.GMap().Get("nickname").(string),
		SendID:    connect.UserName,
		ReceiveID: connect.GroupID,
		Content:   data.Message,
		Time:      gtime.Datetime(),
		Type:      consts.Common,
		Files:     ImageStrings,
	})
	r.SyncChat(connect.GroupID)
}

// 获取指定用户的头像
func (r *WebSocketUnitGroup) GetChatAvatar(username string) string {
	res, err := g.Model("user").Where("username", username).One()
	if err != nil {
		return ""
	}
	return res.GMap().Get("avatar").(string)
}

// 更新消息记录头像
func (r *WebSocketUnitGroup) SyncChatAvatar(group_id string, username string) {
	index := r.getGroupList(group_id)
	if index == -1 {
		return
	}
	ListToken := r.ChatListData[index]
	// 循环发送
	for i, v := range ListToken.ChatList {
		if v.SendID == username {
			ListToken.ChatList[i].Avatar = r.GetChatAvatar(v.SendID)
		}
	}
}

// 同步聊天消息列表
func (r *WebSocketUnitGroup) SyncChat(group_id string) {
	index := r.getGroupList(group_id)
	if index == -1 {
		return
	}
	ListToken := r.ChatListData[index]
	NewDataBytes, err := json.Marshal(r.MsgWrapper(consts.UpdateMsgList, "ok", ListToken.ChatList))
	if err != nil {
		return
	}
	// 循环发送
	for _, v := range ListToken.Users {
		wsRes := r.GetWsForUsername(v)
		if wsRes == nil {
			return
		}
		r.SyncChatAvatar(group_id, v)
		wsRes.Conn.WriteMessage(websocket.TextMessage, NewDataBytes)
	}
}

// 移除指定用户的聊天消息列表
func (r *WebSocketUnitGroup) RemoveChatSlot(un string) {
	index := r.getGroupList(un)
	if index == -1 {
		return
	}
	ListToken := r.ChatListData[index]
	UserName := ListToken.Users[0]
	GroupIDName := ListToken.Users[1]
	link1 := r.GetWsForUsername(UserName)
	link2 := r.GetWsForUsername(GroupIDName)
	if link1 == nil && link2 == nil {
		// 双方都断开连接，保存记录到数据库，并删除内存聊天数据
		// r.ChatListData = append(r.ChatListData[:index-1], r.ChatListData[index:]...)
	}
}

// 移除指定用户的指定聊天消息
func (r *WebSocketUnitGroup) RemoveChat(un string, group_id string, chatIndex int) {
	index := r.getGroupList(group_id)
	if index == -1 {
		return
	}
	ListToken := &r.ChatListData[index]
	list := ListToken.ChatList
	if chatIndex >= 0 && chatIndex < len(list) {
		UserName := ListToken.Users[0]
		res, err := g.Model("user").Where("username", UserName).One()
		if err != nil {
			return
		}
		back, err := array.InsertIntoArray(list, chatIndex, ChatItem{
			Type:    consts.System,
			Content: fmt.Sprintf("%s(%s) 撤回了一条消息！", res.GMap().Get("nickname"), UserName),
		})

		if err == nil {
			// 添加撤回系统消息
			ListToken.ChatList = back
		}

		// 移除指定索引的聊天消息
		ListToken.ChatList = array.RemoveItemFromArray(ListToken.ChatList, chatIndex+1)
		// 同步聊天
		r.SyncChat(group_id)
	}
}
