package docunit

import (
	"context"
	"encoding/json"
	"fmt"
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
	People    []entity.DocumentTemplate
}

// 连接池
type WebSocketPool struct {
	Connections  map[*websocket.Conn]*WebSocketConnection
	RConnections map[string]*WebSocketConnection // 修改的block ，所有连接用户共同维护同一个block
	Lock         sync.Mutex
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
		Connections:  make(map[*websocket.Conn]*WebSocketConnection),
		RConnections: make(map[string]*WebSocketConnection),
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

// 使用block获取block实例
func (r *DocUnit) GetWsForConn(block string) *WebSocketConnection {
	return r.WebSocketPool.RConnections[block]
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
	// 在连接池中找

	if _, ok := pool.RConnections[block_id]; !ok {
		// 更新文档内容
		md := dao.Documents.Ctx(context.Background())
		var docData entity.Documents
		err := md.Clone().Where("block", block_id).Scan(&docData)
		if err != nil {
			conn.Close()
		}
		wsConn.BlockData = docData
		pool.RConnections[block_id] = wsConn
	}
	pool.Connections[conn] = pool.RConnections[block_id]
	pool.Lock.Unlock()

	go func() {

		// 处理连接消息
		defer func() {
			// 保存文档
			pool.Lock.Lock()
			md := g.Model("documents")
			_, err = md.Clone().Where("block", block_id).Update(pool.RConnections[block_id].BlockData)
			if err != nil {
				tw.Tw(context.Background(), "（错误）文档未保存：%s ", err)
			}

			r.UpdateAndSyncPeople(wsConn.Conn, block_id)
			// 从连接池中移除
			tw.Tw(context.Background(), "文档通信退出：%s", wsConn.UserName)
			// 更新协作
			delete(pool.Connections, wsConn.Conn)
			conn.Close()
			pool.Lock.Unlock()
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

// 更新数组并发送
func (r *DocUnit) UpdateAndSyncPeople(conn *websocket.Conn, document_id string) {
	res := r.GetWsForConn(document_id)
	tableName := fmt.Sprintf("`document-%s`", document_id)
	md := g.Model(tableName)
	var pdata []entity.DocumentTemplate
	err := md.Clone().WithAll().Scan(&pdata)
	if err == nil {
		res.People = pdata
		var NewData = scmsg.SCMsgO{Type: consts.HeartBeatServer,
			Message: "",
			Data: g.Map{"doc_data": res.BlockData,
				"people_data": pdata}}
		NewDataBytes, err := json.Marshal(NewData)
		if err != nil {
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, NewDataBytes)
		if err != nil {
			conn.Close()
		}
	}
}

// 处理message
func (r *DocUnit) HandleWebSocketMessage(conn *websocket.Conn, msg []byte) {
	var pool = r.WebSocketPool
	var data scmsg.SCMsgO
	err := json.Unmarshal(msg, &data)
	res := pool.Connections[conn]
	if err != nil || res == nil {
		conn.Close()
		return
	} else if data.Type == consts.CreateChannel {
		tw.Tw(context.Background(), "创建通道：%s", data.Message)
		r.UpdateAndSyncPeople(conn, res.Block)
	} else if data.Type == consts.ChangeContent {
		// tw.Tw(context.Background(), "修改文档内容：%s", data.Message)
		pool.Lock.Lock()
		// 判断是否有权限更改
		for _, v := range res.People {
			if v.Auth > 0 && v.UserId == res.UserName {
				res.BlockData.Content = data.Message
				break
			}
		}
		pool.Lock.Unlock()
	} else if data.Type == consts.ChangeTitle {
		// tw.Tw(context.Background(), "修改文档标题：%s", data.Message)
		res.BlockData.BlockName = data.Message
		md := g.Model("documents")
		_, err = md.Clone().Where("user_id", res.UserName).Where("block", res.Block).Update(res.BlockData)
		tw.Tw(context.Background(), "（123）文档：%s ", res.Block)
		if err != nil {
			tw.Tw(context.Background(), "（错误）文档未保存：%s ", err)
		}
	}
}
