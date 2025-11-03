package game

import (
	"sort"

	"github.com/chenhailong/hong3/models"
)

// CardType 表示牌型
type CardType int

const (
	TypeInvalid    CardType = iota // 无效牌型
	TypeSingle                     // 单张
	TypePair                       // 对子
	TypeStraight                   // 顺子
	TypeConsecutivePairs           // 连对
	TypeThreeOfAKind               // 三张（炸弹）
	TypeFourOfAKind                // 四张（大炸弹）
)

// CardGroup 表示一组出牌
type CardGroup struct {
	Cards []models.Card
	Type  CardType
	Value models.Rank // 用于比较大小的值
}

// ValidateAndCreateCardGroup 验证牌型并创建CardGroup
func ValidateAndCreateCardGroup(cards []models.Card) (*CardGroup, bool) {
	if len(cards) == 0 {
		return nil, false
	}

	// 按点数排序
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank < cards[j].Rank
	})

	// 检查是否包含红桃3（不能作为顺子或连对的一部分）
	for _, card := range cards {
		if card.IsHeartThree() && (len(cards) >= 5 || (len(cards) >= 6 && len(cards)%2 == 0)) {
			return nil, false
		}
	}

	// 单张
	if len(cards) == 1 {
		return &CardGroup{
			Cards: cards,
			Type:  TypeSingle,
			Value: cards[0].Rank,
		}, true
	}

	// 对子
	if len(cards) == 2 && cards[0].Rank == cards[1].Rank {
		return &CardGroup{
			Cards: cards,
			Type:  TypePair,
			Value: cards[0].Rank,
		}, true
	}

	// 三张（炸弹）
	if len(cards) == 3 && cards[0].Rank == cards[1].Rank && cards[1].Rank == cards[2].Rank {
		return &CardGroup{
			Cards: cards,
			Type:  TypeThreeOfAKind,
			Value: cards[0].Rank,
		}, true
	}

	// 四张（大炸弹）
	if len(cards) == 4 && cards[0].Rank == cards[1].Rank && cards[1].Rank == cards[2].Rank && cards[2].Rank == cards[3].Rank {
		return &CardGroup{
			Cards: cards,
			Type:  TypeFourOfAKind,
			Value: cards[0].Rank,
		}, true
	}

	// 顺子（至少5张连续单牌）
	if len(cards) >= 5 {
		isStraight := true
		for i := 1; i < len(cards); i++ {
			if cards[i].Rank != cards[i-1].Rank+1 || cards[i].Rank == models.Three {
				isStraight = false
				break
			}
		}
		if isStraight {
			return &CardGroup{
				Cards: cards,
				Type:  TypeStraight,
				Value: cards[len(cards)-1].Rank,
			}, true
		}
	}

	// 连对（至少3对连续对子）
	if len(cards) >= 6 && len(cards)%2 == 0 {
		isConsecutivePairs := true
		for i := 0; i < len(cards); i += 2 {
			// 检查是否为对子
			if cards[i].Rank != cards[i+1].Rank {
				isConsecutivePairs = false
				break
			}
			// 检查连续性
			if i >= 2 && cards[i].Rank != cards[i-2].Rank+1 {
				isConsecutivePairs = false
				break
			}
			// 不能包含3
			if cards[i].Rank == models.Three {
				isConsecutivePairs = false
				break
			}
		}
		if isConsecutivePairs {
			return &CardGroup{
				Cards: cards,
				Type:  TypeConsecutivePairs,
				Value: cards[len(cards)-1].Rank,
			}, true
		}
	}

	return nil, false
}

// CanBeat 判断当前牌组是否能够打过另一个牌组
func (cg *CardGroup) CanBeat(other *CardGroup) bool {
	// 大炸弹可以打过任何牌
	if cg.Type == TypeFourOfAKind {
		if other.Type == TypeFourOfAKind {
			return cg.Value > other.Value
		}
		return true
	}

	// 炸弹可以打过非炸弹牌型
	if cg.Type == TypeThreeOfAKind && other.Type != TypeFourOfAKind {
		if other.Type == TypeThreeOfAKind {
			return cg.Value > other.Value
		}
		return true
	}

	// 相同牌型比较大小
	if cg.Type == other.Type {
		// 对于顺子和连对，长度必须相同
		if (cg.Type == TypeStraight || cg.Type == TypeConsecutivePairs) && len(cg.Cards) != len(other.Cards) {
			return false
		}
		return cg.Value > other.Value
	}

	return false
}