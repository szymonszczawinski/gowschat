package chat

//
// import (
// 	"testing"
// )
//
// func TestChatRoom_join(t *testing.T) {
// 	peerA := &ChatPeer{
// 		peerId: "peerA",
// 	}
// 	peerB := &ChatPeer{
// 		peerId: "peerB",
// 	}
// 	type fields struct {
// 		peers   map[*ChatPeer]bool
// 		creator *ChatPeer
// 		name    string
// 	}
// 	type args struct {
// 		p *ChatPeer
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{"creator can not join", fields{name: "roomA", creator: peerA, peers: map[*ChatPeer]bool{peerA: true}}, args{peerA}, true},
// 		{"can not join 1st time", fields{name: "roomA", creator: peerA, peers: map[*ChatPeer]bool{peerA: true}}, args{peerB}, false},
// 		{"can not join 2nd time", fields{name: "roomA", creator: peerA, peers: map[*ChatPeer]bool{peerA: true, peerB: true}}, args{peerB}, true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &ChatRoom{
// 				peers:   tt.fields.peers,
// 				creator: tt.fields.creator,
// 				name:    tt.fields.name,
// 			}
// 			if err := r.join(tt.args.p); (err != nil) != tt.wantErr {
// 				t.Errorf("ChatRoom.join() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
