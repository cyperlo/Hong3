package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chenhailong/hong3/game"
)

// Hub 维护活跃的客户端连接和广播消息
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 房间映射
	rooms map[string]map[*Client]bool

	// 游戏映射
	games map[string]*game.Game

	// 广播消息通道
	broadcast chan []byte

	// 注册客户端的通道
	Register chan *Client

	// 注销客户端的通道
	Unregister chan *Client

	// 互斥锁
	mutex sync.Mutex
}

// NewHub 创建一个新的Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		games:      make(map[string]*game.Game),
	}
}

// Run 启动Hub的消息处理循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// 从房间中移除
				if client.roomID != "" {
					if room, ok := h.rooms[client.roomID]; ok {
						delete(room, client)
						if len(room) == 0 {
							delete(h.rooms, client.roomID)
							delete(h.games, client.roomID)
						} else {
							// 通知房间内其他玩家
							h.broadcastToRoom(client.roomID, map[string]interface{}{
								"type":     "player_left",
								"playerID": client.playerID,
							})
						}
					}
				}
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// JoinRoom 将客户端加入房间
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// 如果房间不存在，创建它
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
		h.games[roomID] = game.NewGame(roomID)
	}

	// 使用内部加入逻辑
	h.joinRoomInternal(client, roomID)
}

// sendRoomStateToClient 向客户端发送房间状态
func (h *Hub) sendRoomStateToClient(client *Client, roomID string) {
	room, ok := h.rooms[roomID]
	if !ok {
		return
	}

	// 获取房间内所有玩家信息
	players := make([]map[string]interface{}, 0, len(room))
	g := h.games[roomID]
	for c := range room {
		playerInfo := map[string]interface{}{
			"id":   c.playerID,
			"name": c.playerName,
		}
		// 从游戏中获取玩家的准备状态
		if g != nil {
			for _, p := range g.Players {
				if p != nil && p.ID == c.playerID {
					playerInfo["ready"] = p.Status == game.PlayerStatusReady
					break
				}
			}
		}
		players = append(players, playerInfo)
	}

	// 准备房间状态消息
	roomState := map[string]interface{}{
		"type":    "room_state",
		"room_id": roomID,
		"players": players,
	}

	// 发送房间状态
	data, err := json.Marshal(roomState)
	if err != nil {
		log.Printf("Error marshalling room state: %v", err)
		return
	}

	log.Printf("发送房间状态给玩家 %s: %s", client.playerID, string(data))
	client.send <- data
}

// broadcastToRoomExcept 向房间内除指定客户端外的所有客户端广播消息
func (h *Hub) broadcastToRoomExcept(roomID string, except *Client, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	if room, ok := h.rooms[roomID]; ok {
		for client := range room {
			if client != except {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(room, client)
					delete(h.clients, client)
				}
			}
		}
	}
}

// LeaveRoom 将客户端从房间中移除
func (h *Hub) LeaveRoom(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if client.roomID == "" {
		return
	}

	roomID := client.roomID
	if room, ok := h.rooms[roomID]; ok {
		delete(room, client)
		client.roomID = ""

		// 如果房间为空，删除它
		if len(room) == 0 {
			delete(h.rooms, roomID)
			delete(h.games, roomID)
		} else {
			// 更新其他玩家的房间状态
			for c := range room {
				h.sendRoomStateToClient(c, roomID)
			}
			// 通知房间内其他玩家
			h.broadcastToRoom(roomID, map[string]interface{}{
				"type":     "player_left",
				"playerID": client.playerID,
			})
		}
	}
}

// broadcastToRoom 向房间内所有客户端广播消息
func (h *Hub) broadcastToRoom(roomID string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	if room, ok := h.rooms[roomID]; ok {
		for client := range room {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(room, client)
				delete(h.clients, client)
			}
		}
	}
}

// GetGame 获取房间对应的游戏
func (h *Hub) GetGame(roomID string) *game.Game {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.games[roomID]
}

// HandleGameAction 处理游戏相关的动作
func (h *Hub) HandleGameAction(client *Client, action map[string]interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	roomID := client.roomID
	if roomID == "" {
		client.sendError("未加入房间")
		return
	}

	g, ok := h.games[roomID]
	if !ok {
		client.sendError("游戏不存在")
		return
	}

	actionType, ok := action["action"].(string)
	if !ok {
		client.sendError("无效的动作类型")
		return
	}

	switch actionType {
	case "ready":
		log.Printf("玩家 %s (ID: %s) 准备，房间: %s", client.playerName, client.playerID, roomID)
		
		// 确保玩家在游戏中
		playerInGame := false
		for i, p := range g.Players {
			if p != nil && p.ID == client.playerID {
				playerInGame = true
				log.Printf("玩家 %s 已在游戏中，位置: %d", client.playerID, i)
				break
			}
		}
		
		// 如果玩家不在游戏中，尝试添加
		if !playerInGame {
			log.Printf("玩家 %s 不在游戏中，尝试添加", client.playerID)
			player := &game.Player{
				ID:       client.playerID,
				Name:     client.playerName,
				Status:   game.PlayerStatusWaiting,
				CardCount: 0,
			}
			if err := g.AddPlayer(player); err != nil {
				log.Printf("准备时添加玩家失败: 玩家ID=%s, 错误=%v", client.playerID, err)
				// 打印当前游戏中的所有玩家
				log.Printf("当前游戏中的玩家:")
				for i, p := range g.Players {
					if p != nil {
						log.Printf("  位置 %d: ID=%s, Name=%s", i, p.ID, p.Name)
					} else {
						log.Printf("  位置 %d: 空", i)
					}
				}
				client.sendError(fmt.Sprintf("玩家不在游戏中，无法准备: %v", err))
				return
			}
			log.Printf("成功添加玩家 %s 到游戏", client.playerID)
		}
		
		// 再次确认玩家在游戏中
		err := g.SetPlayerReady(client.playerID)
		if err != nil {
			log.Printf("设置玩家准备状态失败: 玩家ID=%s, 错误=%v", client.playerID, err)
			// 打印当前游戏中的所有玩家
			log.Printf("当前游戏中的玩家:")
			for i, p := range g.Players {
				if p != nil {
					log.Printf("  位置 %d: ID=%s, Name=%s", i, p.ID, p.Name)
				} else {
					log.Printf("  位置 %d: 空", i)
				}
			}
			client.sendError(err.Error())
			return
		}
		
		log.Printf("玩家 %s 准备成功", client.playerID)

		// 广播玩家准备状态
		h.broadcastToRoom(roomID, map[string]interface{}{
			"type":     "player_ready",
			"playerID": client.playerID,
		})

		// 更新所有玩家的房间状态（包含准备状态）
		if room, ok := h.rooms[roomID]; ok {
			for c := range room {
				h.sendRoomStateToClient(c, roomID)
			}
		}

		// 检查是否所有玩家都准备好了
		log.Printf("检查所有玩家准备状态...")
		log.Printf("当前游戏中的玩家:")
		for i, p := range g.Players {
			if p != nil {
				log.Printf("  位置 %d: ID=%s, Name=%s, Status=%d (Ready=%v)", i, p.ID, p.Name, p.Status, p.Status == game.PlayerStatusReady)
			} else {
				log.Printf("  位置 %d: 空", i)
			}
		}
		
		allReady := g.AllPlayersReady()
		log.Printf("所有玩家准备状态: %v", allReady)
		
		if allReady {
			log.Printf("所有玩家已准备，开始游戏...")
			err := g.StartGame()
			if err != nil {
				log.Printf("开始游戏失败: %v", err)
				client.sendError(err.Error())
				return
			}
			log.Printf("游戏开始成功！")

			// 先广播游戏开始
			h.broadcastToRoom(roomID, map[string]interface{}{
				"type":          "game_started",
				"currentPlayer": g.CurrentPlayer,
			})

			// 然后向每个玩家发送他们的手牌（延迟一小段时间，确保 game_started 消息先到达）
			go func() {
				time.Sleep(100 * time.Millisecond)
				h.mutex.Lock()
				defer h.mutex.Unlock()
				if room, ok := h.rooms[roomID]; ok {
					for clientInRoom := range room {
						for _, player := range g.Players {
							if player != nil && player.ID == clientInRoom.playerID {
								select {
								case clientInRoom.send <- h.createGameStateMessage(g, player.ID):
								default:
									log.Printf("无法发送游戏状态给玩家 %s", clientInRoom.playerID)
								}
								break
							}
						}
					}
				}
			}()
		}

	case "play_cards":
		cardIndices, ok := action["card_indices"].([]interface{})
		if !ok {
			client.sendError("无效的牌索引")
			return
		}

		// 转换为int切片
		indices := make([]int, len(cardIndices))
		for i, idx := range cardIndices {
			if idxFloat, ok := idx.(float64); ok {
				indices[i] = int(idxFloat)
			} else {
				client.sendError("无效的牌索引格式")
				return
			}
		}

		err := g.PlayCards(client.playerID, indices)
		if err != nil {
			client.sendError(err.Error())
			return
		}

		// 向每个玩家发送更新后的游戏状态
		for clientInRoom := range h.rooms[roomID] {
			for _, player := range g.Players {
				if player != nil && player.ID == clientInRoom.playerID {
					clientInRoom.send <- h.createGameStateMessage(g, player.ID)
					break
				}
			}
		}

		// 检查游戏是否结束
		if g.Status == game.GameStatusFinished {
			result := g.GetGameResult()
			h.broadcastToRoom(roomID, map[string]interface{}{
				"type":   "game_end",
				"result": result,
			})
		}

	case "pass":
		err := g.Pass(client.playerID)
		if err != nil {
			client.sendError(err.Error())
			return
		}

		// 广播玩家过牌
		h.broadcastToRoom(roomID, map[string]interface{}{
			"type":          "player_pass",
			"playerID":      client.playerID,
			"currentPlayer": g.CurrentPlayer,
		})

		// 如果所有人都过了，清空桌面牌
		if g.TableCards == nil {
			h.broadcastToRoom(roomID, map[string]interface{}{
				"type":          "round_end",
				"currentPlayer": g.CurrentPlayer,
			})
		}
	}
}

// createGameStateMessage 创建游戏状态消息
func (h *Hub) createGameStateMessage(g *game.Game, playerID string) []byte {
	log.Printf("创建游戏状态消息给玩家 %s", playerID)
	// 找到玩家
	var currentPlayer *game.Player
	for _, p := range g.Players {
		if p != nil && p.ID == playerID {
			currentPlayer = p
			break
		}
	}

	if currentPlayer == nil {
		return []byte("{}")
	}

	// 创建其他玩家信息（不包含手牌）
	otherPlayers := make([]map[string]interface{}, 0, 3)
	for _, p := range g.Players {
		if p != nil && p.ID != playerID {
			otherPlayers = append(otherPlayers, map[string]interface{}{
				"id":             p.ID,
				"name":           p.Name,
				"position":       p.Position,
				"status":         p.Status,
				"card_count":     p.CardCount,
				"collected_cards": p.CollectedCards,
			})
		}
	}

	// 创建游戏状态
	gameState := map[string]interface{}{
		"type":          "game_state",
		"status":        g.Status,
		"current_player": g.CurrentPlayer,
		"last_player":   g.LastPlayer,
		"player":        currentPlayer,
		"other_players": otherPlayers,
		"table_cards":   g.TableCards,
	}

	data, _ := json.Marshal(gameState)
	return data
}

// GetRooms 获取所有房间信息
func (h *Hub) GetRooms() []map[string]interface{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	rooms := make([]map[string]interface{}, 0, len(h.rooms))
	for roomID, clients := range h.rooms {
		// 获取房间内的玩家信息
		players := make([]map[string]string, 0, len(clients))
		for client := range clients {
			players = append(players, map[string]string{
				"id":   client.playerID,
				"name": client.playerName,
			})
		}

		// 获取游戏状态
		game := h.games[roomID]
		var gameStatus string
		if game != nil {
			gameStatus = game.GetStatus()
		} else {
			gameStatus = "waiting"
		}

		// 构建房间信息
		roomInfo := map[string]interface{}{
			"id":       roomID,
			"players":  players,
			"status":   gameStatus,
			"capacity": 4, // 每个房间最多4个玩家
		}
		rooms = append(rooms, roomInfo)
	}

	return rooms
}

// CreateRoom 创建新房间并将客户端加入
func (h *Hub) CreateRoom(client *Client, roomID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	log.Printf("创建房间: %s", roomID)

	// 如果房间已存在，直接加入
	if _, ok := h.rooms[roomID]; ok {
		// 房间已存在，使用 JoinRoom 的逻辑
		h.joinRoomInternal(client, roomID)
		return
	}

	// 创建新房间
	h.rooms[roomID] = make(map[*Client]bool)
	h.games[roomID] = game.NewGame(roomID)

	// 将客户端加入房间
	h.rooms[roomID][client] = true
	client.roomID = roomID

	// 将玩家添加到游戏中
	g := h.games[roomID]
	if g != nil {
		player := &game.Player{
			ID:       client.playerID,
			Name:     client.playerName,
			Status:   game.PlayerStatusWaiting,
			CardCount: 0,
		}
		if err := g.AddPlayer(player); err != nil {
			log.Printf("创建房间时添加玩家到游戏失败: 玩家ID=%s, 房间=%s, 错误=%v", client.playerID, roomID, err)
		} else {
			log.Printf("成功添加玩家 %s (%s) 到游戏 %s", client.playerName, client.playerID, roomID)
		}
	} else {
		log.Printf("警告: 房间 %s 的游戏对象为 nil", roomID)
	}

	// 发送房间状态
	h.sendRoomStateToClient(client, roomID)
}

// joinRoomInternal 内部加入房间逻辑（不加锁，需在持有锁的情况下调用）
func (h *Hub) joinRoomInternal(client *Client, roomID string) {
	// 检查房间是否已满（最多4人）
	if len(h.rooms[roomID]) >= 4 {
		client.sendError("房间已满")
		return
	}

	// 检查玩家是否已经在房间中
	if _, exists := h.rooms[roomID][client]; exists {
		// 玩家已在房间中，确保玩家在游戏中
		g := h.games[roomID]
		if g != nil {
			// 检查玩家是否在游戏中
			found := false
			for _, p := range g.Players {
				if p != nil && p.ID == client.playerID {
					found = true
					break
				}
			}
			// 如果不在游戏中，添加玩家
			if !found {
				log.Printf("玩家 %s 在房间中但不在游戏中，尝试添加", client.playerID)
				player := &game.Player{
					ID:       client.playerID,
					Name:     client.playerName,
					Status:   game.PlayerStatusWaiting,
					CardCount: 0,
				}
				if err := g.AddPlayer(player); err != nil {
					log.Printf("加入房间时添加玩家到游戏失败: 玩家ID=%s, 房间=%s, 错误=%v", client.playerID, roomID, err)
				} else {
					log.Printf("成功添加玩家 %s (%s) 到游戏 %s", client.playerName, client.playerID, roomID)
				}
			}
		}
		// 发送房间状态
		h.sendRoomStateToClient(client, roomID)
		return
	}

	// 将客户端加入房间
	h.rooms[roomID][client] = true
	client.roomID = roomID

	// 将玩家添加到游戏中
	g := h.games[roomID]
	if g != nil {
		player := &game.Player{
			ID:       client.playerID,
			Name:     client.playerName,
			Status:   game.PlayerStatusWaiting,
			CardCount: 0,
		}
		if err := g.AddPlayer(player); err != nil {
			log.Printf("加入房间时添加玩家到游戏失败: 玩家ID=%s, 房间=%s, 错误=%v", client.playerID, roomID, err)
		} else {
			log.Printf("成功添加玩家 %s (%s) 到游戏 %s", client.playerName, client.playerID, roomID)
		}
	} else {
		log.Printf("警告: 房间 %s 的游戏对象为 nil，无法添加玩家", roomID)
	}

	// 向新加入的玩家发送完整的房间状态
	h.sendRoomStateToClient(client, roomID)

	// 通知房间内其他玩家有新玩家加入
	h.broadcastToRoomExcept(roomID, client, map[string]interface{}{
		"type":     "player_joined",
		"playerID": client.playerID,
		"name":     client.playerName,
	})
}