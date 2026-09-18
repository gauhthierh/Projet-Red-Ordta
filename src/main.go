package main

import "fmt"

func main() {
	fmt.Println("=== Création de personnage ===")
	c1 := characterCreation()
	c1.mainMenu()
}
