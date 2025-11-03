<script setup>
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue';
import gameStore from '../store/gameStore.js';

const props = defineProps({
  roomId: String,
  playerId: String,
  playerName: String
});

const emit = defineEmits(['game-start', 'leave-room']);
const { state: gameState } = gameStore;
const isReady = ref(false);

// 计算当前房间信息
const roomState = computed(() => {
  return {
    id: gameState.roomId || props.roomId,
    players: gameState.roomPlayers || [],
    status: gameState.gameStatus
  };
});

// 连接WebSocket
onMounted(() => {
  gameStore.connectWebSocket(props.playerId, props.playerName);
  
  // 如果有房间ID，加入房间
  if (props.roomId || gameState.roomId) {
    const roomId = props.roomId || gameState.roomId;
    gameStore.joinRoom(roomId);
  }
  
  // 设置游戏开始回调
  gameState.onGameStart = () => {
    emit('game-start');
  };
});

// 断开WebSocket连接
onBeforeUnmount(() => {
  if (gameState.socket) {
    gameState.socket.close();
  }
});

// 离开房间
const leaveRoom = () => {
  gameStore.leaveRoom();
  emit('leave-room');
};

// 准备/取消准备
const toggleReady = () => {
  isReady.value = !isReady.value;
  gameStore.readyGame();
};

// 监听玩家准备状态变化
watch(() => gameState.roomPlayers, (newPlayers) => {
  // 更新自己的准备状态显示
  const currentPlayer = newPlayers.find(p => p.id === props.playerId);
  if (currentPlayer) {
    isReady.value = currentPlayer.ready || false;
  }
}, { deep: true });
</script>

<template>
  <div class="room-container">
    <div class="room-header">
      <h2>房间: {{ roomState.id }}</h2>
      <button class="leave-button" @click="leaveRoom">离开房间</button>
    </div>
    
    <div class="room-content">
      <div class="players-section">
        <h3>玩家列表</h3>
        <div class="players-list">
          <div 
            v-for="player in roomState.players" 
            :key="player.id"
            class="player-item"
            :class="{ 'current-player': player.id === playerId, 'player-ready': player.ready }"
          >
            <div class="player-name">{{ player.name }}</div>
            <div class="player-status">{{ player.ready ? '已准备' : '未准备' }}</div>
          </div>
        </div>
      </div>
      
      <div class="game-info">
        <h3>游戏规则</h3>
        <div class="rules-card">
          <p>- 4人扑克牌游戏，无大小王</p>
          <p>- 根据红桃3分布自动判定队伍（1v3或2v2）</p>
          <p>- 牌型：单张、对子、顺子、连队、炸弹（3张）、大炸弹（4张）</p>
          <p>- 首轮必须由拥有红桃4的玩家出牌</p>
          <p>- 出牌类型必须与当前桌面牌型相同，且点数更大</p>
          <p>- 炸弹可以压制任意非炸弹牌型</p>
        </div>
      </div>
    </div>
    
    <div class="room-actions">
      <button 
        class="ready-button" 
        :class="{ 'ready-active': isReady }"
        @click="toggleReady"
      >
        {{ isReady ? '取消准备' : '准备' }}
      </button>
    </div>
    
    <div v-if="gameState.error" class="error-message">
      {{ gameState.error }}
    </div>
  </div>
</template>

<style scoped>
.room-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 1rem;
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.leave-button {
  background-color: #d32f2f;
}

.leave-button:hover {
  background-color: #b71c1c;
}

.room-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.players-section, .game-info {
  background-color: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h3 {
  margin-bottom: 1rem;
  color: #1b5e20;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.player-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.current-player {
  background-color: #e8f5e9;
  border-left: 4px solid #1b5e20;
}

.player-ready {
  background-color: #e8f5e9;
}

.player-name {
  font-weight: bold;
}

.player-status {
  color: #666;
}

.rules-card p {
  margin-bottom: 0.5rem;
}

.room-actions {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.ready-button {
  padding: 0.75rem 2rem;
  font-size: 1.1rem;
}

.ready-active {
  background-color: #d32f2f;
}

.ready-active:hover {
  background-color: #b71c1c;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #ffebee;
  color: #d32f2f;
  border-radius: 4px;
  text-align: center;
}

@media (max-width: 768px) {
  .room-content {
    grid-template-columns: 1fr;
  }
}
</style>