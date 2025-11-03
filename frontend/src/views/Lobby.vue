<template>
  <div class="lobby">
    <div class="header">
      <h1>游戏大厅</h1>
      <div v-if="gameState.connected" class="player-info">
        欢迎, {{ gameState.playerName }}
      </div>
    </div>

    <div v-if="!gameState.connected" class="login-form">
      <input
        v-model="playerName"
        type="text"
        placeholder="输入你的名字"
        @keyup.enter="handleConnect"
      />
      <button
        @click="handleConnect"
        :disabled="!playerName || gameState.connecting"
      >
        {{ gameState.connecting ? '连接中...' : '进入游戏' }}
      </button>
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
import { ref, watch } from 'vue';
import gameStore from '../store/gameStore';
import { v4 as uuidv4 } from 'uuid';

export default {
  name: 'Lobby',
  emits: ['join-room', 'create-room'],
  setup(props, { emit }) {
    const playerName = ref('');
    const { state: gameState } = gameStore;
    const hasJumpedToRoom = ref(false); // 防止重复跳转
    
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

    const handleConnect = () => {
      if (playerName.value) {
        const playerId = uuidv4();
        gameStore.connectWebSocket(playerId, playerName.value);
      }
    };

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
      playerName,
      gameState,
      handleConnect,
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
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.player-info {
  font-size: 1.1em;
  color: #666;
}

.login-form {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.login-form input {
  flex: 1;
  padding: 8px;
  font-size: 1em;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.room-controls {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.create-room-btn,
.refresh-btn {
  padding: 8px 16px;
  font-size: 1em;
  cursor: pointer;
  border: none;
  border-radius: 4px;
  background-color: #4CAF50;
  color: white;
  transition: background-color 0.3s;
}

.create-room-btn:hover,
.refresh-btn:hover {
  background-color: #45a049;
}

.create-room-btn:disabled,
.refresh-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.loading,
.no-rooms {
  text-align: center;
  padding: 20px;
  color: #666;
}

.rooms {
  display: grid;
  gap: 15px;
}

.room-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.room-item:hover {
  background-color: #e8e8e8;
}

.room-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.room-id {
  font-weight: bold;
  color: #333;
}

.player-count {
  color: #666;
  font-size: 0.9em;
}

.join-btn {
  padding: 6px 12px;
  font-size: 0.9em;
  cursor: pointer;
  border: none;
  border-radius: 4px;
  background-color: #2196F3;
  color: white;
  transition: background-color 0.3s;
}

.join-btn:hover {
  background-color: #1976D2;
}

.join-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.room {
  background-color: #f5f5f5;
  padding: 20px;
  border-radius: 4px;
}

.player-list {
  margin: 20px 0;
}

.player-list ul {
  list-style: none;
  padding: 0;
}

.player-list li {
  padding: 8px;
  margin: 5px 0;
  background-color: white;
  border-radius: 4px;
}

.current-player {
  font-weight: bold;
  background-color: #e3f2fd;
}

.leave-btn {
  padding: 8px 16px;
  font-size: 1em;
  cursor: pointer;
  border: none;
  border-radius: 4px;
  background-color: #f44336;
  color: white;
  transition: background-color 0.3s;
}

.leave-btn:hover {
  background-color: #d32f2f;
}

.error {
  margin-top: 20px;
  padding: 10px;
  background-color: #ffebee;
  color: #c62828;
  border-radius: 4px;
}
</style>