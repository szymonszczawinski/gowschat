package chat

//
// import (
// 	"gowschat/server/api"
// 	"gowschat/server/chat/room"
// 	"gowschat/server/chatjson"
// 	"time"
// )
//
// func (chat *ChatServer) createMessageRoom(r *room.ChatRoom) (api.IMessageSerializable, error) {
// 	peers := []string{}
// 	for peer := range r.peers {
// 		peers = append(peers, peer.peerId)
// 	}
// 	switch p.PeerType {
// 	case api.PeerTypeJson:
// 		message := &chatjson.MessageRoomJson{
// 			RoomName: r.name,
// 			Peers:    peers,
// 		}
// 		return message, nil
// 	default:
// 		return nil, api.ErrUnsupporterPeerType
// 	}
// }
//
// func (chat *ChatServer) createMessageRoomList() (api.IMessageSerializable, error) {
// 	rooms := []string{}
// 	for room := range p.chat.rooms {
// 		rooms = append(rooms, room.name)
// 	}
// 	switch p.PeerType {
// 	case api.PeerTypeJson:
//
// 		message := &chatjson.MessageRoomListJson{
// 			Rooms: rooms,
// 		}
// 		return message, nil
// 	default:
// 		return nil, api.ErrUnsupporterPeerType
// 	}
// }
//
// func (chat *ChatServer) createMessageOut(m api.IMessage) (api.IMessageSerializable, error) {
// 	switch v := m.(type) {
// 	case api.IMessageIn:
// 		switch p.PeerType {
// 		case api.PeerTypeJson:
// 			messgeOut := &chatjson.MessageOutJson{
// 				Sent: time.Now(),
// 				MessageInJson: chatjson.MessageInJson{
// 					Message: v.GetMessage(),
// 					From:    v.GetFrom(),
// 				},
// 			}
//
// 			return messgeOut, nil
// 		default:
// 			return nil, api.ErrUnsupporterPeerType
// 		}
// 	default:
// 		return nil, api.ErrIncorrectMessageTypeIn
//
// 	}
// }
//
// func (chat *ChatServer) createMessageError(appError error) (api.IMessageSerializable, error) {
// 	switch p.PeerType {
// 	case api.PeerTypeJson:
// 		message := &chatjson.MessageErrorJson{
// 			Error: appError.Error(),
// 		}
// 		return message, nil
// 	default:
// 		return nil, api.ErrUnsupporterPeerType
// 	}
// }
