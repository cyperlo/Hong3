package game

import (
	"errors"
	"fmt"
	"sync"

	"github.com/chenhailong/hong3/models"
)

// PlayerStatus 表示玩家状态
type PlayerStatus int

const (
	PlayerStatusWaiting  PlayerStatus = iota // 等待中
	PlayerStatusReady                        // 已准备
	PlayerStatusPlaying                      // 游戏中
	PlayerStatusFinished                     // 已完成
)

// GameStatus 表示游戏状态
type GameStatus int

const (
	GameStatusWaiting  GameStatus = iota // 等待中
	GameStatusPlaying                    // 游戏中
	GameStatusFinished                   // 已结束
)

// TeamType 表示队伍类型
type TeamType int

const (
	TeamTypeUnknown TeamType = iota // 未知
	TeamType1v3                     // 1v3
	TeamType2v2                     // 2v2
)

// Player 表示游戏中的玩家
type Player struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Status   PlayerStatus   `json:"status"`
	Cards    []models.Card  `json:"cards,omitempty"` // 手牌
	Team     int            `json:"-"`               // 队伍编号（1或2）
	Position int            `json:"position"`        // 座位位置（0-3）
	CardCount int           `json:"card_count"`      // 手牌数量
	CollectedCards int      `json:"collected_cards"` // 收集的牌数
}

// Game 表示一局游戏
type Game struct {
	ID            string                `json:"id"`
	Status        GameStatus            `json:"status"`
	Players       [4]*Player            `json:"players"`
	CurrentPlayer int                   `json:"current_player"` // 当前出牌玩家索引
	LastPlayer    int                   `json:"last_player"`    // 最后出牌玩家索引
	TableCards    *CardGroup            `json:"table_cards"`    // 桌面上的牌
	TeamType      TeamType              `json:"-"`              // 队伍类型
	FinishedOrder []int                 `json:"finished_order"` // 完成顺序
	mutex         sync.Mutex
}

// NewGame 创建新游戏
func NewGame(id string) *Game {
	return &Game{
		ID:            id,
		Status:        GameStatusWaiting,
		CurrentPlayer: -1,
		LastPlayer:    -1,
		FinishedOrder: make([]int, 0, 4),
	}
}

// AddPlayer 添加玩家到游戏
func (g *Game) AddPlayer(player *Player) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Status != GameStatusWaiting {
		return errors.New("游戏已开始，无法加入")
	}

	// 检查玩家是否已经存在
	for _, p := range g.Players {
		if p != nil && p.ID == player.ID {
			// 玩家已存在，更新信息
			p.Name = player.Name
			p.Status = player.Status
			return nil
		}
	}

	// 查找空位
	for i, p := range g.Players {
		if p == nil {
			player.Position = i
			g.Players[i] = player
			return nil
		}
	}

	return errors.New("房间已满")
}

// RemovePlayer 从游戏中移除玩家
func (g *Game) RemovePlayer(playerID string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for i, p := range g.Players {
		if p != nil && p.ID == playerID {
			g.Players[i] = nil
			break
		}
	}
}

// SetPlayerReady 设置玩家准备状态
func (g *Game) SetPlayerReady(playerID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Status != GameStatusWaiting {
		return errors.New("游戏已开始")
	}

	for _, p := range g.Players {
		if p != nil && p.ID == playerID {
			p.Status = PlayerStatusReady
			return nil
		}
	}

	return errors.New("玩家不在游戏中")
}

// AllPlayersReady 检查是否所有玩家都已准备
func (g *Game) AllPlayersReady() bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	playerCount := 0
	readyCount := 0

	for _, p := range g.Players {
		if p != nil {
			playerCount++
			if p.Status == PlayerStatusReady {
				readyCount++
			}
		}
	}

	return playerCount == 4 && readyCount == 4
}

// StartGame 开始游戏
func (g *Game) StartGame() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Status != GameStatusWaiting {
		return errors.New("游戏已开始")
	}

	// 检查所有玩家是否准备好（不需要再次加锁，因为已经在锁内）
	playerCount := 0
	readyCount := 0
	for _, p := range g.Players {
		if p != nil {
			playerCount++
			if p.Status == PlayerStatusReady {
				readyCount++
			}
		}
	}
	
	if playerCount != 4 || readyCount != 4 {
		return fmt.Errorf("玩家未全部准备: 玩家数=%d, 已准备数=%d", playerCount, readyCount)
	}

	// 创建并洗牌
	deck := models.NewDeck()
	deck.Shuffle()

	// 发牌
	hands := deck.Deal()
	for i, p := range g.Players {
		if p != nil {
			p.Cards = hands[i]
			p.CardCount = len(p.Cards)
			p.Status = PlayerStatusPlaying
		}
	}

	// 确定队伍
	g.determineTeams()

	// 找到持有红桃4的玩家作为首个出牌者
	g.findFirstPlayer()

	g.Status = GameStatusPlaying
	return nil
}

// determineTeams 根据红桃3的分布确定队伍
func (g *Game) determineTeams() {
	heartThreeCount := make(map[int]int)
	
	// 统计每个玩家手中红桃3的数量
	for i, p := range g.Players {
		if p == nil {
			continue
		}
		
		for _, card := range p.Cards {
			if card.IsHeartThree() {
				heartThreeCount[i]++
			}
		}
	}
	
	// 判断队伍类型
	if len(heartThreeCount) == 1 {
		// 一个玩家有两张红桃3，1v3模式
		g.TeamType = TeamType1v3
		for i := range g.Players {
			if heartThreeCount[i] == 2 {
				g.Players[i].Team = 1 // 单人队伍
			} else {
				g.Players[i].Team = 2 // 三人队伍
			}
		}
	} else {
		// 两个玩家各有一张红桃3，2v2模式
		g.TeamType = TeamType2v2
		team := 1
		for i := range g.Players {
			if heartThreeCount[i] > 0 {
				g.Players[i].Team = team
				team = 3 - team // 切换队伍（1->2或2->1）
			}
		}
		
		// 为没有红桃3的玩家分配队伍
		for i := range g.Players {
			if heartThreeCount[i] == 0 {
				// 找到同队友
				for j := range g.Players {
					if i != j && heartThreeCount[j] > 0 && g.Players[j].Team != 0 {
						g.Players[i].Team = g.Players[j].Team
						break
					}
				}
			}
		}
	}
}

// findFirstPlayer 找到持有红桃4的玩家作为首个出牌者
func (g *Game) findFirstPlayer() {
	for i, p := range g.Players {
		if p == nil {
			continue
		}
		
		// 检查是否有四个4（特殊规则）
		fourCount := 0
		for _, card := range p.Cards {
			if card.Rank == models.Four {
				fourCount++
			}
		}
		
		if fourCount == 4 {
			g.CurrentPlayer = i
			g.LastPlayer = i
			return
		}
		
		// 检查是否有红桃4
		for _, card := range p.Cards {
			if card.IsHeartFour() {
				g.CurrentPlayer = i
				g.LastPlayer = i
				return
			}
		}
	}
}

// PlayCards 玩家出牌
func (g *Game) PlayCards(playerID string, cardIndices []int) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Status != GameStatusPlaying {
		return errors.New("游戏未开始或已结束")
	}

	// 找到玩家
	var playerIndex int = -1
	var player *Player
	for i, p := range g.Players {
		if p != nil && p.ID == playerID {
			playerIndex = i
			player = p
			break
		}
	}

	if playerIndex == -1 || player == nil {
		return errors.New("玩家不在游戏中")
	}

	if playerIndex != g.CurrentPlayer {
		return errors.New("不是该玩家的回合")
	}

	// 检查玩家是否已经完成游戏
	if player.Status == PlayerStatusFinished {
		return errors.New("玩家已经出完牌")
	}

	// 如果不出牌（过）
	if len(cardIndices) == 0 {
		// 如果是最后出牌的玩家，不能过
		if playerIndex == g.LastPlayer {
			return errors.New("最后出牌的玩家不能过")
		}
		
		// 更新当前玩家
		g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
		// 跳过已经完成的玩家
		for g.Players[g.CurrentPlayer] == nil || g.Players[g.CurrentPlayer].Status == PlayerStatusFinished {
			g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
		}
		
		return nil
	}

	// 检查索引是否有效
	if len(cardIndices) > len(player.Cards) {
		return errors.New("选择的牌数超过手牌数量")
	}

	// 获取选中的牌
	selectedCards := make([]models.Card, 0, len(cardIndices))
	for _, idx := range cardIndices {
		if idx < 0 || idx >= len(player.Cards) {
			return errors.New("无效的牌索引")
		}
		selectedCards = append(selectedCards, player.Cards[idx])
	}

	// 验证牌型
	cardGroup, valid := ValidateAndCreateCardGroup(selectedCards)
	if !valid {
		return errors.New("无效的牌型")
	}

	// 首轮必须出红桃4
	if g.TableCards == nil {
		hasHeartFour := false
		for _, card := range selectedCards {
			if card.IsHeartFour() {
				hasHeartFour = true
				break
			}
		}

		// 检查是否有四个4的特殊规则
		fourCount := 0
		for _, card := range player.Cards {
			if card.Rank == models.Four {
				fourCount++
			}
		}

		if fourCount < 4 && !hasHeartFour {
			return errors.New("首轮必须出包含红桃4的牌")
		}
	} else {
		// 检查是否能打过桌面上的牌
		if !cardGroup.CanBeat(g.TableCards) {
			return errors.New("出的牌无法打过桌面上的牌")
		}
	}

	// 更新桌面牌
	g.TableCards = cardGroup

	// 从玩家手牌中移除出的牌
	for i := len(cardIndices) - 1; i >= 0; i-- {
		idx := cardIndices[i]
		player.Cards = append(player.Cards[:idx], player.Cards[idx+1:]...)
	}
	player.CardCount = len(player.Cards)

	// 更新最后出牌玩家
	g.LastPlayer = playerIndex

	// 检查玩家是否出完牌
	if len(player.Cards) == 0 {
		player.Status = PlayerStatusFinished
		g.FinishedOrder = append(g.FinishedOrder, playerIndex)
		
		// 检查游戏是否结束
		if g.checkGameEnd() {
			g.Status = GameStatusFinished
			return nil
		}
	}

	// 更新当前玩家
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	// 跳过已经完成的玩家
	for g.Players[g.CurrentPlayer] == nil || g.Players[g.CurrentPlayer].Status == PlayerStatusFinished {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	}

	return nil
}

// Pass 玩家选择不出牌（过）
func (g *Game) Pass(playerID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Status != GameStatusPlaying {
		return errors.New("游戏未开始或已结束")
	}

	// 找到玩家
	var playerIndex int = -1
	for i, p := range g.Players {
		if p != nil && p.ID == playerID {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return errors.New("玩家不在游戏中")
	}

	if playerIndex != g.CurrentPlayer {
		return errors.New("不是该玩家的回合")
	}

	// 如果是最后出牌的玩家，不能过
	if playerIndex == g.LastPlayer {
		return errors.New("最后出牌的玩家不能过")
	}

	// 更新当前玩家
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	// 跳过已经完成的玩家
	for g.Players[g.CurrentPlayer] == nil || g.Players[g.CurrentPlayer].Status == PlayerStatusFinished {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	}

	// 检查是否所有人都过了
	passCount := 0
	for i := 0; i < 4; i++ {
		if i != g.LastPlayer && (g.Players[i] == nil || g.Players[i].Status == PlayerStatusFinished) {
			passCount++
		}
	}

	// 如果所有其他玩家都过了，最后出牌的玩家获得本轮牌权
	if passCount == 3 {
		// 清空桌面牌
		g.TableCards = nil
		// 最后出牌的玩家重新开始
		g.CurrentPlayer = g.LastPlayer
	}

	return nil
}

// checkGameEnd 检查游戏是否结束
func (g *Game) checkGameEnd() bool {
	if g.TeamType == TeamType1v3 {
		// 1v3模式，任何一方出完牌即结束
		return len(g.FinishedOrder) > 0
	} else if g.TeamType == TeamType2v2 {
		// 2v2模式，检查是否有一个队伍的两名玩家都出完牌
		teamFinished := make(map[int]int)
		for _, idx := range g.FinishedOrder {
			team := g.Players[idx].Team
			teamFinished[team]++
			if teamFinished[team] == 2 {
				return true
			}
		}
		
		// 或者检查是否只剩下一名玩家
		activePlayers := 0
		for _, p := range g.Players {
			if p != nil && p.Status == PlayerStatusPlaying {
				activePlayers++
			}
		}
		return activePlayers <= 1
	}
	
	return false
}

// CollectTableCards 收集桌面上的牌
func (g *Game) CollectTableCards() {
	if g.TableCards == nil || g.LastPlayer < 0 {
		return
	}
	
	player := g.Players[g.LastPlayer]
	if player != nil {
		player.CollectedCards += len(g.TableCards.Cards)
	}
	
	g.TableCards = nil
}

// GetGameResult 获取游戏结果
func (g *Game) GetGameResult() map[string]interface{} {
	result := make(map[string]interface{})
	
	if g.Status != GameStatusFinished {
		result["status"] = "游戏未结束"
		return result
	}
	
	// 计算胜利队伍
	var winningTeam int
	if g.TeamType == TeamType1v3 {
		// 1v3模式，看谁先出完牌
		winningTeam = g.Players[g.FinishedOrder[0]].Team
	} else {
		// 2v2模式
		teamFinished := make(map[int]int)
		for _, idx := range g.FinishedOrder {
			team := g.Players[idx].Team
			teamFinished[team]++
			if teamFinished[team] == 2 {
				winningTeam = team
				break
			}
		}
		
		// 如果没有队伍两人都出完牌，比较剩余牌数
		if winningTeam == 0 {
			team1Cards := 0
			team2Cards := 0
			for _, p := range g.Players {
				if p != nil {
					if p.Team == 1 {
						team1Cards += p.CollectedCards
					} else {
						team2Cards += p.CollectedCards
					}
				}
			}
			
			if team1Cards > team2Cards {
				winningTeam = 1
			} else {
				winningTeam = 2
			}
		}
	}
	
	result["winning_team"] = winningTeam
	result["finished_order"] = g.FinishedOrder
	
	playerResults := make([]map[string]interface{}, 0, 4)
	for _, p := range g.Players {
		if p != nil {
			playerResult := map[string]interface{}{
				"id":             p.ID,
				"name":           p.Name,
				"position":       p.Position,
				"team":           p.Team,
				"collected_cards": p.CollectedCards,
				"is_winner":      p.Team == winningTeam,
			}
			playerResults = append(playerResults, playerResult)
		}
	}
	result["players"] = playerResults
	
	return result
}

// GetStatus 获取游戏状态的字符串表示
func (g *Game) GetStatus() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	switch g.Status {
	case GameStatusWaiting:
		return "waiting"
	case GameStatusPlaying:
		return "playing"
	case GameStatusFinished:
		return "finished"
	default:
		return "unknown"
	}
}