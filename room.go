package main

import (
	"math/rand"
	"sync"
	"time"
)

type RoomID string

type Room struct {
	sync.RWMutex
	roomInfo
}

type roomInfo struct {
	ID         RoomID
	Players    []string
	RoomLeader string
	// All players are in and rounds can start
	IsReady bool
	// The current round number. This is used so that the client can request
	// round info.
	RoundNumber int

	Round       Round
	LastUpdated time.Time
}

// roomInfo is a snapshot of the room's state. It does not send any information
// about the current round.
func (ri roomInfo) copy() roomInfo {
	return roomInfo{
		ID:          ri.ID,
		Players:     append([]string{}, ri.Players...),
		RoomLeader:  ri.RoomLeader,
		IsReady:     ri.IsReady,
		RoundNumber: ri.RoundNumber,
		LastUpdated: ri.LastUpdated,
	}
}

func (r *Room) IsLocked() bool {
	r.Lock()
	defer r.Unlock()
	return r.IsReady
}

func (r *Room) GetRoomInfo() roomInfo {
	r.Lock()
	defer r.Unlock()
	return r.roomInfo
}

func (r *Room) AddPlayer(player string) error {
	r.Lock()
	defer r.Unlock()
	if r.IsReady {
		return ErrRoomIsLocked
	}
	for _, p := range r.Players {
		if p == player {
			return ErrPlayerAlreadyInRequestedRoom
		}
	}
	r.Players = append(r.Players, player)
	return nil
}

func (r *Room) StartNewRound() {
	r.Lock()
	defer r.Unlock()
	chameleon := r.Players[rand.Intn(len(r.Players))]
	r.Round = NewRound(chameleon)
	r.RoundNumber++
}

func (r *Room) GetRound() Round {
	r.Lock()
	defer r.Unlock()
	return r.Round
}

func NewRoom(id string, leader string) Room {
	return Room{
		roomInfo: roomInfo{
			ID:          RoomID(id),
			Players:     []string{leader},
			RoomLeader:  leader,
			IsReady:     false,
			RoundNumber: 0,
			LastUpdated: time.Now().Round(time.Second),
		},
	}
}
