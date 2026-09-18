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
		"potion de vie": 3,
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

func (c Character) accessInventory() {
	fmt.Println("Inventaire :")
	vide := true
	for nom, quantite := range c.Inventaire {
		if quantite > 0 {
			fmt.Printf("- %s : %d\n", nom, quantite)
			vide = false
		}
	}
	if vide {
		fmt.Println("(inventaire vide)")
	}
}

func main() {
	c1 := initCharacter("Guillaume", "Elfe", 1, 100, 40)
	c1.mainMenu()
}
