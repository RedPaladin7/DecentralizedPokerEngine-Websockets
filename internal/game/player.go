package game

import (
	"fmt"
	"sort"

	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/crypto"
	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/deck"
	"github.com/RedPaladin7/DecentralizedPokerEngine-Websockets.git/internal/protocol"
	"github.com/sirupsen/logrus"
)

type PlayerState struct {
	ListenAddr 			string 
	RotationID 			int 
	IsReady 			bool 
	IsActive 			bool 
	IsFolded 			bool 
	CurrentRoundBet 	int 
	IsAllIn 			bool 
	Stack 				int 
	TotalBetThisHand 	int 
}

type PlayerStateResponse struct {
	PlayerID 		string 	`json:"player_id"`
	RotationID 		int 	`json:"rotation_id"`
	Stack 			int 	`json:"stack"`
	CurrentBet 		int 	`json:"current_bet"`
	IsActive 		bool 	`json:"is_active"`
	IsFolded 		bool 	`json:"is_folded"`
	IsAllIn 		bool 	`json:"is_all_in"`
	IsReady 		bool 	`json:"is_ready"`
	IsDealer 		bool 	`json:"is_dealer"`
	IsCurrentTurn 	bool 	`json:"is_current_turn"`
}

type TableStateResponse struct {
	Status 			string 			`json:"status"`
	MyHand 			[]CardResponse 	`json:"my_hand"`
	CommunityCards 	[]CardResponse 	`json:"community_cards"`
	Pot 			int 			`json:"pot"`
	HighestBet 		int 			`json:"highest_bet"`
	MinRaise 		int 			`json:"min_raise"`
	ValidActions 	[]string 		`json:"valid_actions"`
	IsMyTurn 		bool 			`json:"is_my_turn"`
	MyStack 		int 			`json:"my_stack"`
	CurrentTurnID 	int 			`json:"current_turn_id"`
	MyPlayerID 		int 			`json:"my_player_id"`
	DealerID 		int 			`json:"dealer_id"`
	SmallBlind 		int 			`json:"small_blind"`
	BigBlind 		int 			`json:"big_blind"`
}

type CardResponse struct {
	Suit 	string 	`json:"suit"`
	Value 	int 	`json:"value"`
	Display string 	`json:"display"`
}

func (g *Game) SetPlayerReady(addr string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	state, ok := g.playerStates[addr]
	if !ok {
		return fmt.Errorf("player %s not found", addr)
	}
	if !state.IsReady {
		state.RotationID = g.nextRotationID
		g.rotationMap[state.RotationID] = addr 
		g.nextRotationID++ 
		state.IsReady = true 
		logrus.Infof("Player %s is ready (Rotation ID: %d)", addr, state.RotationID)
	}

	g.sendToPlayers(protocol.TypePlayerReady, protocol.PlayerReadyPayload{
		PlayerID: addr,
	}, g.getOtherPlayers()...)

	myID := g.playerStates[g.listenAddr].RotationID
	if len(g.getReadyPlayers()) >= 3 && g.currentStatus == GameStatusWaiting && g.currentDealerID == myID {
		g.StartNewHand()
	}
	return nil
}

func (g *Game) StartNewHand() {
	activeReadyPlayers := g.getReadyActivePlayers()
	if len(activeReadyPlayers) < 3 {
		g.setStatus(GameStatusWaiting)
		logrus.Warn("Not enough players to start a new hand")
		return 
	}
	logrus.Info("===Starting New Hand===")

	g.rotationMap = make(map[int]string)
	g.nextRotationID = 0
	g.myHand = make([]deck.Card, 0, 2)
	g.communityCards = make([]deck.Card, 0, 5)
	g.lastRaiseAmount = BigBlind 
	g.currentPot = 0
	g.highestBet = 0
	g.sidePots = []SidePot{}
	g.revealedKeys = make(map[string]*crypto.CardKeys)
	g.foldedPlayerKeys = make(map[string]*crypto.CardKeys)

	// recheck the sorting logic 
	sort.Strings(activeReadyPlayers)
	for _, addr := range activeReadyPlayers {
		state := g.playerStates[addr]
		state.RotationID = g.nextRotationID
		state.IsFolded = false 
		state.CurrentRoundBet = 0
		state.TotalBetThisHand = 0 
		state.IsAllIn = false 
		g.rotationMap[state.RotationID] = addr
		g.nextRotationID++
	}
	g.advanceDealer()
	g.postBlinds()
	g.setStatus(GameStatusDealing)
	g.InitiateShuffleAndDeal()
}

func (g *Game) postBlinds() {
	sbID := g.getNextActivePlayerID(g.currentDealerID)
	sbAddr :+ g.rotationMap[sbID]
	g.upd
}

func (g *Game) updatePlayerState(addr string, action PlayerAction, value int)

