package models

import (
	"math/rand"
)

// Suit 表示扑克牌花色
type Suit string

const (
	Hearts   Suit = "hearts"   // 红桃
	Diamonds Suit = "diamonds" // 方块
	Clubs    Suit = "clubs"    // 梅花
	Spades   Suit = "spades"   // 黑桃
)

// Rank 表示扑克牌点数
type Rank int

const (
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
	Two   Rank = 15
	Three Rank = 16
)

// Card 表示一张扑克牌
type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

// IsHeartThree 判断是否为红桃3
func (c *Card) IsHeartThree() bool {
	return c.Suit == Hearts && c.Rank == Three
}

// IsHeartFour 判断是否为红桃4
func (c *Card) IsHeartFour() bool {
	return c.Suit == Hearts && c.Rank == Four
}

// GetRankString 获取牌面的字符串表示
func (c *Card) GetRankString() string {
	switch c.Rank {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	case Two:
		return "2"
	case Three:
		return "3"
	default:
		return string(rune('0' + int(c.Rank)))
	}
}

// Deck 表示一副扑克牌
type Deck []Card

// NewDeck 创建一副新的扑克牌（不含大小王）
func NewDeck() Deck {
	var deck Deck
	suits := []Suit{Hearts, Diamonds, Clubs, Spades}
	ranks := []Rank{Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace, Two, Three}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}

	return deck
}

// Shuffle 洗牌
func (d *Deck) Shuffle() {
	// 使用Fisher-Yates洗牌算法
	cards := *d
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

// Deal 发牌给4名玩家
func (d *Deck) Deal() [4][]Card {
	var hands [4][]Card
	for i := range hands {
		hands[i] = make([]Card, 0, 13)
	}

	cards := *d
	for i, card := range cards {
		playerIndex := i % 4
		hands[playerIndex] = append(hands[playerIndex], card)
	}

	return hands
}