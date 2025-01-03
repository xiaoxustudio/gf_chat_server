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
	NickName  string                `json:"nickname"`
	Avatar    string                `json:"avatar"`
	SendID    string                `json:"send_id"`
	ReceiveID string                `json:"receive_id"`
	Content   string                `json:"content"`
	Time      string                `json:"time"`
	Files     []string              `json:"files"`
	Type      consts.ChatItemType   `json:"type"`
	SendData  *entity.GroupTemplate `json:"send_data"`
}

// 群组聊天单元
type ChatToken struct {
	Users     []string // 连接的成员
	ChatList  []ChatItem
	GroupInfo *entity.Groups // 群组信息
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const HeartBeatDelay = time.Duration(10) * time.Second // 10s

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
	defer func() {
		// 连接关闭时从连接池中移除
		pool.Lock.Lock()
		// 将users中剔除用户
		var group_id = wsConn.GroupID
		index := r.getGroupList(group_id)
		if index == -1 {
			return
		}
		ListToken := &r.ChatListData[index]
		var NewUsers = make([]string, len(ListToken.Users))
		for _, v := range ListToken.Users {
			if v != wsConn.UserName {
				NewUsers = append(NewUsers, v)
			}
		}
		ListToken.Users = NewUsers
		delete(pool.Connections, wsConn)
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
			conn.Close()
			return
		} else if len(resSelf) == 0 {
			// "你还未加入群聊"
			conn.Close()
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
			Users:    []string{connect.UserName},
			ChatList: []ChatItem{},
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
	md = g.Model(fmt.Sprintf("group-%s", connect.GroupID))
	var mdata entity.GroupTemplate
	err = md.Clone().Where("user_id", connect.UserName).With(&entity.User{}).Scan(&mdata)
	if err != nil {
		return
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
		SendData: &entity.GroupTemplate{
			UserId:   mdata.UserId,
			Auth:     mdata.Auth,
			AddTime:  mdata.AddTime,
			UserData: mdata.UserData,
		},
	})
	r.SyncChat(connect.GroupID)
}

// 传入group_id，获取所有用户记录
func (r *WebSocketUnitGroup) getChatUserData(group_id string) g.Map {
	md := g.Model(fmt.Sprintf("group-%s", group_id))
	var mdata []*entity.GroupTemplate
	arrs := make(g.Map, len(mdata))

	err := md.With(&entity.User{}).Scan(&mdata)
	if err != nil {
		return arrs
	}
	for _, data := range mdata {
		username := data.UserId
		arrs[username] = data
	}
	return arrs
}

// 更新同步用户记录
func (r *WebSocketUnitGroup) SyncChatUserData(group_id string) {
	index := r.getGroupList(group_id)
	if index == -1 {
		return
	}
	ListToken := &r.ChatListData[index]
	data := r.getChatUserData(group_id)
	// 循环发送
	for i, v := range ListToken.ChatList {
		var item = &ListToken.ChatList[i]
		// 普通消息类型
		if item.Type == consts.Common {
			item.SendData = data[v.SendData.UserId].(*entity.GroupTemplate)
		}
	}
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
	// 更新user_data
	r.SyncChatUserData(group_id)
	// 循环发送
	for _, v := range ListToken.Users {
		wsRes := r.GetWsForUsername(v)
		if wsRes != nil {
			r.SyncChatAvatar(group_id, v)
			wsRes.Conn.WriteMessage(websocket.TextMessage, NewDataBytes)
		}
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

	// 如果chatIndex无效
	if chatIndex < 0 || chatIndex >= len(list) {
		return
	}

	// 检查要移除的消息是否属于指定用户
	if list[chatIndex].SendData.UserId != un {
		return
	}
	res, err := g.Model("user").Where("username", un).One()

	if err != nil {
		return
	}

	// 插入撤回系统消息
	nickname := res.GMap().Get("nickname")
	withdrawMessage := ChatItem{
		Type:    consts.System,
		Content: fmt.Sprintf("%s(%s) 撤回了一条消息！", nickname, un),
	}

	back, err := array.InsertIntoArray(list, chatIndex, withdrawMessage)

	if err == nil {
		// 添加撤回系统消息
		ListToken.ChatList = back
	}

	// 移除指定索引的聊天消息
	ListToken.ChatList = array.RemoveItemFromArray(ListToken.ChatList, chatIndex+1)

	// 同步聊天
	r.SyncChat(group_id)
}
