package main

type Round struct {
	// The secret word
	SecretWord string
	// The entries that everyone sees.
	Entries []string
	// The name of the player that is the chameleon.
	ChameleonPlayer string
}

func NewRound(chameleon string) Round {
	category := GetSmallCategory()
	return Round{
		SecretWord:      category.GetRandomEntry(),
		Entries:         category.Entries(),
		ChameleonPlayer: chameleon,
	}
}
