package main

import (
	"fmt"
	"slices"
	"time"
)

// Noms des objets : définis une seule fois pour éviter les fautes de frappe
// dans les clés de l'inventaire.
const (
	itemPotionDeVie     = "Potion de vie"
	itemPotionDePoison  = "Potion de poison"
	itemLivreBouleDeFeu = "Livre de Sort : Boule de Feu"
)

// Noms des sorts.
const (
	sortCoupDePoing = "Coup de poing"
	sortBouleDeFeu  = "Boule de Feu"
)

// useItem applique l'effet d'un objet choisi dans l'inventaire.
func (c *Character) useItem(item string) {
	switch item {
	case itemPotionDeVie:
		c.takePot()
	case itemPotionDePoison:
		c.poisonPot()
	case itemLivreBouleDeFeu:
		// spellBook vérifie si le sort est déjà connu AVANT que le livre
		// soit retiré : un livre inutile reste dans l'inventaire.
		if !c.spellBook() {
			fmt.Printf("Vous connaissez déjà le sort %s : le livre reste dans l'inventaire.\n", sortBouleDeFeu)
			return
		}
		c.removeInventory(itemLivreBouleDeFeu)
		fmt.Printf("Vous apprenez le sort %s !\n", sortBouleDeFeu)
	default:
		fmt.Printf("%s n'a pas d'effet utilisable.\n", item)
	}
}

// spellBook ajoute le sort Boule de Feu à la liste des sorts. Un même sort
// ne s'apprend qu'une fois : renvoie false s'il est déjà connu, sans rien
// modifier.
func (c *Character) spellBook() bool {
	if slices.Contains(c.Skill, sortBouleDeFeu) {
		return false
	}
	c.Skill = append(c.Skill, sortBouleDeFeu)
	return true
}

func (c *Character) takePot() {
	// À pleins PV, la potion serait gâchée : elle est refusée et reste
	// dans l'inventaire (même principe que le livre d'un sort déjà connu).
	if c.Inventaire[itemPotionDeVie] > 0 && c.PVActuel >= c.PVMax {
		fmt.Println("Vous êtes déjà en pleine santé")
		return
	}
	if c.removeInventory(itemPotionDeVie) {
		avant := c.PVActuel
		c.PVActuel += 50
		if c.PVActuel >= c.PVMax {
			c.PVActuel = c.PVMax
		}
		// Le gain annoncé est calculé après le plafond : à 70 / 100, +30.
		fmt.Printf("Vous buvez une %s (+%d PV)\n", itemPotionDeVie, c.PVActuel-avant)
		fmt.Printf("PV : %d / %d\n", c.PVActuel, c.PVMax)
	} else {
		fmt.Println("Aucune potion dans l'inventaire")
	}
}

// poisonPot boit une potion de poison : 10 dégâts par seconde pendant 3 s.
// Le poison s'arrête si le personnage meurt.
func (c *Character) poisonPot() {
	if !c.removeInventory(itemPotionDePoison) {
		fmt.Println("Aucune potion de poison dans l'inventaire")
		return
	}
	fmt.Printf("Vous buvez une %s...\n", itemPotionDePoison)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.PVActuel -= 10
		if c.PVActuel < 0 {
			c.PVActuel = 0 // les PV ne s'affichent jamais en négatif
		}
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", c.Nom, c.PVActuel, c.PVMax)
		if c.isDead() {
			fmt.Println("Le poison cesse de faire effet.")
			return
		}
	}
}
