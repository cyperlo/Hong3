<script setup>
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue';
import PlayingCard from '../components/PlayingCard.vue';
import gameStore from '../store/gameStore.js';
import authStore from '../store/authStore.js';

const props = defineProps({
  roomId: String,
  playerId: String,
  playerName: String
});

const emit = defineEmits(['game-end']);
const { state: gameState } = gameStore;

const isMyTurn = computed(() => {
  return gameState.currentPlayer === gameState.player?.position;
});


// 发牌动画相关
const isDealing = ref(false);
const dealtCards = ref([]);
const showCards = ref(false);

// 对手牌进行排序（从大到小：3 > 2 > A > K > ... > 4）
const sortCards = (cards) => {
  if (!cards || cards.length === 0) return [];
  
  const sorted = [...cards].sort((a, b) => {
    // 先按点数排序（从大到小）
    if (a.rank !== b.rank) {
      return b.rank - a.rank;
    }
    // 点数相同，按花色排序（红桃 > 方块 > 黑桃 > 梅花）
    const suitOrder = { 'hearts': 4, 'diamonds': 3, 'spades': 2, 'clubs': 1 };
    return suitOrder[b.suit] - suitOrder[a.suit];
  });
  
  return sorted;
};

// 计算排序后的手牌
const sortedCards = computed(() => {
  if (!gameState.player || !gameState.player.cards) {
    return [];
  }
  return sortCards(gameState.player.cards);
});

// 监听手牌变化，触发发牌动画
watch(() => gameState.player?.cards, (newCards, oldCards) => {
  // 如果是从无到有（初始发牌），触发发牌动画
  if (newCards && newCards.length > 0 && (!oldCards || oldCards.length === 0)) {
    startDealingAnimation(newCards);
  } else if (newCards && newCards.length > 0) {
    // 手牌已存在
    if (oldCards && oldCards.length !== newCards.length) {
      // 手牌数量变化（出牌后），直接显示，不需要动画
      showCards.value = true;
      dealtCards.value = sortCards(newCards);
      isDealing.value = false;
    } else if (!showCards.value) {
      // 第一次显示
      showCards.value = true;
      dealtCards.value = sortCards(newCards);
    }
  } else if (!newCards || newCards.length === 0) {
    // 手牌被清空
    showCards.value = false;
    dealtCards.value = [];
    isDealing.value = false;
  }
}, { deep: true, immediate: true });



onMounted(() => {
  // 如果还未连接，使用 authStore 中的用户信息连接
  if (!gameState.connected && !gameState.connecting) {
    if (authStore.state.user) {
      gameStore.connectWebSocket(authStore.state.user.id, authStore.state.user.name);
    } else if (props.playerId && props.playerName) {
      // 如果没有认证信息，使用 props（备用方案）
      gameStore.connectWebSocket(props.playerId, props.playerName);
    }
  }
  
  if (props.roomId) {
    gameStore.joinRoom(props.roomId);
  }
});

// 断开WebSocket连接
onBeforeUnmount(() => {
  if (gameState.socket) {
    gameState.socket.close();
  }
});

// 选择/取消选择卡牌
const toggleCardSelection = (index) => {
  // index 是排序后的索引，需要转换为原始索引
  // 由于我们使用排序后的牌，这里直接使用排序后的索引
  // 但需要确保 gameStore 中的 selectedCards 使用正确的索引
  const originalIndex = gameState.player?.cards?.findIndex((c, i) => {
    const sorted = sortedCards.value;
    return sorted[index] && c.suit === sorted[index].suit && c.rank === sorted[index].rank;
  }) ?? index;
  
  gameStore.toggleCardSelection(originalIndex);
};

// 检查卡片是否被选中（使用排序后的索引）
const isCardSelected = (sortedIndex) => {
  if (!gameState.selectedCards || !gameState.player?.cards) {
    return false;
  }
  
  const sorted = sortedCards.value;
  const card = sorted[sortedIndex];
  if (!card) return false;
  
  // 找到原始索引
  const originalIndex = gameState.player.cards.findIndex((c) => 
    c.suit === card.suit && c.rank === card.rank
  );
  
  return gameState.selectedCards.includes(originalIndex);
};

// 出牌
const playCards = () => {
  gameStore.playCards();
};

// 过牌
const pass = () => {
  gameStore.pass();
};

// 提示出牌
const hint = () => {
  gameStore.hint();
};

// 获取玩家位置类名
const getPlayerPositionClass = (position) => {
  const myPosition = gameState.player?.position || 0;
  const relativePosition = (position - myPosition + 4) % 4;
  
  switch (relativePosition) {
    case 0: return 'bottom';
    case 1: return 'right';
    case 2: return 'top';
    case 3: return 'left';
    default: return '';
  }
};

// 格式化牌组显示
const formatTableCards = (cardGroup) => {
  if (!cardGroup || !cardGroup.cards) {
    return '无';
  }
  
  return cardGroup.cards.map(card => {
    return `${card.suit === 'hearts' ? '♥' : card.suit === 'diamonds' ? '♦' : card.suit === 'clubs' ? '♣' : '♠'}${getCardRankDisplay(card.rank)}`;
  }).join(' ');
};

// 获取牌面显示
const getCardRankDisplay = (rank) => {
  switch (rank) {
    case 11: return 'J';
    case 12: return 'Q';
    case 13: return 'K';
    case 14: return 'A';
    case 15: return '2';
    case 16: return '3';
    default: return rank.toString();
  }
};

// 开始发牌动画
const startDealingAnimation = (cards) => {
  isDealing.value = true;
  showCards.value = false;
  dealtCards.value = [];
  
  const sorted = sortCards(cards);
  
  // 逐张发牌动画
  sorted.forEach((card, index) => {
    setTimeout(() => {
      dealtCards.value.push(card);
      
      // 最后一张牌发完后，显示所有牌
      if (index === sorted.length - 1) {
        setTimeout(() => {
          isDealing.value = false;
          showCards.value = true;
        }, 150);
      }
    }, index * 80); // 每张牌间隔 80ms
  });
};
</script>

<template>
  
  <div class="game-table">
    <!-- 顶部：其他玩家区域 -->
    <div class="players-top-row">
      <div 
        v-for="player in gameState.otherPlayers" 
        :key="player.id"
        :class="['player-area', getPlayerPositionClass(player.position)]"
      >
        <div class="player-info">
          <div class="player-name">{{ player.name }}</div>
          <div class="card-count">{{ player.card_count }}张</div>
        </div>
        <div class="player-cards">
          <div 
            v-for="i in Math.min(player.card_count, 5)" 
            :key="i" 
            class="card-back"
          ></div>
          <div v-if="player.card_count > 5" class="card-count-more">+{{ player.card_count - 5 }}</div>
        </div>
        <div 
          v-if="gameState.currentPlayer === player.position" 
          class="player-turn-indicator"
        >
          ⭐ 出牌中
        </div>
      </div>
    </div>
    
    <!-- 中央：桌面牌区 -->
    <div class="game-center">
      <div class="table-cards">
        <div class="table-cards-label">桌面牌</div>
        <div class="table-cards-content">
          {{ gameState.tableCards ? formatTableCards(gameState.tableCards) : '暂无出牌' }}
        </div>
      </div>
    </div>
    
    <!-- 玩家手牌区 -->
    <div class="player-hand" v-if="gameState.player && gameState.player.cards && gameState.player.cards.length > 0">
      <div class="hand-cards-container" :class="{ 'dealing': isDealing }">
        <div 
          v-for="(card, index) in sortedCards" 
          :key="`${card.suit}-${card.rank}-${index}`"
          class="hand-card"
          :class="{ 
            selected: isCardSelected(index),
            'card-dealing': isDealing && dealtCards.some(c => c.suit === card.suit && c.rank === card.rank),
            'card-visible': !isDealing || showCards || dealtCards.some(c => c.suit === card.suit && c.rank === card.rank)
          }"
          :style="isDealing && dealtCards.some(c => c.suit === card.suit && c.rank === card.rank) ? { animationDelay: `${dealtCards.findIndex(c => c.suit === card.suit && c.rank === card.rank) * 0.08}s` } : {}"
          @click="toggleCardSelection(index)"
        >
          <PlayingCard :suit="card.suit" :rank="card.rank" />
        </div>
      </div>
      
      <div class="action-buttons">
        <button 
          @click="playCards" 
          :disabled="!isMyTurn || !gameState.selectedCards || gameState.selectedCards.length === 0"
        >
          出牌
        </button>
        <button 
          @click="pass" 
          :disabled="!isMyTurn || (gameState.player && gameState.lastPlayer === gameState.player.position)"
        >
          过
        </button>
        <button @click="hint">提示</button>
      </div>
    </div>
    
    <!-- 加载或等待状态 -->
    <div v-if="!gameState.player || !gameState.player.cards" class="waiting-game">
      <div class="waiting-message">
        <p>等待游戏开始...</p>
        <p v-if="gameState.gameStatus === 'playing'">正在加载游戏数据...</p>
      </div>
    </div>
    
    <!-- 游戏状态指示器 -->
    <div class="game-status">
      <div v-if="isMyTurn" class="my-turn">轮到您出牌</div>
    </div>
  </div>
</template>

<style scoped>
.game-table {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%);
  padding: 0;
  margin: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}


/* 顶部玩家区域 */
.players-top-row {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 0.4rem 0.3rem;
  min-height: 90px;
  max-height: 100px;
  flex-wrap: nowrap;
  gap: 0.3rem;
  flex-shrink: 0;
}

.player-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;
  flex: 1;
  max-width: 25%;
  height: 100%;
}


.player-info {
  background-color: rgba(255, 255, 255, 0.95);
  padding: 0.4rem 0.6rem;
  border-radius: 8px;
  text-align: center;
  margin-bottom: 0.4rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  width: 100%;
}

.player-name {
  font-weight: bold;
  font-size: 0.85rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-count {
  font-size: 0.75rem;
  color: #666;
  margin-top: 0.2rem;
}

.card-count-more {
  width: 25px;
  height: 35px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: bold;
  color: #d32f2f;
  border: 1px solid #d32f2f;
}

.player-cards {
  display: flex;
  gap: 3px;
  justify-content: center;
  flex-wrap: wrap;
}

.card-back {
  width: 25px;
  height: 35px;
  background: linear-gradient(135deg, #b71c1c 0%, #c62828 100%);
  border-radius: 4px;
  border: 1px solid white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.player-turn-indicator {
  background-color: #ffeb3b;
  color: #333;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-size: 0.7rem;
  margin-top: 0.3rem;
  font-weight: bold;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* 中央游戏区域 */
.game-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 0.5rem;
  min-height: 0;
  overflow: hidden;
}

.table-cards {
  background-color: rgba(255, 255, 255, 0.95);
  padding: 0.8rem 1.2rem;
  border-radius: 12px;
  min-width: 200px;
  max-width: 90%;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  border: 2px solid rgba(255, 255, 255, 0.5);
}

.table-cards-label {
  font-weight: bold;
  margin-bottom: 0.5rem;
  font-size: 1rem;
  color: #333;
}

.table-cards-content {
  font-size: 1.2rem;
  color: #d32f2f;
  font-weight: bold;
  min-height: 1.5rem;
}

/* 底部手牌区域 */
.player-hand {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0.4rem 0.3rem;
  background-color: rgba(0, 0, 0, 0.25);
  border-radius: 0;
  flex-shrink: 0;
  max-height: 40vh;
  min-height: 180px;
}

.hand-cards-container {
  display: flex;
  flex-wrap: nowrap;
  overflow-x: auto;
  overflow-y: hidden;
  gap: 6px;
  padding: 8px 10px;
  width: 100%;
  max-width: 100%;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  flex: 1;
  min-height: 0;
  justify-content: center;
  align-items: flex-end;
  box-sizing: border-box;
}

.hand-cards-container::-webkit-scrollbar {
  height: 6px;
}

.hand-cards-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.hand-cards-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.5);
  border-radius: 3px;
}

.hand-card {
  flex-shrink: 0;
  width: 65px;
  height: 95px;
  transition: all 0.3s ease;
  cursor: pointer;
  touch-action: manipulation;
  opacity: 0;
  transform: translateY(100px) scale(0.8);
  will-change: transform, opacity;
}

.hand-card.card-visible {
  opacity: 1;
  transform: translateY(0) scale(1);
  animation: none;
}

.hand-card.card-dealing {
  animation: dealCard 0.4s ease-out forwards;
}

@keyframes dealCard {
  0% {
    opacity: 0;
    transform: translateY(100px) scale(0.8) rotate(-10deg);
  }
  50% {
    opacity: 0.8;
    transform: translateY(-20px) scale(1.1) rotate(5deg);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1) rotate(0deg);
  }
}

.hand-card.selected {
  transform: translateY(-12px) scale(1.08);
  z-index: 10;
  filter: drop-shadow(0 4px 8px rgba(255, 235, 59, 0.8));
}

.hand-card.selected.card-visible {
  transform: translateY(-12px) scale(1.08);
}

.action-buttons {
  display: flex;
  gap: 0.6rem;
  margin-top: 0.6rem;
  width: 100%;
  justify-content: center;
  padding: 0 0.5rem;
  flex-wrap: nowrap;
  flex-shrink: 0;
}

.action-buttons button {
  flex: 1;
  min-width: 0;
  padding: 0.7rem 0.8rem;
  font-size: 0.95rem;
  font-weight: bold;
  border-radius: 8px;
  border: none;
  background: linear-gradient(135deg, #d32f2f 0%, #f44336 100%);
  color: white;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  transition: all 0.2s;
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

.action-buttons button:disabled {
  background: #ccc;
  cursor: not-allowed;
  opacity: 0.6;
}

.action-buttons button:not(:disabled):active {
  transform: scale(0.95);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.game-status {
  position: fixed;
  top: 10px;
  left: 10px;
  z-index: 100;
}

.my-turn {
  background: linear-gradient(135deg, #ffeb3b 0%, #ffc107 100%);
  color: #333;
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  font-weight: bold;
  font-size: 1rem;
  box-shadow: 0 4px 12px rgba(255, 235, 59, 0.6);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}

.waiting-game {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%);
}

.waiting-message {
  text-align: center;
  color: white;
  font-size: 1.5rem;
}

.waiting-message p {
  margin: 1rem 0;
}

</style>