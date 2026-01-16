package protocol

import (
	"encoding/json"
	"time"
)

type MessageType string 

const (
	TypeHandshake 		MessageType = "handshake"
	TypePeerList 		MessageType = "peer_list"
	TypePlayerAction 	MessageType = "player_action"
	TypePlayerReady 	MessageType = "player_ready"
	TypeEncDeck 		MessageType = "enc_deck"
	TypeGameState 		MessageType = "game_state"
	TypeShuffleStatus 	MessageType = "shuffle_status"
	TypeGetRPC			MessageType = "get_rpc"
	TypeRPCResponse 	MessageType = "rpc_response"
	TypeRevealKeys 		MessageType = "reveal_keys"
	TypeShowdownResult 	MessageType = "showdown_result"
	TypeError 			MessageType = "error"
	TypePing 			MessageType = "ping"
	TypePong 			MessageType = "pong"
)

type Message struct {
	Type 		MessageType    `json:"type"`
	From 		string         `json:"from"`
	Payload 	json.RawMessage `json:"payload"`
	TimeStamp 	time.Time `json:"timestamp"`
}

func NewMessage(from string, msgType MessageType, payload interface{}) (*Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err 
	}
	return &Message{
		Type: 		msgType,
		From: 		from,
		Payload: 	data,
		TimeStamp: 	time.Now(),
	}, nil
}

type PlayerReadyPayload struct {
	PlayerID string `json:"player_id"`
}