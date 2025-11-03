<script setup>
import { ref, onMounted } from 'vue';
import Lobby from './views/Lobby.vue';
import Room from './views/Room.vue';
import Game from './views/Game.vue';
import GameResult from './views/GameResult.vue';

// 游戏状态
const currentView = ref('lobby'); // lobby, room, game, result
const roomId = ref('');
const playerId = ref('');
const playerName = ref('');

// 生成随机玩家ID和名称
onMounted(() => {
  playerId.value = 'player_' + Math.random().toString(36).substring(2, 10);
  playerName.value = '玩家' + Math.floor(Math.random() * 1000);
});

// 切换视图
const switchView = (view, params = {}) => {
  currentView.value = view;
  if (params.roomId) {
    roomId.value = params.roomId;
  }
};
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>Hong3 - 红桃3扑克牌游戏</h1>
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
        @leave-room="switchView('lobby')"
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
    
    <footer class="app-footer">
      <p>© 2025 Hong3 Card Game</p>
    </footer>
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
}

.app-header {
  background-color: #d32f2f;
  color: white;
  padding: 1rem;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.app-content {
  flex: 1;
  padding: 1rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.app-footer {
  background-color: #333;
  color: white;
  text-align: center;
  padding: 1rem;
  margin-top: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-content {
    padding: 0.5rem;
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
