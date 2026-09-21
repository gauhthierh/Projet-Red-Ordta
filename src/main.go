package main

import (
	"fmt"
	"ordta/library"
)

/* Ce fichier constitue le point d'entrée du programme et lance la création du personnage puis le menu principal. */

/* La fonction main démarre le jeu en créant le personnage avant d'ouvrir le menu principal. */
func main() {
	fmt.Println("=== Création de personnage ===")
	c1 := library.CharacterCreation()
	c1.MainMenu()
}
