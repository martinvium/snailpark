package main

// const (
// 	AttrToughness = 1
// 	AttrPower     = 2
// 	AttrCost      = 3

// 	AttrType     = 4
// 	AttrLocation = 5
// 	AttrColor    = 6

// 	TagPlayerId    = 7
// 	TagTitle       = 8
// 	TagDescription = 9
// )

// // Seems like bullshit to me... maybe we should have different types, and
// // maybe we should have the variant type instead of different collections...
// type StackItem struct {
// 	Added     bool
// 	Removed   bool
// 	Card      *Card
// 	Key       int    // optional?
// 	Attribute int    // optional?
// 	Tag       string // optional?
// }

// // Maybe they should just be the types that we intend to send to the client e.g.
// type AddCard struct{}
// type RemoveCard struct{}
// type ChangeAttr struct{}
// type ChangeTag struct{}
// type StartGame struct{}

// // then we just make stack interface{}...?

// // also do we not want above to accept lists of things? but then should they
// // be lists for different types, then they all become roughly as complex as
// // the original.

// func DoStuff() {
// 	ChangeTag(g, AttrToughness-2)
// }

// func ChangeAttr(g *Game, id string, key int, value int) {
// 	card := CardById(g, entityId)
// 	if card == nil {
// 		return
// 	}

// 	if card.Attributes[key] == value {
// 		return
// 	}

// 	card.Attributes[key] = value

// 	g.stack = append(g.stack, &StackItem{card, key, Attribute: value})
// }

// func ChangeTag(g *Game, id string, key int, value string) {

// }
