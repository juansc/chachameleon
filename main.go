package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	theGame := game{engine: NewEngine()}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/room", CORSMiddleware(), theGame.createRoom)
	r.GET("/room/:room_id/", CORSMiddleware(), theGame.getRoomInfo)
	r.POST("/room/:room_id/", CORSMiddleware(), theGame.joinRoom)
	r.POST("/room/:room_id/round", theGame.startNewRound)
	r.GET("/room/:room_id/round/:round_number/", theGame.getRoundInfo)
	r.GET("/reset", theGame.reset)
	r.Use(CORSMiddleware())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-methods, access-control-allow-origin, x-player")
		c.Header("Access-Control-expose-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-methods, access-control-allow-origin, x-player")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type game struct {
	engine *engine
}

type request struct {
	Token       string            `json:"token"`
	RequestType string            `json:"request_type"`
	Args        map[string]string `json:"args"`
	Player      string            `json:"player"`
}

// A leader can request the following actions on a room
// 1. Create a room
// 2. Ready the room
// 3.

type roomResponse struct {
	RoomID       string   `json:"room_id"`
	IsLeader     bool     `json:"is_leader"`
	OtherPlayers []string `json:"other_players"`
}

type roundResponse struct {
	Error string `json:"error"`
	// Only players will see this. A spectator cannot see this.
	IsChameleon bool `json:"is_chameleon"`
	// Everyone, including spectators, will see this.
	Words []string `json:"words"`
	// All the non-chameleon players will see this.
	SecretWord string `json:"secret_word"`
}

// People join and request to join a room.
// If the room doesn't exist, create a new room.
// A player can request to join a room.
// A player can request new round info. This will either
// fail because the round is not ready or there is one already ongoing,
// of it will return the round info.
// The master player will have a new round request. This will reset the room
// and then allow other players new round requests to get new info.
//
//
//

func (g *game) createRoom(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Player == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player is required"})
		return
	}

	roomID, err := g.engine.CreateRoom(req.Player)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "room created", "room_id": string(roomID)})
}

func (g *game) joinRoom(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("player", req.Player, "requested to join room ", c.Param("room_id"))

	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	if err := g.engine.JoinRoom(req.Player, roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "joined room"})
}

func (g *game) getRoomInfo(c *gin.Context) {
	reqPlayer := c.Request.Header.Get("X-Player")

	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	roomInfo, err := g.engine.GetRoomInfo(reqPlayer, RoomID(roomID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type roomInfoResponse struct {
		IsLeader     bool     `json:"is_leader"`
		Players      []string `json:"players"`
		CurrentRound int      `json:"current_round"`
		IsLocked     bool     `json:"is_locked"`
	}

	resp := roomInfoResponse{
		IsLeader:     roomInfo.RoomLeader == reqPlayer,
		Players:      roomInfo.Players,
		CurrentRound: roomInfo.RoundNumber,
		IsLocked:     roomInfo.IsReady,
	}

	c.JSON(http.StatusOK, &resp)
}

func (g *game) startNewRound(c *gin.Context) {
	reqPlayer := c.Request.Header.Get("X-Player")

	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	roundNumber, err := g.engine.StartNewRound(reqPlayer, RoomID(roomID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "new round started", "round_number": roundNumber})
}

func (g *game) getRoundInfo(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}
	roundNumberStr := c.Param("round_number")
	if roundNumberStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "round_number is required"})
		return
	}

	roundNumber, err := strconv.Atoi(roundNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "round_number must be a number"})
		return
	}

	roundInfo, err := g.engine.GetRoundInfo(req.Player, RoomID(roomID), roundNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := roundResponse{
		Words: roundInfo.Entries,
	}

	if req.Player == roundInfo.ChameleonPlayer {
		resp.IsChameleon = true
	} else {
		resp.SecretWord = roundInfo.SecretWord
	}
	c.JSON(http.StatusOK, &resp)
}

func (g *game) reset(_ *gin.Context) {
	g.engine.Reset()
}
