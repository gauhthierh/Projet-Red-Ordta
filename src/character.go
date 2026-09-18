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

// isDead vérifie si le personnage est mort (PV à 0 ou moins : les dégâts
// peuvent faire passer sous 0). Si oui, il ressuscite avec 50 % de ses PV
// max. Renvoie true s'il est mort, pour que l'appelant puisse s'arrêter.
func (c *Character) isDead() bool {
	if c.PVActuel > 0 {
		return false
	}
	fmt.Printf("%s est mort !\n", c.Nom)
	c.PVActuel = c.PVMax / 2
	fmt.Printf("%s ressuscite avec %d / %d PV.\n", c.Nom, c.PVActuel, c.PVMax)
	return true
}
