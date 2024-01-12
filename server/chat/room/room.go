package room

import (
	"gowschat/server/api"
	"sync"
)

type ChatRoom struct {
	peers   map[api.IChatPeer]bool
	creator api.IChatPeer
	Name    string
	muPeers sync.RWMutex
}

func NewChatRoom(name string, creator api.IChatPeer) *ChatRoom {
	room := &ChatRoom{
		Name:    name,
		peers:   map[api.IChatPeer]bool{},
		creator: creator,
	}
	room.peers[creator] = true

	return room
}

func (r *ChatRoom) Join(p api.IChatPeer) error {
	r.muPeers.Lock()
	defer r.muPeers.Unlock()
	_, exist := r.peers[p]
	if exist {
		return api.ErrRoomAlreadyJoined
	}
	r.peers[p] = true
	return nil
}
