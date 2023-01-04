package main

import (
	"sync"

	"github.com/google/uuid"
)

type engine struct {
	sync.RWMutex
	Players     []string
	PlayerRooms map[string]RoomID
	Rooms       map[RoomID]*Room
}

func NewEngine() *engine {
	return &engine{
		Players:     []string{},
		PlayerRooms: map[string]RoomID{},
		Rooms:       map[RoomID]*Room{},
	}
}

func (e *engine) CreateRoom(player string) (RoomID, error) {
	if _, ok := e.PlayerRooms[player]; ok {
		// Player is already in a room
		return "", ErrPlayerAlreadyInRoom
	}
	id := uuid.New().String()[:4]
	e.Lock()
	defer e.Unlock()
	// Create a new room.
	room := NewRoom(id, player)
	e.Rooms[RoomID(id)] = &room
	e.PlayerRooms[player] = room.ID
	return room.ID, nil
}

func (e *engine) JoinRoom(player string, roomID string) error {
	e.Lock()
	defer e.Unlock()

	room, ok := e.Rooms[RoomID(roomID)]
	if !ok {
		return ErrRoomDoesNotExist
	}
	if err := room.AddPlayer(player); err != nil {
		return err
	}
	e.PlayerRooms[player] = room.ID
	return nil
}

func (e *engine) GetRoomInfo(player string, roomID RoomID) (roomInfo, error) {
	e.Lock()
	defer e.Unlock()

	if playerRoomID, ok := e.PlayerRooms[player]; !ok || roomID != playerRoomID {
		return roomInfo{}, ErrPlayerNotInRoom
	}

	room, ok := e.Rooms[roomID]
	if !ok {
		return roomInfo{}, ErrRoomDoesNotExist
	}
	return room.GetRoomInfo(), nil
}

func (e *engine) DestroyRoom(player string, roomID RoomID) error {
	e.Lock()
	defer e.Unlock()

	if playerRoomID, ok := e.PlayerRooms[player]; !ok || roomID != playerRoomID {
		return ErrPlayerNotInRoom
	}

	room, ok := e.Rooms[roomID]
	if !ok {
		return ErrRoomDoesNotExist
	}
	room.Lock()
	defer room.Unlock()

	if room.RoomLeader != player {
		return ErrPlayerNotRoomLeader
	}

	playersInRoom := room.Players
	// Remove all the players as being in the room
	for _, playerInRoom := range playersInRoom {
		delete(e.PlayerRooms, playerInRoom)
	}

	// Delete the room
	delete(e.Rooms, roomID)

	return nil
}

func (e *engine) StartNewRound(player string, roomID RoomID) (int, error) {
	e.Lock()
	defer e.Unlock()

	if playerRoomID, ok := e.PlayerRooms[player]; !ok || roomID != playerRoomID {
		return 0, ErrPlayerNotInRoom
	}

	room, ok := e.Rooms[roomID]
	if !ok {
		return 0, ErrRoomDoesNotExist
	}

	if room.RoomLeader != player {
		return 0, ErrPlayerNotRoomLeader
	}
	room.StartNewRound()
	return room.roomInfo.RoundNumber, nil
}

func (e *engine) GetRoundInfo(player string, roomID RoomID, round int) (Round, error) {
	e.Lock()
	defer e.Unlock()

	if playerRoomID, ok := e.PlayerRooms[player]; !ok || roomID != playerRoomID {
		return Round{}, ErrPlayerNotInRoom
	}

	room, ok := e.Rooms[roomID]
	if !ok {
		return Round{}, ErrRoomDoesNotExist
	}
	room.Lock()
	defer room.Unlock()

	if room.RoundNumber != round {
		return Round{}, ErrRoundDoesNotExist
	}

	return room.GetRound(), nil
}
