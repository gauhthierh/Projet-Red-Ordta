package main

import (
	"fmt"
)

func isDead(character *Character) {
	if character.pvactuel <= 0 {
		fmt.Println(character.nom, "est mort!")
		character.pvactuel = character.pvmax / 2
		fmt.Println(character.nom, "a maintenant", character.pvactuel, "points de vie")
	} else {
		fmt.Println(character.nom, "est vivant")
	}
}
