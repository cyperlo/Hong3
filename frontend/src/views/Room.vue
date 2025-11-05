<script setup>
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue';
import gameStore from '../store/gameStore.js';
import authStore from '../store/authStore.js';

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

// 连接WebSocket（如果还未连接）
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

// 断开WebSocket连接和清理房间状态
onBeforeUnmount(() => {
  // 清理房间状态
  gameStore.leaveRoom();
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

// 暴露方法给父组件
defineExpose({
  leaveRoom
});
</script>

<template>
  <div class="room-container">
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
  max-width: 900px;
  margin: 0 auto;
  padding: 1rem;
  min-height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
}

.room-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
  flex: 1;
}

.players-section, .game-info {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border: 1px solid #f0f0f0;
}

h3 {
  margin-bottom: 1.25rem;
  color: #1b5e20;
  font-size: 1.3rem;
  font-weight: bold;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid #e8f5e9;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.player-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  background: linear-gradient(135deg, #f5f5f5 0%, #fafafa 100%);
  border-radius: 8px;
  margin-bottom: 0.5rem;
  transition: all 0.3s;
}

.player-item:hover {
  transform: translateX(4px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.current-player {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  border-left: 4px solid #1b5e20;
  font-weight: 600;
}

.player-ready {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}

.player-name {
  font-weight: bold;
}

.player-status {
  color: #666;
}

.rules-card {
  line-height: 1.8;
}

.rules-card p {
  margin-bottom: 0.75rem;
  color: #555;
  padding-left: 1rem;
  position: relative;
}

.rules-card p::before {
  content: '•';
  position: absolute;
  left: 0;
  color: #1b5e20;
  font-weight: bold;
  font-size: 1.2rem;
}

.room-actions {
  display: flex;
  justify-content: center;
  margin-top: auto;
  padding: 1rem 0;
  position: sticky;
  bottom: 0;
  background-color: #f0f2f5;
  z-index: 10;
}

.ready-button {
  padding: 1rem 3rem;
  font-size: 1.2rem;
  font-weight: bold;
  min-height: 56px;
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s;
  letter-spacing: 1px;
}

.ready-button:not(.ready-active) {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
}

.ready-active {
  background: linear-gradient(135deg, #d32f2f 0%, #b71c1c 100%);
}

.ready-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.ready-button:active {
  transform: translateY(0);
}

.error-message {
  margin-top: 1rem;
  padding: 1rem;
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  color: #c62828;
  border-radius: 8px;
  text-align: center;
  border-left: 4px solid #c62828;
  font-weight: 500;
}

@media (max-width: 768px) {
  .room-container {
    padding: 0.5rem;
    min-height: calc(100vh - 180px);
  }
  
  .room-content {
    grid-template-columns: 1fr;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  .room-actions {
    padding: 0.75rem 0;
    position: sticky;
    bottom: 0;
    background-color: #f0f2f5;
    margin-top: auto;
  }
  
  .ready-button {
    width: 100%;
    max-width: 300px;
    padding: 1rem 2rem;
    font-size: 1.2rem;
  }
}
</style>