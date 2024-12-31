package docunit

import (
	"context"
	"encoding/json"
	"gf_chat_server/internal/consts"
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
	Conn     *websocket.Conn
	User     entity.User
	UserName string
	Content  string
}

// 连接池
type WebSocketPool struct {
	Connections map[*WebSocketConnection]bool
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
type WebSocketUnit struct {
	WebSocketPool *WebSocketPool
}

func New() WebSocketUnit {
	pool := &WebSocketPool{
		Connections: make(map[*WebSocketConnection]bool),
	}
	return WebSocketUnit{WebSocketPool: pool}
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
		Content:  "",
	}
	// 将连接添加到连接池
	pool.Lock.Lock()
	pool.Connections[wsConn] = true
	pool.Lock.Unlock()

	// 处理连接消息
	defer func() {
		// 连接关闭时从连接池中移除
		pool.Lock.Lock()
		tw.Tw(context.Background(), "文档通信退出：%s", wsConn.UserName)
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
		r.HandleWebSocketMessage(conn, message)
	}
}

// 处理message
func (r *WebSocketUnit) HandleWebSocketMessage(conn *websocket.Conn, msg []byte) {
	var data scmsg.SCMsgO
	err := json.Unmarshal(msg, &data)
	res := r.GetWsForConn(conn)

	if err != nil || res == nil {
		return
	} else if data.Type == consts.CreateChannel {
		tw.Tw(context.Background(), "创建通道：%s", data.Message)
	} else if data.Type == consts.ChangeContent {
		tw.Tw(context.Background(), "修改文档内容：%s", data.Message)
	}
}
