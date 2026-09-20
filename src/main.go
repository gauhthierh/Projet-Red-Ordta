package main

import (
	"fmt"
	"ordta/library"
)

func main() {
	fmt.Println("=== Création de personnage ===")
	c1 := library.CharacterCreation()
	c1.MainMenu()
}
