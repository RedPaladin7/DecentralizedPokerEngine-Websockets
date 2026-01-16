package game

import (
	"fmt"

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
	if len(g.getReadyPlayers()) >= 2 && g.currentStatus == GameStatusWaiting && g.currentDealerID == myID {
		g.StartNewHand()
	}
	return nil
}

func (g *Game) StartNewHand() {
	
}

