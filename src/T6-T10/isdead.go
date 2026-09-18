package main

import (
	"fmt"
)

func isDead(character *Character) {
	if character.PVActuel <= 0 {
		fmt.Println(character.Nom, "est mort!")
		character.PVActuel = character.PVMax / 2
		fmt.Println(character.Nom, "a maintenant", character.PVActuel, "points de vie")
	} else {
		fmt.Println(character.Nom, "est vivant")
	}
}
