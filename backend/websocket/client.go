package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入超时
	writeWait = 10 * time.Second

	// 读取超时
	pongWait = 60 * time.Second

	// 发送ping的频率
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

// Client 是WebSocket连接的中间人
type Client struct {
	hub *Hub

	// WebSocket连接
	conn *websocket.Conn

	// 发送消息的缓冲通道
	send chan []byte

	// 玩家ID
	playerID string

	// 玩家名称
	playerName string

	// 房间ID
	roomID string
}

// NewClient 创建一个新的客户端
func NewClient(hub *Hub, conn *websocket.Conn, playerID, playerName string) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		playerID:   playerID,
		playerName: playerName,
	}
}

// ReadPump 从WebSocket连接中泵取消息
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("连接异常关闭: %v", err)
			} else {
				log.Printf("连接正常关闭 (ID: %s)", c.playerID)
			}
			break
		}

		// 解析消息
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Printf("解析消息错误: %v", err)
			continue
		}

		// 处理消息
		c.handleMessage(data)
	}
}

// WritePump 将消息泵送到WebSocket连接
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 直接发送单个消息
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

			// 发送队列中的其他消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				message := <-c.send
				if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(message map[string]interface{}) {
	log.Printf("收到消息: %v", message)
	
	messageType, ok := message["type"].(string)
	if !ok {
		c.sendError("无效的消息类型")
		return
	}

	switch messageType {
	case "join_room":
		roomID, ok := message["room_id"].(string)
		if !ok {
			c.sendError("无效的房间ID")
			return
		}
		c.hub.JoinRoom(c, roomID)

	case "leave_room":
		c.hub.LeaveRoom(c)

	case "game_action":
		c.hub.HandleGameAction(c, message)

	case "create_room":
		log.Printf("创建房间请求")
		// 生成一个唯一的房间ID
		roomID := generateRoomID()
		log.Printf("生成房间ID: %s", roomID)
		
		// 创建并加入房间
		c.hub.CreateRoom(c, roomID)
		
		// 发送房间创建成功的消息
		response := map[string]interface{}{
			"type":    "room_created",
			"room_id": roomID,
		}
		data, _ := json.Marshal(response)
		c.send <- data
		log.Printf("已发送房间创建成功消息")

	default:
		c.sendError("未知的消息类型")
	}
}

// sendError 发送错误消息给客户端
func (c *Client) sendError(message string) {
	response := map[string]interface{}{
		"type":  "error",
		"error": message,
	}
	data, _ := json.Marshal(response)
	c.send <- data
}

// generateRoomID 生成一个唯一的房间ID
func generateRoomID() string {
	// 简单实现，实际应用中可能需要更复杂的逻辑
	return time.Now().Format("20060102150405")
}