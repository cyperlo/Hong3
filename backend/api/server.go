package api

import (
	"log"
	"net/http"

	"github.com/chenhailong/hong3/auth"
	"github.com/chenhailong/hong3/websocket"
	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

// Server 表示API服务器
type Server struct {
	router *gin.Engine
	hub    *websocket.Hub
}

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
	EnableCompression: true,
}

// NewServer 创建一个新的API服务器
func NewServer() *Server {
	router := gin.Default()
	hub := websocket.NewHub()
	go hub.Run()

	server := &Server{
		router: router,
		hub:    hub,
	}

	server.setupRoutes()
	return server
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 允许跨域
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由
	s.router.POST("/api/register", s.handleRegister)
	s.router.POST("/api/login", s.handleLogin)
	s.router.GET("/api/me", s.handleGetMe)
	s.router.GET("/ws", s.handleWebSocket)
	s.router.GET("/api/health", s.handleHealth)
	s.router.GET("/api/rooms", s.handleGetRooms)
}

// Run 启动服务器
func (s *Server) Run(addr string) error {
	log.Printf("Starting server on %s", addr)
	return s.router.Run(addr)
}

// handleWebSocket 处理WebSocket连接
func (s *Server) handleWebSocket(c *gin.Context) {
	playerID := c.Query("player_id")
	playerName := c.Query("player_name")

	if playerID == "" || playerName == "" {
		log.Printf("WebSocket连接失败: 缺少玩家信息 (ID: %s, Name: %s)", playerID, playerName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少玩家ID或名称"})
		return
	}

	// 记录连接请求
	log.Printf("收到WebSocket连接请求 (ID: %s, Name: %s)", playerID, playerName)

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket连接升级失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法升级到WebSocket连接"})
		return
	}

	// 创建新的客户端
	client := websocket.NewClient(s.hub, conn, playerID, playerName)
	log.Printf("创建新客户端 (ID: %s, Name: %s)", playerID, playerName)

	// 注册客户端
	s.hub.Register <- client
	log.Printf("注册客户端到Hub (ID: %s, Name: %s)", playerID, playerName)

	// 启动客户端的读写循环
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("客户端读取循环panic: %v (ID: %s)", r, playerID)
			}
		}()
		client.ReadPump()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("客户端写入循环panic: %v (ID: %s)", r, playerID)
			}
		}()
		client.WritePump()
	}()

	log.Printf("WebSocket连接建立成功 (ID: %s, Name: %s)", playerID, playerName)
}

// handleHealth 处理健康检查
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleGetRooms 获取房间列表
func (s *Server) handleGetRooms(c *gin.Context) {
	rooms := s.hub.GetRooms()
	c.JSON(http.StatusOK, rooms)
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// handleRegister 处理注册
func (s *Server) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	store := auth.GetStore()
	user, err := store.Register(req.Username, req.Password, req.Name)
	if err != nil {
		if err == auth.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
		},
		"message": "注册成功",
	})
}

// handleLogin 处理登录
func (s *Server) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	store := auth.GetStore()
	user, token, err := store.Login(req.Username, req.Password)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
		},
		"message": "登录成功",
	})
}

// handleGetMe 获取当前用户信息
func (s *Server) handleGetMe(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		// 尝试从query参数获取
		token = c.Query("token")
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供token"})
		return
	}

	// 移除 "Bearer " 前缀（如果有）
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	store := auth.GetStore()
	user, err := store.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
		},
	})
}