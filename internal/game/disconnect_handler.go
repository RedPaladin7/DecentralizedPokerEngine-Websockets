package game

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DisconnectTimeout = 5 * time.Minute
)

type DisconnectHandler struct {
	game 				*Game
	disconnectTimers 	map[string]*time.Timer
	reconnectChannels 	map[string]chan bool 
	mu 					sync.RWMutex
	logger 				*logrus.Logger 
}

func NewDisconnectHandler(game *Game) *DisconnectHandler {
	return &DisconnectHandler{
		game: 				game,
		disconnectTimers: 	make(map[string]*time.Timer),
		reconnectChannels: 	make(map[string]chan bool),
		logger: 			logrus.New(),
	}
}