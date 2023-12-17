package server

import (
	"errors"
	"sync"
)

var ErrRoomAlreadyJoined = errors.New("room already joined")

type ChatRoom struct {
	peers   map[*ChatPeer]bool
	creator *ChatPeer
	name    string
	muPeers sync.RWMutex
}

func NewChatRoom(name string, creator *ChatPeer) *ChatRoom {
	room := &ChatRoom{
		name:    name,
		peers:   map[*ChatPeer]bool{},
		creator: creator,
	}
	room.peers[creator] = true

	return room
}

func (r *ChatRoom) join(p *ChatPeer) error {
	r.muPeers.Lock()
	defer r.muPeers.Unlock()
	_, exist := r.peers[p]
	if exist {
		return ErrRoomAlreadyJoined
	}
	r.peers[p] = true
	return nil
}
