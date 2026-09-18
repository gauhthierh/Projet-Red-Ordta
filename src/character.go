package main

import "fmt"

type Character struct {
	Nom        string
	Classe     string
	Niveau     int
	PVMax      int
	PVActuel   int
	Inventaire map[string]int
}

func initCharacter(nom string, classe string, niveau int, pvmax int, pvactuel int) Character {
	inventaire := map[string]int{
		itemPotionDeVie: 3,
	}

	return Character{
		Nom:        nom,
		Classe:     classe,
		Niveau:     niveau,
		PVMax:      pvmax,
		PVActuel:   pvactuel,
		Inventaire: inventaire,
	}
}

func (c Character) displayInfo() {
	fmt.Printf("Nom : %s\n", c.Nom)
	fmt.Printf("Classe : %s\n", c.Classe)
	fmt.Printf("Niveau : %d\n", c.Niveau)
	fmt.Printf("Pv : %d / %d\n", c.PVActuel, c.PVMax)
	c.accessInventory()
}

func (c *Character) isDead() {
	if c.PVActuel <= 0 {
		fmt.Println(c.Nom, "est mort!")
		c.PVActuel = c.PVMax / 2
		fmt.Println(c.Nom, "a maintenant", c.PVActuel, "points de vie")
	} else {
		fmt.Println(c.Nom, "est vivant")
	}
}
