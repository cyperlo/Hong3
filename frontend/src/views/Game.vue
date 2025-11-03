<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue';
import PlayingCard from '../components/PlayingCard.vue';
import gameStore from '../store/gameStore.js';

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

// 连接WebSocket
onMounted(() => {
  gameStore.connectWebSocket(props.playerId, props.playerName);
  
  // 如果有房间ID，加入房间
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
  gameStore.toggleCardSelection(index);
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
    <div class="player-hand" v-if="gameState.player && gameState.player.cards">
      <div class="hand-cards-container">
        <div 
          v-for="(card, index) in gameState.player.cards" 
          :key="`${card.suit}-${card.rank}-${index}`"
          class="hand-card"
          :class="{ selected: gameState.selectedCards && gameState.selectedCards.includes(index) }"
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
  position: relative;
  width: 100%;
  min-height: 100vh;
  background: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%);
  padding: 0.5rem;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 顶部玩家区域（横屏时显示在顶部） */
.players-top-row {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 0.5rem;
  min-height: 100px;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.player-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 80px;
  flex: 1;
  max-width: 120px;
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
  padding: 1rem;
  min-height: 200px;
}

.table-cards {
  background-color: rgba(255, 255, 255, 0.95);
  padding: 1rem 1.5rem;
  border-radius: 12px;
  min-width: 250px;
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
  padding: 0.5rem;
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 12px 12px 0 0;
}

.hand-cards-container {
  display: flex;
  flex-wrap: nowrap;
  overflow-x: auto;
  overflow-y: hidden;
  gap: 8px;
  padding: 10px;
  width: 100%;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
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
  width: 70px;
  height: 100px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.hand-card.selected {
  transform: translateY(-15px) scale(1.05);
  z-index: 10;
  filter: drop-shadow(0 4px 8px rgba(255, 235, 59, 0.8));
}

.action-buttons {
  display: flex;
  gap: 0.8rem;
  margin-top: 0.8rem;
  width: 100%;
  justify-content: center;
  padding: 0 1rem;
  flex-wrap: wrap;
}

.action-buttons button {
  flex: 1;
  min-width: 80px;
  max-width: 120px;
  padding: 0.8rem 1.2rem;
  font-size: 1rem;
  font-weight: bold;
  border-radius: 8px;
  border: none;
  background: linear-gradient(135deg, #d32f2f 0%, #f44336 100%);
  color: white;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  transition: all 0.2s;
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

/* 横屏优化（手机横屏模式） */
@media (orientation: landscape) and (max-height: 600px) {
  .game-table {
    min-height: 100vh;
  }
  
  .players-top-row {
    min-height: 80px;
    padding: 0.3rem;
  }
  
  .player-area {
    min-width: 70px;
    max-width: 100px;
  }
  
  .player-info {
    padding: 0.3rem 0.4rem;
  }
  
  .player-name {
    font-size: 0.75rem;
  }
  
  .card-count {
    font-size: 0.65rem;
  }
  
  .card-back {
    width: 20px;
    height: 28px;
  }
  
  .game-center {
    min-height: 150px;
    padding: 0.5rem;
  }
  
  .table-cards {
    padding: 0.8rem 1rem;
    min-width: 200px;
  }
  
  .hand-card {
    width: 60px;
    height: 85px;
  }
  
  .hand-card.selected {
    transform: translateY(-10px) scale(1.03);
  }
  
  .action-buttons {
    gap: 0.5rem;
    margin-top: 0.5rem;
    padding: 0 0.5rem;
  }
  
  .action-buttons button {
    padding: 0.6rem 1rem;
    font-size: 0.9rem;
    min-width: 70px;
    max-width: 100px;
  }
}

/* 竖屏手机优化 */
@media (orientation: portrait) and (max-width: 768px) {
  .hand-card {
    width: 55px;
    height: 80px;
  }
  
  .action-buttons button {
    padding: 0.7rem 1rem;
    font-size: 0.95rem;
  }
}
</style>