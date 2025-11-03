<script setup>
import { computed } from 'vue';

const props = defineProps({
  suit: {
    type: String,
    required: true,
    validator: (value) => ['hearts', 'diamonds', 'clubs', 'spades'].includes(value)
  },
  rank: {
    type: Number,
    required: true,
    validator: (value) => value >= 4 && value <= 16
  }
});

const suitSymbol = computed(() => {
  switch (props.suit) {
    case 'hearts': return '♥';
    case 'diamonds': return '♦';
    case 'clubs': return '♣';
    case 'spades': return '♠';
    default: return '';
  }
});

const rankDisplay = computed(() => {
  switch (props.rank) {
    case 11: return 'J';
    case 12: return 'Q';
    case 13: return 'K';
    case 14: return 'A';
    case 15: return '2';
    case 16: return '3';
    default: return props.rank.toString();
  }
});

const isRed = computed(() => {
  return props.suit === 'hearts' || props.suit === 'diamonds';
});

const isHeartThree = computed(() => {
  return props.suit === 'hearts' && props.rank === 16;
});

const isHeartFour = computed(() => {
  return props.suit === 'hearts' && props.rank === 4;
});
</script>

<template>
  <div class="playing-card" :class="{ 'red-card': isRed, 'heart-three': isHeartThree, 'heart-four': isHeartFour }">
    <div class="card-corner top-left">
      <div class="card-rank">{{ rankDisplay }}</div>
      <div class="card-suit">{{ suitSymbol }}</div>
    </div>
    
    <div class="card-center">
      <div class="card-suit-big">{{ suitSymbol }}</div>
    </div>
    
    <div class="card-corner bottom-right">
      <div class="card-rank">{{ rankDisplay }}</div>
      <div class="card-suit">{{ suitSymbol }}</div>
    </div>
  </div>
</template>

<style scoped>
.playing-card {
  width: 100%;
  height: 100%;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 5px;
  position: relative;
  color: black;
  font-family: Arial, sans-serif;
}

.red-card {
  color: #d32f2f;
}

.heart-three {
  background-color: #ffebee;
  box-shadow: 0 0 10px #ff5252;
}

.heart-four {
  background-color: #e8f5e9;
  box-shadow: 0 0 10px #69f0ae;
}

.card-corner {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-weight: bold;
}

.top-left {
  position: absolute;
  top: 5px;
  left: 5px;
}

.bottom-right {
  position: absolute;
  bottom: 5px;
  right: 5px;
  transform: rotate(180deg);
}

.card-rank {
  font-size: 1.2em;
}

.card-suit {
  font-size: 1em;
}

.card-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.card-suit-big {
  font-size: 2.5em;
}

@media (max-width: 768px) {
  .card-rank {
    font-size: 1em;
  }
  
  .card-suit {
    font-size: 0.8em;
  }
  
  .card-suit-big {
    font-size: 2em;
  }
}
</style>