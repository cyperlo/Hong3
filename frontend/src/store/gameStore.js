import { reactive } from 'vue';

const gameState = reactive({
  connected: false,
  connecting: false,
  playerId: null,
  playerName: null,
  roomId: null,
  roomPlayers: [],
  error: null,
  availableRooms: [],
  loadingRooms: false,
  gameStatus: null, // 游戏状态: 'waiting', 'playing', 'finished'
  onGameStart: null, // 游戏开始回调
  // 游戏数据
  player: null, // 当前玩家的数据（包含手牌）
  otherPlayers: [], // 其他玩家数据
  selectedCards: [], // 选中的手牌索引
  tableCards: null, // 桌面上的牌
  currentPlayer: null, // 当前出牌玩家的位置
  lastPlayer: null, // 最后出牌玩家的位置
});

const ws = {
  connection: null,
};

// 获取后端服务器地址（动态获取，支持通过 IP 访问）
const getBackendUrl = () => {
  const hostname = window.location.hostname;
  const port = '8080';
  // 根据当前页面协议选择 WebSocket 协议
  const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  return {
    ws: `${wsProtocol}://${hostname}:${port}`,
    http: `${window.location.protocol}//${hostname}:${port}`,
  };
};

const connectWebSocket = async (playerId, playerName) => {
  if (gameState.connecting || gameState.connected) {
    return;
  }

  gameState.connecting = true;
  gameState.error = null;

  try {
    const backend = getBackendUrl();
    const wsUrl = `${backend.ws}/ws?player_id=${playerId}&player_name=${encodeURIComponent(playerName)}`;
    console.log('连接 WebSocket:', wsUrl);
    ws.connection = new WebSocket(wsUrl);

    ws.connection.onopen = () => {
      console.log('WebSocket连接已建立');
      gameState.connected = true;
      gameState.connecting = false;
      gameState.playerId = playerId;
      gameState.playerName = playerName;
      // 连接成功后立即获取房间列表
      fetchRooms();
    };

    ws.connection.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        console.log('收到消息:', message);
        handleMessage(message);
      } catch (error) {
        console.error('消息解析错误:', error);
        gameState.error = '消息解析错误';
      }
    };

    ws.connection.onclose = () => {
      console.log('WebSocket连接已关闭');
      gameState.connected = false;
      gameState.connecting = false;
      gameState.roomId = null;
      gameState.roomPlayers = [];
    };

    ws.connection.onerror = (error) => {
      console.error('WebSocket错误:', error);
      gameState.error = 'WebSocket连接错误';
      gameState.connected = false;
      gameState.connecting = false;
    };
  } catch (error) {
    console.error('连接错误:', error);
    gameState.error = '连接错误';
    gameState.connecting = false;
  }
};

const handleMessage = (message) => {
  switch (message.type) {
    case 'room_state':
      updateRoomState(message);
      break;
    case 'room_created':
      // 房间创建成功，等待 room_state 消息
      console.log('房间创建成功:', message.room_id);
      // 刷新房间列表，让其他玩家能看到新创建的房间
      fetchRooms();
      break;
    case 'player_joined':
      handlePlayerJoined(message);
      break;
    case 'player_left':
      handlePlayerLeft(message);
      break;
    case 'player_ready':
      handlePlayerReady(message);
      break;
    case 'game_started':
      handleGameStarted(message);
      break;
    case 'game_state':
      handleGameState(message);
      break;
    case 'error':
      gameState.error = message.message || message.error;
      console.error('收到错误消息:', message);
      break;
    default:
      console.log('未处理的消息类型:', message.type, message);
  }
};

const updateRoomState = (message) => {
  console.log('更新房间状态:', message);
  gameState.roomId = message.room_id;
  // 确保每个玩家都有 ready 属性
  const players = (message.players || []).map(p => ({
    ...p,
    ready: p.ready || false
  }));
  gameState.roomPlayers = players;
  // 房间状态更新后刷新房间列表
  fetchRooms();
};

const handlePlayerJoined = (message) => {
  console.log('玩家加入:', message);
  // 如果当前用户在房间中，刷新房间列表以获取最新玩家信息
  if (gameState.roomId) {
    fetchRooms();
    // 也可以直接添加玩家到列表（如果消息包含完整信息）
    if (message.playerID && message.name) {
      const existingPlayer = gameState.roomPlayers.find(p => p.id === message.playerID);
      if (!existingPlayer) {
        gameState.roomPlayers.push({
          id: message.playerID,
          name: message.name
        });
      }
    }
  }
};

const handlePlayerLeft = (message) => {
  console.log('玩家离开:', message);
  // 从房间玩家列表中移除离开的玩家
  if (gameState.roomId && message.playerID) {
    gameState.roomPlayers = gameState.roomPlayers.filter(
      p => p.id !== message.playerID
    );
    // 刷新房间列表
    fetchRooms();
  }
};

const sendMessage = (message) => {
  if (!ws.connection || ws.connection.readyState !== WebSocket.OPEN) {
    gameState.error = 'WebSocket未连接';
    return;
  }

  try {
    console.log('发送消息:', message);
    ws.connection.send(JSON.stringify(message));
  } catch (error) {
    console.error('发送消息错误:', error);
    gameState.error = '发送消息错误';
  }
};

const createRoom = () => {
  sendMessage({ type: 'create_room' });
};

const joinRoom = (roomId) => {
  sendMessage({ type: 'join_room', room_id: roomId });
};

const leaveRoom = () => {
  if (gameState.roomId) {
    sendMessage({ type: 'leave_room' });
    gameState.roomId = null;
    gameState.roomPlayers = [];
    gameState.gameStatus = null;
    // 清除游戏数据
    gameState.player = null;
    gameState.otherPlayers = [];
    gameState.selectedCards = [];
    gameState.tableCards = null;
    gameState.currentPlayer = null;
    gameState.lastPlayer = null;
  }
};

// 选择/取消选择卡牌
const toggleCardSelection = (index) => {
  if (!gameState.selectedCards) {
    gameState.selectedCards = [];
  }
  const idx = gameState.selectedCards.indexOf(index);
  if (idx > -1) {
    gameState.selectedCards.splice(idx, 1);
  } else {
    gameState.selectedCards.push(index);
  }
};

// 出牌
const playCards = () => {
  if (!gameState.roomId || !gameState.selectedCards || gameState.selectedCards.length === 0) {
    return;
  }
  sendMessage({
    type: 'game_action',
    action: 'play_cards',
    card_indices: gameState.selectedCards,
  });
  gameState.selectedCards = [];
};

// 过牌
const pass = () => {
  if (!gameState.roomId) {
    return;
  }
  sendMessage({
    type: 'game_action',
    action: 'pass',
  });
};

// 提示出牌
const hint = () => {
  // TODO: 实现提示功能
  console.log('提示功能待实现');
};

const readyGame = () => {
  if (!gameState.roomId) {
    gameState.error = '未加入房间';
    return;
  }
  sendMessage({
    type: 'game_action',
    action: 'ready',
  });
};

const handlePlayerReady = (message) => {
  console.log('玩家准备:', message);
  // 更新玩家准备状态
  if (gameState.roomId && message.playerID) {
    const player = gameState.roomPlayers.find(p => p.id === message.playerID);
    if (player) {
      player.ready = true;
    } else {
      // 如果玩家不在列表中，刷新房间列表
      fetchRooms();
    }
  }
};

const handleGameStarted = (message) => {
  console.log('游戏开始:', message);
  gameState.gameStatus = 'playing';
  if (message.currentPlayer !== undefined) {
    gameState.currentPlayer = message.currentPlayer;
  }
  // 触发游戏开始回调
  if (gameState.onGameStart) {
    console.log('触发游戏开始回调');
    gameState.onGameStart();
  } else {
    console.warn('游戏开始回调未设置');
  }
};

const handleGameState = (message) => {
  console.log('游戏状态更新:', message);
  // 更新游戏状态
  if (message.status) {
    gameState.gameStatus = message.status;
  }
  
  // 更新当前玩家数据（包含手牌）
  if (message.player) {
    gameState.player = message.player;
    // 初始化选中卡片数组
    if (!gameState.selectedCards) {
      gameState.selectedCards = [];
    }
  }
  
  // 更新其他玩家数据
  if (message.other_players) {
    gameState.otherPlayers = message.other_players;
  }
  
  // 更新桌面牌
  if (message.table_cards !== undefined) {
    gameState.tableCards = message.table_cards;
  }
  
  // 更新当前出牌玩家
  if (message.current_player !== undefined) {
    gameState.currentPlayer = message.current_player;
  }
  
  // 更新最后出牌玩家
  if (message.last_player !== undefined) {
    gameState.lastPlayer = message.last_player;
  }
  
  // 如果收到手牌数据但游戏状态还是 waiting，更新为 playing
  if (message.player && message.player.cards && message.player.cards.length > 0) {
    if (gameState.gameStatus !== 'playing') {
      console.log('收到手牌数据，更新游戏状态为 playing');
      gameState.gameStatus = 'playing';
    }
    // 如果游戏已开始但没有触发回调，尝试触发
    if (gameState.onGameStart) {
      console.log('收到手牌数据，触发游戏开始回调');
      const callback = gameState.onGameStart;
      // 延迟一下确保状态已更新
      setTimeout(() => {
        if (gameState.gameStatus === 'playing' && gameState.player && gameState.player.cards) {
          callback();
        }
      }, 200);
    }
  }
};

const fetchRooms = async () => {
  if (gameState.loadingRooms) {
    return;
  }

  gameState.loadingRooms = true;
  try {
    const backend = getBackendUrl();
    const response = await fetch(`${backend.http}/api/rooms`);
    if (!response.ok) {
      throw new Error('获取房间列表失败');
    }
    const rooms = await response.json();
    console.log('获取到房间列表:', rooms);
    // 确保 rooms 是一个数组
    gameState.availableRooms = Array.isArray(rooms) ? rooms : [];
  } catch (error) {
    console.error('获取房间列表错误:', error);
    gameState.error = '获取房间列表失败';
  } finally {
    gameState.loadingRooms = false;
  }
};

export default {
  state: gameState,
  connectWebSocket,
  createRoom,
  joinRoom,
  leaveRoom,
  fetchRooms,
  readyGame,
  toggleCardSelection,
  playCards,
  pass,
  hint,
};