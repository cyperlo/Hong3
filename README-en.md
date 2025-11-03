# Hong3 - Hearts Three Card Game

A real-time 4-player card game based on WebSocket technology.

## 🎮 Game Introduction

Hong3 (Hearts Three) is a classic 4-player card game using a standard deck without Jokers. The game automatically determines teams (1v3 or 2v2) based on the distribution of Hearts Three, but the frontend does not display team information, adding strategic depth and fun to the game.

## ✨ Key Features

- 🎯 **Real-time Multiplayer**: WebSocket-based real-time gaming experience
- 📱 **Responsive Design**: Perfectly adapted for mobile and desktop browsers
- 🎨 **Modern UI**: Beautiful and smooth user interface
- 🤝 **Auto Match**: Support for creating and joining rooms
- 🎲 **Smart Rules**: Complete card type validation and game logic
- 🔐 **Hidden Teams**: Automatic team assignment based on Hearts Three distribution

## 🎯 Game Rules

### Cards and Rank

- **Cards**: 4 to 3, no Jokers
- **Rank** (from low to high): `4 < 5 < 6 < 7 < 8 < 9 < 10 < J < Q < K < A < 2 < 3`

### Card Types

- **Single**: One card
- **Pair**: Two cards of the same rank
- **Straight**: At least 5 consecutive cards, cannot include 3
- **Consecutive Pairs**: At least 3 consecutive pairs, cannot include 3
- **Bomb**: 3 cards of the same rank
- **Big Bomb**: 4 cards of the same rank

### Ranking Rules

- Straights/Consecutive Pairs are ranked by the highest card
- Bomb > Non-bomb card types
- Big Bomb > Regular Bomb

### Team Determination

- **One player has both Hearts Three** → 1v3 (one vs three)
- **Hearts Three split between two players (one each)** → 2v2 teams

**Note**: The frontend does not display team information. Players must deduce teammates through their hand and table cards.

### Playing Rules

- Playing order: Counter-clockwise
- First round: Must be played by the player who has Hearts Four (can be played alone or as part of a straight/consecutive pairs)
- Special rule: If a player has all four 4s, they can play any card in the first round
- Card type must match the current table cards and be higher in rank
- Bombs can beat any non-bomb card type
- 4-card bombs can beat 3-card bombs

### Round End

- If no one plays (all pass), the last player wins the round and collects cards
- The round restarts with the last player

### Win Conditions

**1v3 Mode**:
- Single player wins if they play all cards first

**2v2 Mode**:
- Team wins if both players finish their cards
- If one player from each team finishes, the remaining two continue
- After the game ends, remaining cards are assigned based on collected cards to determine the winning team

## 🛠️ Tech Stack

### Backend
- **Go** - Programming language
- **Gin** - Web framework
- **WebSocket** - Real-time communication
- **Redis** - Cache and session storage
- **PostgreSQL** - Database

### Frontend
- **Vue 3** - Frontend framework
- **Vite** - Build tool
- **WebSocket** - Real-time communication

## 📦 Installation & Running

### Requirements

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 7+

### Start Backend

```bash
cd backend

# Install dependencies
go mod download

# Run server
go run main.go
```

Backend runs on `:8080` by default.

### Start Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

Frontend runs on `http://localhost:5173` by default.

### Production Build

```bash
# Build frontend
cd frontend
npm run build

# Built files are in frontend/dist directory
```

## 📁 Project Structure

```
Hong3/
├── backend/              # Go backend
│   ├── api/             # API routes
│   ├── game/            # Game logic
│   ├── models/          # Data models
│   ├── websocket/       # WebSocket handling
│   └── main.go          # Entry point
├── frontend/            # Vue frontend
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── views/       # Page views
│   │   ├── store/       # State management
│   │   └── main.js      # Entry point
│   └── package.json
├── .gitignore
├── README.md            # Chinese docs
└── README-en.md         # English docs
```

## 🎮 Usage

1. **Enter Game**:
   - Open browser and visit frontend address
   - Enter player name and connect

2. **Create/Join Room**:
   - Click "Create Room" to create a new room
   - Or click "Refresh Room List" to view available rooms and join

3. **Ready Up**:
   - Wait for all 4 players to join
   - All players click "Ready" button
   - Game starts automatically

4. **Playing**:
   - Select cards from your hand (click cards)
   - Click "Play" to play selected cards
   - Or click "Pass" to skip the round
   - Use "Hint" button for suggestions

## 🔧 Configuration

### Backend Configuration

Default backend settings:
- Port: `8080`
- WebSocket path: `/ws`

Can be modified via environment variables or config file (to be implemented)

### Frontend Configuration

Default frontend settings:
- Dev server port: `5173`
- Backend API address: Automatically detected from current domain

To modify backend address, edit the `getBackendUrl` function in `frontend/src/store/gameStore.js`.

## 📱 Mobile Support

- ✅ Responsive layout for all screen sizes
- ✅ Landscape optimization for mobile landscape mode
- ✅ Touch scrolling for horizontal card browsing
- ✅ Support for IP access for local network multiplayer

## 🐛 Issues

If you encounter any problems or have suggestions, please submit an Issue.

## 📄 License

MIT License

## 👥 Contributing

Pull requests are welcome!

---

**Enjoy the game and good luck!** 🎉

