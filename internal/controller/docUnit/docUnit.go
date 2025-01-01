package docunit

import (
	"context"
	"encoding/json"
	"gf_chat_server/internal/consts"
	"gf_chat_server/internal/dao"
	"gf_chat_server/internal/model/entity"
	scmsg "gf_chat_server/utility/scMsg"
	"gf_chat_server/utility/token"
	"gf_chat_server/utility/tw"
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

// 文档WebSocket连接结构体
type WebSocketConnection struct {
	Conn      *websocket.Conn
	User      entity.User
	Block     string
	UserName  string
	BlockData entity.Documents
}

// 连接池
type WebSocketPool struct {
	Connections map[*websocket.Conn]*WebSocketConnection
	Lock        sync.Mutex
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
type DocUnit struct {
	WebSocketPool *WebSocketPool
}

func New() DocUnit {
	pool := &WebSocketPool{
		Connections: make(map[*websocket.Conn]*WebSocketConnection),
	}
	return DocUnit{WebSocketPool: pool}
}

func (r *DocUnit) ValidToken(username string, Request *http.Request) bool {
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

// 使用Conn获取Ws连接
func (r *DocUnit) GetWsForConn(conn *websocket.Conn) *WebSocketConnection {
	var pool = r.WebSocketPool
	for i, wsConn := range pool.Connections {
		if conn == i {
			return wsConn
		}
	}
	return nil
}

// 使用Username获取Ws连接
func (r *DocUnit) GetWsForUsername(un string) *WebSocketConnection {
	var pool = r.WebSocketPool
	for _, wsConn := range pool.Connections {
		if wsConn.UserName == un {
			return wsConn
		}
	}
	return nil
}

func (r *DocUnit) HandleWebSocket(ResponseWriter http.ResponseWriter, Request *http.Request) {
	var query = Request.URL.Query()
	var username = query.Get("user")
	var block_id = query.Get("block")
	vaild := r.ValidToken(username, Request)
	if !vaild {
		return
	}
	conn, err := upgrader.Upgrade(ResponseWriter, Request, nil)
	if err != nil {
		return
	}
	var pool = r.WebSocketPool
	var udata entity.User
	err = g.Model("user").Where("username", username).Scan(&udata)
	if err != nil {
		conn.Close()
		return
	}
	// 创建新的WebSocket连接实例
	wsConn := &WebSocketConnection{
		Conn:     conn,
		User:     udata,
		UserName: udata.Username,
		Block:    block_id,
	}
	// 将连接添加到连接池
	pool.Lock.Lock()
	pool.Connections[conn] = wsConn
	pool.Lock.Unlock()

	go func() {

		// 处理连接消息
		defer func() {
			// 保存文档
			pool.Lock.Lock()
			md := g.Model("documents")
			_, err = md.Clone().Where("user_id", wsConn.UserName).Where("block", wsConn.Block).Update(wsConn.BlockData)
			if err != nil {
				tw.Tw(context.Background(), "（错误）文档未保存：%s ", err)
			}
			// 从连接池中移除
			tw.Tw(context.Background(), "文档通信退出：%s", wsConn.UserName)
			delete(pool.Connections, wsConn.Conn)
			pool.Lock.Unlock()
			conn.Close()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break // 退出循环
			}
			// 处理消息
			r.HandleWebSocketMessage(conn, message)
		}
	}()
	// r.sendHeartbeat()
}

// 处理message
func (r *DocUnit) HandleWebSocketMessage(conn *websocket.Conn, msg []byte) {
	var pool = r.WebSocketPool
	var data scmsg.SCMsgO
	err := json.Unmarshal(msg, &data)
	res := r.GetWsForConn(conn)
	if err != nil || res == nil {
		conn.Close()
		return
	} else if data.Type == consts.CreateChannel {
		tw.Tw(context.Background(), "创建通道：%s", data.Message)
		// 验证
		pool.Lock.Lock()
		md := dao.Documents.Ctx(context.Background())
		var docData entity.Documents
		err := md.Clone().Where("block", res.Block).Scan(&docData)
		if err == nil {
			// 初始化
			res.BlockData = docData
			var NewData = scmsg.SCMsgO{Type: consts.HeartBeatServer, Message: docData.Content, Data: docData}
			NewDataBytes, err := json.Marshal(NewData)
			if err != nil {
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, NewDataBytes)
			if err != nil {
				conn.Close()
			}
			pool.Lock.Unlock()
		} else {
			conn.Close()
			return
		}
	} else if data.Type == consts.ChangeContent {
		// tw.Tw(context.Background(), "修改文档内容：%s", data.Message)
		res.BlockData.Content = data.Message
	} else if data.Type == consts.ChangeTitle {
		// tw.Tw(context.Background(), "修改文档标题：%s", data.Message)
		res.BlockData.BlockName = data.Message
		md := g.Model("documents")
		ws := r.GetWsForConn(conn)
		_, err = md.Clone().Where("user_id", ws.UserName).Where("block", ws.Block).Update(ws.BlockData)
		if err != nil {
			tw.Tw(context.Background(), "（错误）文档未保存：%s ", err)
		}
	}
}
