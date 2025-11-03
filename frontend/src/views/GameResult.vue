<script setup>
import { ref } from 'vue';

const props = defineProps({
  result: Object
});

const emit = defineEmits(['return-to-lobby']);

const returnToLobby = () => {
  emit('return-to-lobby');
};
</script>

<template>
  <div class="game-result">
    <h2>游戏结束</h2>
    
    <div class="result-card">
      <div class="winner-info">
        <h3>胜利方</h3>
        <div class="winner-team">
          {{ result.winningTeam === 'team1' ? '队伍1' : '队伍2' }}
        </div>
      </div>
      
      <div class="players-info">
        <h3>玩家信息</h3>
        <div class="players-list">
          <div v-for="player in result.players" :key="player.id" class="player-item">
            <div class="player-name">{{ player.name }}</div>
            <div class="player-stats">
              <div>剩余牌数: {{ player.remainingCards }}</div>
              <div>累计牌数: {{ player.collectedCards }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <button class="return-button" @click="returnToLobby">返回大厅</button>
  </div>
</template>

<style scoped>
.game-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

h2 {
  color: #1b5e20;
  margin-bottom: 2rem;
}

.result-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  width: 100%;
  margin-bottom: 2rem;
}

.winner-info {
  text-align: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.winner-team {
  font-size: 1.5rem;
  font-weight: bold;
  color: #d32f2f;
  margin-top: 0.5rem;
}

.players-info h3 {
  margin-bottom: 1rem;
}

.players-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.player-item {
  background-color: #f5f5f5;
  border-radius: 4px;
  padding: 1rem;
}

.player-name {
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.player-stats {
  color: #666;
  font-size: 0.9rem;
}

.return-button {
  background-color: #1b5e20;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.return-button:hover {
  background-color: #2e7d32;
}

@media (max-width: 768px) {
  .game-result {
    padding: 1rem;
  }
  
  .result-card {
    padding: 1rem;
  }
  
  .players-list {
    grid-template-columns: 1fr;
  }
}
</style>