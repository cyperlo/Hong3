<script setup>
import { ref, onMounted, watch, computed } from 'vue';
import Login from './views/Login.vue';
import Lobby from './views/Lobby.vue';
import Room from './views/Room.vue';
import Game from './views/Game.vue';
import GameResult from './views/GameResult.vue';
import authStore from './store/authStore';

// 游戏状态
const currentView = ref('login'); // login, lobby, room, game, result
const roomId = ref('');
const playerId = ref('');
const playerName = ref('');

// 检查是否已登录
const isAuthenticated = computed(() => authStore.state.isAuthenticated);

// 从认证信息中获取玩家信息
onMounted(() => {
  if (isAuthenticated.value && authStore.state.user) {
    // 已登录，使用用户信息
    playerId.value = authStore.state.user.id;
    playerName.value = authStore.state.user.name;
    currentView.value = 'lobby';
  } else {
    // 未登录，显示登录页面
    currentView.value = 'login';
  }
});

// 处理登录成功
const handleLoginSuccess = () => {
  if (authStore.state.user) {
    playerId.value = authStore.state.user.id;
    playerName.value = authStore.state.user.name;
    currentView.value = 'lobby';
  }
};

// 切换视图
const switchView = (view, params = {}) => {
  currentView.value = view;
  if (params.roomId) {
    roomId.value = params.roomId;
  }
};

// 离开房间处理
const handleLeaveRoom = () => {
  // 直接切换视图，Room 组件会在 onBeforeUnmount 中清理
  // 避免循环调用：如果调用 roomComponent.value.leaveRoom() 会 emit 事件，又会触发 handleLeaveRoom
  switchView('lobby');
  roomId.value = '';
};

// 监听视图变化，切换 body 的全屏类
watch(currentView, (newView) => {
  if (newView === 'game') {
    document.body.classList.add('game-fullscreen');
  } else {
    document.body.classList.remove('game-fullscreen');
  }
});
</script>

<template>
  <div class="app-container" :class="{ 'game-view': currentView === 'game' }">
    <!-- 登录页面 -->
    <Login 
      v-if="currentView === 'login'" 
      @login-success="handleLoginSuccess"
    />
    
    <!-- 游戏页面 -->
    <template v-else>
      <header v-if="currentView !== 'game'" class="app-header" :class="{ 'room-header-mode': currentView === 'room', 'lobby-header-mode': currentView === 'lobby' }">
        <div class="header-content">
          <h1 class="app-title">
            <span class="title-red">红</span><span class="title-number">3</span>
          </h1>
          <div v-if="currentView === 'room'" class="room-header-info">
            <div class="room-badge">
              <span class="room-icon">🏠</span>
              <span class="room-id-text">房间 {{ roomId }}</span>
            </div>
            <button class="header-leave-btn" @click="handleLeaveRoom">
              <span class="leave-icon">←</span>
              离开
            </button>
          </div>
        </div>
      </header>
      
      <main class="app-content">
        <Lobby 
          v-if="currentView === 'lobby'" 
          :playerId="playerId" 
          :playerName="playerName"
          @join-room="switchView('room', { roomId: $event })"
          @create-room="switchView('room', { roomId: $event })"
        />
        
        <Room 
          v-if="currentView === 'room'" 
          :roomId="roomId"
          :playerId="playerId" 
          :playerName="playerName"
          @game-start="switchView('game')"
          @leave-room="handleLeaveRoom"
        />
        
        <Game 
          v-if="currentView === 'game'" 
          :roomId="roomId"
          :playerId="playerId" 
          :playerName="playerName"
          @game-end="switchView('result', { result: $event })"
        />
        
        <GameResult 
          v-if="currentView === 'result'" 
          @back-to-lobby="switchView('lobby')"
        />
      </main>
      
      <footer v-if="currentView !== 'game'" class="app-footer">
        <p>© 2025 红3</p>
      </footer>
    </template>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Microsoft YaHei', Arial, sans-serif;
  background-color: #f0f2f5;
  color: #333;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  width: 100%;
  position: relative;
  /* 全屏显示 */
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}

.app-container.game-view {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  min-height: 100vh;
  overflow: hidden;
}

/* 移动端横屏优化 */
@media screen and (max-width: 768px) {
  .app-container {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
  }
  
  .app-content {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
  }
}

/* 游戏界面时，给 body 添加全屏类 */
.app-container.game-view ~ * {
  /* 这个选择器用于确保游戏界面时 body 是全屏的 */
}

.app-header {
  background: linear-gradient(135deg, #d32f2f 0%, #b71c1c 100%);
  color: white;
  padding: 1.5rem 1rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.app-header.lobby-header-mode {
  background: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%);
  box-shadow: 0 2px 8px rgba(27, 94, 32, 0.15);
}

.app-header.room-header-mode {
  padding: 1rem;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.app-title {
  font-size: 2.5rem;
  font-weight: bold;
  margin: 0;
  letter-spacing: 2px;
  display: inline-block;
}

.app-header.room-header-mode .app-title {
  font-size: 2rem;
}

.title-red {
  color: #ffeb3b;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.app-header.lobby-header-mode .title-red {
  color: #ffeb3b;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
}

.title-number {
  color: white;
  margin-left: 0.2rem;
}

/* 房间信息显示 */
.room-header-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.room-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.2);
  padding: 0.5rem 1rem;
  border-radius: 20px;
  backdrop-filter: blur(10px);
  font-weight: 600;
}

.room-icon {
  font-size: 1.2rem;
}

.room-id-text {
  font-size: 1rem;
}

.header-leave-btn {
  background: rgba(255, 255, 255, 0.25);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.5);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  backdrop-filter: blur(10px);
}

.header-leave-btn:hover {
  background: rgba(255, 255, 255, 0.35);
  border-color: rgba(255, 255, 255, 0.7);
  transform: translateX(-2px);
}

.header-leave-btn:active {
  transform: translateX(0);
}

.leave-icon {
  font-size: 1.1rem;
  font-weight: bold;
}

.app-content {
  flex: 1;
  padding: 1rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  overflow: hidden;
}

.app-container.game-view .app-content {
  padding: 0;
  max-width: 100%;
  margin: 0;
  height: 100%;
  overflow: hidden;
}

.app-footer {
  background-color: #1a1a1a;
  color: rgba(255, 255, 255, 0.7);
  text-align: center;
  padding: 0.75rem;
  margin-top: auto;
  font-size: 0.85rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-content {
    padding: 0.5rem;
  }
  
  .app-title {
    font-size: 1.8rem;
  }
  
  .app-header {
    padding: 1rem 0.5rem;
  }
  
  .app-header.room-header-mode {
    padding: 0.75rem 0.5rem;
  }
  
  .header-content {
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .room-header-info {
    width: 100%;
    justify-content: space-between;
  }
  
  .room-badge {
    font-size: 0.9rem;
    padding: 0.4rem 0.8rem;
  }
  
  .header-leave-btn {
    padding: 0.4rem 0.9rem;
    font-size: 0.9rem;
  }
}

button {
  background-color: #d32f2f;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #b71c1c;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 1rem;
}
</style>
