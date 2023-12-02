package server

import "sync"

type ChatServer struct {
	peers   map[*ChatPeer]bool
	muPeers sync.RWMutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		peers: map[*ChatPeer]bool{},
	}
}

func (chat *ChatServer) AddPeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	chat.peers[peer] = true
}

func (chat *ChatServer) RemovePeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	if _, exist := chat.peers[peer]; exist {
		peer.con.Close()
		delete(chat.peers, peer)
	}
}

func (chat *ChatServer) Run() {
}
