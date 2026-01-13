package deck

import "github.com/chehsunliu/poker"

func EvaluateBestHand(holeCards, communityCards []Card) (int32, string) {
	var cardStrings []string 
	for _, c := range append(holeCards, communityCards...){
		cardStrings = append(cardStrings, cardToString(c))

	}
	cards := make([]poker.Card, len(cardStrings))
	for i, s := range cardStrings {
		cards[i] = poker.NewCard(s)
	}
	result := poker.Evaluate(cards)
	return int32(result), getHandType(result)
}

func cardToString(c Card) string {
	ranks := map[int]string{
		1:  "A",
		11: "J",
		12: "Q",
		13: "K",
		10: "T",
		9: "9",
		8: "8",
		7: "7",
		6: "6",
		5: "5",
		4: "4",
		3: "3",
		2: "2",
	}
	suits := map[string]string{
		"Spades": "s", "Hearts": "h", "Clubs": "c", "Diamonds": "d",
	}
	return ranks[c.Value] + suits[c.Suit.String()]
}

func getHandType(rank int32) string {
	return poker.RankString(rank)
}