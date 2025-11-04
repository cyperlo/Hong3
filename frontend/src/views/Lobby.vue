<template>
  <div class="lobby">
    <div class="header">
      <h1>游戏大厅</h1>
      <div v-if="gameState.connected" class="player-info">
        欢迎, {{ gameState.playerName }}
      </div>
    </div>

    <div v-if="!gameState.connected" class="connecting-status">
      <div class="connecting-message">
        <p>正在连接游戏服务器...</p>
        <div class="loading-spinner"></div>
      </div>
    </div>

    <div v-else-if="!gameState.roomId" class="room-list">
      <div class="room-controls">
        <button
          @click="handleCreateRoom"
          :disabled="gameState.loadingRooms"
          class="create-room-btn"
        >
          创建房间
        </button>
        <button
          @click="handleRefreshRooms"
          :disabled="gameState.loadingRooms"
          class="refresh-btn"
        >
          刷新房间列表
        </button>
      </div>

      <div v-if="gameState.loadingRooms" class="loading">
        加载房间列表中...
      </div>
      <div v-else-if="gameState.availableRooms.length === 0" class="no-rooms">
        暂无可用房间
      </div>
      <div v-else class="rooms">
        <div
          v-for="room in gameState.availableRooms"
          :key="room.id"
          class="room-item"
        >
          <div class="room-info">
            <span class="room-id">房间 #{{ room.id }}</span>
            <span class="player-count" v-if="room.players">
              玩家: {{ room.players.length }}/{{ room.capacity }}
            </span>
          </div>
          <button
            @click="handleJoinRoom(room.id)"
            :disabled="!room.players || room.players.length >= room.capacity"
            class="join-btn"
          >
            加入
          </button>
        </div>
      </div>
    </div>

    <div v-else class="room">
      <h2>房间 #{{ gameState.roomId }}</h2>
      <div class="player-list">
        <h3>玩家列表:</h3>
        <ul>
          <li
            v-for="player in gameState.roomPlayers"
            :key="player.id"
            :class="{ 'current-player': player.id === gameState.playerId }"
          >
            {{ player.name }}
            {{ player.id === gameState.playerId ? '(你)' : '' }}
          </li>
        </ul>
      </div>
      <button @click="handleLeaveRoom" class="leave-btn">
        离开房间
      </button>
    </div>

    <div v-if="gameState.error" class="error">
      {{ gameState.error }}
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue';
import gameStore from '../store/gameStore';
import authStore from '../store/authStore';

export default {
  name: 'Lobby',
  emits: ['join-room', 'create-room'],
  setup(props, { emit }) {
    const { state: gameState } = gameStore;
    const hasJumpedToRoom = ref(false); // 防止重复跳转
    
    // 组件挂载时自动连接 WebSocket（使用登录时的用户信息）
    onMounted(() => {
      if (authStore.state.user && !gameState.connected && !gameState.connecting) {
        // 使用登录时的用户 ID 和昵称连接 WebSocket
        gameStore.connectWebSocket(authStore.state.user.id, authStore.state.user.name);
      }
    });
    
    // 监听房间状态变化，自动跳转到 Room 视图
    watch(() => gameState.roomId, (newRoomId) => {
      if (newRoomId && gameState.roomPlayers.length > 0 && !hasJumpedToRoom.value) {
        // 已有房间且玩家列表不为空，跳转到房间视图
        hasJumpedToRoom.value = true;
        setTimeout(() => {
          emit('join-room', newRoomId);
        }, 300); // 延迟确保状态已更新
      }
    });

    const handleCreateRoom = () => {
      gameStore.createRoom();
      // 监听房间创建，跳转到房间视图
      const checkRoom = setInterval(() => {
        if (gameState.roomId) {
          clearInterval(checkRoom);
          emit('create-room', gameState.roomId);
        }
      }, 100);
      // 5秒后停止检查
      setTimeout(() => clearInterval(checkRoom), 5000);
    };

    const handleJoinRoom = (roomId) => {
      gameStore.joinRoom(roomId);
      // 跳转到房间视图
      setTimeout(() => {
        // 延迟一下确保 room_state 消息已收到
        if (gameState.roomId) {
          emit('join-room', roomId || gameState.roomId);
        }
      }, 100);
    };

    const handleLeaveRoom = () => {
      gameStore.leaveRoom();
    };

    const handleRefreshRooms = () => {
      gameStore.fetchRooms();
    };

    return {
      gameState,
      handleCreateRoom,
      handleJoinRoom,
      handleLeaveRoom,
      handleRefreshRooms,
    };
  },
};
</script>

<style scoped>
.lobby {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.header {
  text-align: center;
  margin-bottom: 2rem;
}

.header h1 {
  font-size: 2rem;
  color: #1b5e20;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.player-info {
  font-size: 1rem;
  color: #666;
  background-color: #e8f5e9;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  display: inline-block;
  margin-top: 0.5rem;
}

.connecting-status {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.connecting-message {
  text-align: center;
  color: #666;
}

.connecting-message p {
  margin-bottom: 1rem;
  font-size: 1.1rem;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #4CAF50;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.room-controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  justify-content: center;
}

.create-room-btn,
.refresh-btn {
  padding: 0.875rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(76, 175, 80, 0.3);
}

.create-room-btn:hover,
.refresh-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.4);
}

.create-room-btn:disabled,
.refresh-btn:disabled {
  background: #e0e0e0;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.loading,
.no-rooms {
  text-align: center;
  padding: 3rem;
  color: #666;
  font-size: 1.1rem;
}

.rooms {
  display: grid;
  gap: 1rem;
}

.room-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  border: 2px solid transparent;
}

.room-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  border-color: #e0e0e0;
}

.room-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.room-id {
  font-weight: 600;
  color: #1b5e20;
  font-size: 1.1rem;
}

.player-count {
  color: #666;
  font-size: 0.9rem;
}

.join-btn {
  padding: 0.625rem 1.5rem;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
  color: white;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.3);
}

.join-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.4);
}

.join-btn:disabled {
  background: #e0e0e0;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.room {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.player-list {
  margin: 1.5rem 0;
}

.player-list h3 {
  color: #1b5e20;
  margin-bottom: 1rem;
  font-size: 1.2rem;
}

.player-list ul {
  list-style: none;
  padding: 0;
}

.player-list li {
  padding: 0.875rem 1rem;
  margin: 0.5rem 0;
  background-color: #f5f5f5;
  border-radius: 8px;
  transition: all 0.3s;
}

.current-player {
  font-weight: 600;
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  border-left: 4px solid #1b5e20;
}

.leave-btn {
  padding: 0.875rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #f44336 0%, #d32f2f 100%);
  color: white;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(244, 67, 54, 0.3);
}

.leave-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(244, 67, 54, 0.4);
}

.error {
  margin-top: 1.5rem;
  padding: 1rem;
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  color: #c62828;
  border-radius: 8px;
  border-left: 4px solid #c62828;
  font-weight: 500;
}

@media (max-width: 768px) {
  .lobby {
    padding: 1rem 0.5rem;
  }
  
  .header h1 {
    font-size: 1.5rem;
  }
  
  .login-form {
    padding: 1.5rem;
  }
  
  .room-controls {
    flex-direction: column;
  }
  
  .create-room-btn,
  .refresh-btn {
    width: 100%;
  }
  
  .room-item {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .join-btn {
    width: 100%;
  }
}
</style>