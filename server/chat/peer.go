package chat

import (
	"errors"
	"fmt"
	"gowschat/server/api"
	"gowschat/server/chatjson"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type (
	PeerStatus string
	PeerType   string
)

const (
	PeerTypeJson  = "json"
	PeerTypeProto = "proto"
	PeerTypeWeb   = "web"
	PeerTypeUnset = "unset"

	PeerStatusOnline  = "online"
	PeerStatusOffline = "offline"
)

type (
	ChatPeer struct {
		server      *ChatServer
		con         *websocket.Conn
		parser      api.IParser
		serializer  api.ISerializer
		outgoing    chan api.MessageSerializable
		peerId      string
		rooms       map[*ChatRoom]bool
		PeerType    PeerType
		status      PeerStatus
		registerred bool
	}
	UserCredentials struct {
		email    string
		password string
	}
	ChatUser struct {
		*ChatPeer
		UserCredentials
	}
)

func NewChatPeer(chatServer *ChatServer, con *websocket.Conn) *ChatPeer {
	return &ChatPeer{
		server:      chatServer,
		con:         con,
		outgoing:    make(chan api.MessageSerializable),
		peerId:      uuid.NewString(),
		PeerType:    PeerTypeUnset,
		rooms:       map[*ChatRoom]bool{},
		status:      PeerStatusOnline,
		registerred: false,
	}
}

func NewUserCredentials(email, password string) UserCredentials {
	return UserCredentials{
		email:    email,
		password: password,
	}
}

func NewChatUser(p *ChatPeer, c UserCredentials) *ChatUser {
	return &ChatUser{
		UserCredentials: c,
		ChatPeer:        p,
	}
}

func (p *ChatPeer) ReadMessages() {
	defer func() {
		p.server.DisconnectPeer(p)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := p.con.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("ERROR :: SetReadDeadline:", err)
		return
	}
	p.con.SetReadLimit(512)
	// Configure how to handle Pong responses
	p.con.SetPongHandler(p.pongHandler)

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := p.readMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR :: reading message: %v", err)
			}
			p.status = PeerStatusOffline
			break // Break the loop to close conn & Cleanup
		}

		event, err := p.parser.ParseEvent(messageType, payload)
		if err != nil {
			log.Println("ERROR :: parseEvent:", err)
		} else {
			// Route the Event
			if err := p.server.RouteEvent(event, p); err != nil {
				log.Println("ERROR :: routeEvent", err)
			}
		}
	}
}

func (p *ChatPeer) WriteMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		p.server.DisconnectPeer(p)
		ticker.Stop()
	}()
	for {
		select {
		case messageOut, ok := <-p.outgoing:
			if !ok {
				if err := p.close(); err != nil {
					log.Println("ERROR :: connection closed", err)
				}
				return
			}
			event, err := p.createEvent(messageOut)
			if err != nil {
				log.Println("ERROR :: createEvent:", err)
				continue
			}
			data, err := event.Serialize()
			if err != nil {
				log.Println("ERROR :: Event.Serialize:", err)
				continue
			}

			if err := p.writeMessage(data); err != nil {
				log.Println("ERROR :: writeMessage:", err)
			}
		case <-ticker.C:
			log.Println("ping")
			// Send the Ping
			if err := p.con.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ERROR :: PING", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

func (p *ChatPeer) isRegisterred() bool {
	return p.registerred
}

// pongHandler is used to handle PongMessages for the Client
func (p *ChatPeer) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong", pongMsg)
	return p.con.SetReadDeadline(time.Now().Add(pongWait))
}

func (p *ChatPeer) initPeer(mesageType int) {
	if p.PeerType == PeerTypeUnset {
		switch mesageType {
		case websocket.TextMessage:
			p.PeerType = PeerTypeJson
			p.parser = &chatjson.ParserJson{}
		case websocket.BinaryMessage:
			p.PeerType = PeerTypeProto
		}
	}
}

func (p *ChatPeer) writeMessage(data []byte) error {
	if p.PeerType == PeerTypeJson {
		return p.con.WriteMessage(websocket.TextMessage, data)
	}
	if p.PeerType == PeerTypeProto {
		return p.con.WriteMessage(websocket.BinaryMessage, data)
	}
	return api.ErrUnsupporterPeerType
}

func (p *ChatPeer) readMessage() (int, []byte, error) {
	messageType, payload, err := p.con.ReadMessage()
	p.initPeer(messageType)
	return messageType, payload, err
}

func (p *ChatPeer) close() error {
	return p.con.WriteMessage(websocket.CloseMessage, nil)
}

func (p *ChatPeer) writeError(err error) error {
	if p.PeerType == PeerTypeJson {
		return p.con.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
	if p.PeerType == PeerTypeProto {
		return p.con.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
	}
	return api.ErrUnsupporterPeerType
}

func (p *ChatPeer) createMessageRoomList() (api.MessageSerializable, error) {
	rooms := []string{}
	for room := range p.server.rooms {
		rooms = append(rooms, room.name)
	}
	switch p.PeerType {
	case PeerTypeJson:

		message := &chatjson.MessageRoomListJson{
			Rooms: rooms,
		}
		return message, nil
	default:
		return nil, api.ErrUnsupporterPeerType
	}
}

func (p *ChatPeer) createMessageRoom(r *ChatRoom) (api.MessageSerializable, error) {
	peers := []string{}
	for peer := range r.peers {
		peers = append(peers, peer.peerId)
	}
	switch p.PeerType {
	case PeerTypeJson:
		message := &chatjson.MessageRoomJson{
			RoomName: r.name,
			Peers:    peers,
		}
		return message, nil
	default:
		return nil, api.ErrUnsupporterPeerType
	}
}

func (p *ChatPeer) createMessageError(appError error) (api.MessageSerializable, error) {
	switch p.PeerType {
	case PeerTypeJson:
		message := &chatjson.MessageErrorJson{
			Error: appError.Error(),
		}
		return message, nil
	default:
		return nil, api.ErrUnsupporterPeerType
	}
}

func (p *ChatPeer) createEvent(m api.MessageSerializable) (api.Event, error) {
	messagedata, err := m.Serialize()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if p.PeerType == PeerTypeJson {
		return m.CreateEvent(messagedata)
	}
	if p.PeerType == PeerTypeProto {
		return nil, errors.Join(api.ErrUnsupporterPeerType, fmt.Errorf("%v", p.PeerType))
	}
	return nil, api.ErrUnsupporterPeerType
}

func (p *ChatPeer) createMessageOut(m api.Message) (api.MessageSerializable, error) {
	switch v := m.(type) {
	case api.MessageIn:
		switch p.PeerType {
		case PeerTypeJson:
			messgeOut := &chatjson.MessageOutJson{
				Sent: time.Now(),
				MessageInJson: chatjson.MessageInJson{
					Message: v.GetMessage(),
					From:    v.GetFrom(),
				},
			}

			return messgeOut, nil
		default:
			return nil, api.ErrUnsupporterPeerType
		}
	default:
		return nil, api.ErrIncorrectMessageTypeIn

	}
}
