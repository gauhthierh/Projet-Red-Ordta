package main

import (
	"fmt"
	"time"
)

func poisonPot(character *Character) {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		character.PVActuel -= 10
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", character.Nom, character.PVActuel, character.PVMax)
		isDead(character)
	}
}
