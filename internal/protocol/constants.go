package protocol

const (
	GameVariantTexasHoldem = "TEXAS_HOLDEM"
)

const (
	ErrCodeInvalidMessage    = "INVALID_MESSAGE"
	ErrCodeInvalidAction     = "INVALID_ACTION"
	ErrCodeNotYourTurn       = "NOT_YOUR_TURN"
	ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
	ErrCodeGameNotStarted    = "GAME_NOT_STARTED"
	ErrCodePlayerNotFound    = "PLAYER_NOT_FOUND"
	ErrCodeAlreadyInGame     = "ALREADY_IN_GAME"
	ErrCodeGameFull          = "GAME_FULL"
	ErrCodeInternalError     = "INTERNAL_ERROR"
)

const (
	ActionFold = "fold"
	ActionCheck = "check"
	ActionCall = "call"
	ActionRaise = "raise"
	ActionAllIn = "all_in"
	ActionBet = "bet"
)

const (
	StateWaiting = "WAITING"
	StateDealing = "DEALING"
	StatePreFlop = "PREFLOP"
	StateFlop = "FLOP"
	StateTurn = "TURN"
	StateRiver = "RIVER"
	StateShowdown = "SHOWDOWN"
)

const (
	DefaultSmallBlind = 10 
	DefaultBigBlind = 20 
	DefaultStack = 1000
	DefaultMaxPlayers = 6
)