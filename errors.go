package main

import "errors"

var (
	ErrPlayerAlreadyInRoom          = errors.New("player is already in a room")
	ErrPlayerAlreadyInRequestedRoom = errors.New("player is already in the requested room")
	ErrPlayerNotInRoom              = errors.New("action cannot be done because player is not in room")
	ErrRoomDoesNotExist             = errors.New("room does not exist")
	ErrPlayerNotRoomLeader          = errors.New("player is not the room leader")
	ErrRoomIsLocked                 = errors.New("room is locked")
	ErrRoundDoesNotExist            = errors.New("round deos not exist")
)
