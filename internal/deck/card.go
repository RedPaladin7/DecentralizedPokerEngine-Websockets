package deck 

import "fmt"

type Suit int 

const (
	Hearts Suit = iota 
	Diamonds 
	Clubs 
	Spades 
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	case Spades:
		return "Spades"
	}
	return "Unknown"
}

func (s Suit) Symbol() string {
	switch s {
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}

type Card struct {
	Suit Suit 
	Value int 
}

func NewCard(suit Suit, value int) Card {
	return Card{
		Suit: suit,
		Value: value,
	}
}

func NewCardFromByte(b byte) Card {
	value := int(b/4) + 1 
	suit := Suit(b % 4)
	return Card{
		Suit: suit,
		Value: value,
	}
}

func (c Card) ToByte() byte {
	return byte((c.Value-1)*4 + int(c.Suit))
}

func (c Card) String() string {
	var valueName string 
	switch c.Value {
	case 1:
		valueName = "A"
	case 11:
		valueName = "J"
	case 12:
		valueName = "Q"
	case 13:
		valueName = "K"
	default:
		valueName = fmt.Sprintf("%d", c.Value)
	}
	return valueName + c.Suit.Symbol()
}

func (c Card) FullName() string {
	var valueName string
	switch c.Value {
	case 14:
		valueName = "Ace"
	case 13:
		valueName = "King"
	case 12:
		valueName = "Queen"
	case 11:
		valueName = "Jack"
	default:
		valueName = fmt.Sprintf("%d", c.Value)
	}
	return valueName + " of " + c.Suit.String()
}

func (c Card) ToBytes() []byte {
	return []byte{c.ToByte()}
}

func (c Card) IsValid() bool {
	return c.Value >= 1 && c.Value <= 13 && c.Suit >= Hearts && c.Suit <= Spades
}

func (c Card) Compare(other Card) int {
	if c.Value > other.Value {
		return 1
	} else if c.Value < other.Value {
		return -1
	} 
	return 0
}