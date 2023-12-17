package server

import "sync"

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
