package game

import (
	"encoding/json"
	"sync"

	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/blockchain"
	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/crypto"
	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/deck"
	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/protocol"
	"github.com/sirupsen/logrus"
)

const (
	SmallBlind = 10 
	BigBlind = 20 
)

type Game struct {
	lock 				sync.RWMutex
	listenAddr 			string 
	broadcastFunc 		BroadcastFunc
	playerStates 		map[string]*PlayerState
	rotationMap 		map[int]string 
	nextRotationID 		int 
	currentDealerID 	int 
	currentPlayerTurn 	int 
	currentStatus 		GameStatus 
	currentPot 			int 
	highestBet 			int 
	lastRaiserID 		int 
	lastRaiseAmount 	int 

	deckKeys 			*crypto.CardKeys
	foldedPlayerKeys 	map[string]*crypto.CardKeys
	revealedKeys 		map[string]*crypto.CardKeys
	currentDeck 		[][]byte 
	myHand 				[]deck.Card
	communityCards 		[]deck.Card

	sidePots 			[]SidePot 

	blockchain 			*blockchain.BlockchainClient 
	blockchainEnabled	bool 

	DisconnectHandler 	*DisconnectHandler
}

type BroadcastFunc func(data []byte, target ...string)

type SidePot struct {
	Amount int 
	Cap int 
	EligiblePlayers []string 
}

func NewGame(addr string, broadcast BroadcastFunc, bc *blockchain.BlockchainClient) *Game {
	keys, _ := crypto.GenerateCardKeys()
	g := &Game{
		listenAddr: addr,
		broadcastFunc: broadcast,
		playerStates: make(map[string]*PlayerState),
		rotationMap: make(map[int]string),
		currentStatus: GameStatusWaiting,
		deckKeys: keys,
		foldedPlayerKeys: make(map[string]*crypto.CardKeys),
		revealedKeys: make(map[string]*crypto.CardKeys),
		myHand: make([]deck.Card, 0, 2),
		communityCards: make([]deck.Card, 0, 5),
		sidePots: []SidePot{},
		blockchain: bc,
		blockchainEnabled: bc != nil,
	}

	g.DisconnectHandler = NewDisconnectHandler(g)

	go g.loop()
	return g
}

// TODO 
func (g *Game) loop() {}

func (g *Game) GetStatus() GameStatus {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.currentStatus
}

func (g *Game) setStatus(status GameStatus) {
	g.currentStatus = status
	logrus.Infof("Game status changed to: %s", status.String())
}

func (g *Game) PlayerCount() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.playerStates)
}

func (g *Game) ActivePlayerCount() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	count := 0
	for _, state := range g.playerStates {
		if state.IsActive {
			count++
		}
	}
	return count
}

func (g *Game) GetAllPlayers() []PlayerStateResponse {
	g.lock.RLock()
	defer g.lock.RUnlock()

	players := make([]PlayerStateResponse, 0)
	for i := 0; i < g.nextRotationID; i++ {
		addr, ok := g.rotationMap[i]
		if !ok {
			continue
		}

		state, ok := g.playerStates[addr]
		if !ok {
			continue
		}

		players = append(players, PlayerStateResponse{
			PlayerID:      state.ListenAddr,
			RotationID:    state.RotationID,
			Stack:         state.Stack,
			CurrentBet:    state.CurrentRoundBet,
			IsActive:      state.IsActive,
			IsFolded:      state.IsFolded,
			IsAllIn:       state.IsAllIn,
			IsReady:       state.IsReady,
			IsDealer:      state.RotationID == g.currentDealerID,
			IsCurrentTurn: state.RotationID == g.currentPlayerTurn,
		})
	}
	return players
}

func (g *Game) GetTableState(clientID string) TableStateResponse {
	g.lock.RLock()
	defer g.lock.RUnlock()

	myState, exists := g.playerStates[clientID]
	if !exists {
		return TableStateResponse{
			Status: g.currentStatus.String(),
		}
	}

	validActions := g.getValidActions(clientID)
	actionStrings := make([]string, len(validActions))
	for i, action := range validActions {
		actionStrings[i] = action.String()
	}

	myHandResp := make([]CardResponse, 0)
	if len(g.myHand) > 0 {
		myHandResp = make([]CardResponse, len(g.myHand))
		for i, card := range g.myHand {
			myHandResp[i] = CardResponse{
				Suit:    card.Suit.String(),
				Value:   card.Value,
				Display: card.String(),
			}
		}
	}

	communityCardResp := make([]CardResponse, len(g.communityCards))
	for i, card := range g.communityCards {
		communityCardResp[i] = CardResponse{
			Suit:    card.Suit.String(),
			Value:   card.Value,
			Display: card.String(),
		}
	}

	minRaise := g.highestBet + g.lastRaiseAmount
	if g.highestBet == 0 {
		minRaise = BigBlind
	}

	return TableStateResponse{
		Status:          g.currentStatus.String(),
		MyHand:          myHandResp,
		CommunityCards:  communityCardResp,
		Pot:             g.currentPot,
		HighestBet:      g.highestBet,
		MinRaise:        minRaise,
		ValidActions:    actionStrings,
		IsMyTurn:        myState.RotationID == g.currentPlayerTurn,
		MyStack:         myState.Stack,
		CurrentTurnID:   g.currentPlayerTurn,
		MyPlayerID:      myState.RotationID,
		DealerID:        g.currentDealerID,
		SmallBlind:      SmallBlind,
		BigBlind:        BigBlind,
	}
}

func (g *Game) HandleMessage(from string, msg *protocol.Message) error {
	switch msg.Type {
	case protocol.TypePlayerReady:
		return g.handleMessageReady(from)
	}
}

func (g *Game) handleMessageReady(from string) error {
	logrus.Infof("Player %s is ready", from)
	return g.SetPlayerReady(from)
}

func (g *Game) broadcast(data []byte, targets ...string) {
	if g.broadcastFunc != nil {
		g.broadcastFunc(data, targets...)
	}
}

func (g *Game) sendToPlayers(msgType protocol.MessageType, payload interface{}, targets ...string) error {
	msg, err := protocol.NewMessage(g.listenAddr, msgType, payload)
	if err != nil {
		return err
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	g.broadcast(data, targets...)
	return nil
}

func (g *Game) getOtherPlayers() []string {
	others := make([]string, 0)
	for addr := range g.playerStates {
		if addr != g.listenAddr {
			others = append(others, addr)
		}
	}
	return others
}

func (g *Game) getReadyPlayers() []string {
	ready := make([]string, 0)
	for addr, state := range g.playerStates {
		if state.IsReady {
			ready = append(ready, addr)
		}
	}
	return ready
}