package game

type PlayerAction int 

const (
	PlayerActionFold PlayerAction = iota 
	PlayerActionCheck 
	PlayerActionCall 
	PlayerActionRaise 
	PlayerActionBet
)