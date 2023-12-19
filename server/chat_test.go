package server

import (
	"reflect"
	"testing"
)

func TestChatServer_createRoom(t *testing.T) {
	roomA := &ChatRoom{name: "roomA"}
	// roomB := &ChatRoom{name: "roomB"}
	creator := &ChatPeer{peerId: "creator"}
	type fields struct {
		rooms map[*ChatRoom]bool
	}
	type args struct {
		name    string
		creator *ChatPeer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"can add new room", fields{rooms: map[*ChatRoom]bool{}}, args{name: "roomA", creator: creator}, false},
		{"can not add same room", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomA", creator: creator}, true},
		{"can add different room", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomB", creator: creator}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &ChatServer{
				rooms: tt.fields.rooms,
			}
			if err := chat.createRoom(tt.args.name, tt.args.creator); (err != nil) != tt.wantErr {
				t.Errorf("ChatServer.createRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChatServer_joinRoom(t *testing.T) {
	peerA := &ChatPeer{peerId: "peerA"}
	peerB := &ChatPeer{peerId: "peerB"}
	roomA := &ChatRoom{name: "roomA", creator: peerA, peers: map[*ChatPeer]bool{peerA: true}}
	// roomB := &ChatRoom{name: "roomB"}

	type fields struct {
		rooms map[*ChatRoom]bool
	}
	type args struct {
		name   string
		joiner *ChatPeer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ChatRoom
		wantErr bool
	}{
		{"can join room", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomA", joiner: peerB}, roomA, false},
		{"can not join non existing room", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomB", joiner: peerB}, nil, true},
		{"can not join room same room twice", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomA", joiner: peerA}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &ChatServer{
				rooms: tt.fields.rooms,
			}
			got, err := chat.joinRoom(tt.args.name, tt.args.joiner)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatServer.joinRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChatServer.joinRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatServer_getRoom(t *testing.T) {
	peerA := &ChatPeer{peerId: "peerA"}
	roomA := &ChatRoom{name: "roomA", creator: peerA, peers: map[*ChatPeer]bool{peerA: true}}
	type fields struct {
		rooms map[*ChatRoom]bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ChatRoom
		wantErr bool
	}{
		{"can get room", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomA"}, roomA, false},
		{"can not get non existing room", fields{rooms: map[*ChatRoom]bool{roomA: true}}, args{name: "roomB"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &ChatServer{
				rooms: tt.fields.rooms,
			}
			got, err := chat.getRoom(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatServer.getRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChatServer.getRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}
