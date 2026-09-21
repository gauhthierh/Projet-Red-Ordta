package library

import (
	"fmt"
	"slices"
	"time"
)

// Noms des objets : définis une seule fois pour éviter les fautes de frappe
// dans les clés de l'inventaire.
const (
	ItemPotionDeVie       = "Potion de vie"
	ItemPotionDePoison    = "Potion de poison"
	ItemLivreBouleDeFeu   = "Livre de Sort : Boule de Feu"
	ItemFourrureDeLoup    = "Fourrure de loup"
	ItemPeauDeTroll       = "Peau de Troll"
	ItemCuirDeSanglier    = "Cuir de sanglier"
	ItemPlumeDeCorbeau    = "Plume de corbeau"
	ItemChapeauAventurier = "Chapeau de l'aventurier"
	ItemTuniqueAventurier = "Tunique de l'aventurier"
	ItemBottesAventurier  = "Bottes de l'aventurier"
	EmplacementTete       = "tête"
	EmplacementTorse      = "torse"
	EmplacementPied       = "pied"
)

// Noms des sorts.
const (
	SortCoupDePoing = "Coup de poing"
	SortBouleDeFeu  = "Boule de Feu"
)

// useItem applique l'effet d'un objet choisi dans l'inventaire.
func (c *Character) useItem(item string) {
	switch item {
	case ItemPotionDeVie:
		c.takePot()
	case ItemPotionDePoison:
		c.poisonPot()
	case ItemLivreBouleDeFeu:
		// spellBook vérifie si le sort est déjà connu AVANT que le livre
		// soit retiré : un livre inutile reste dans l'inventaire.
		if !c.spellBook() {
			fmt.Printf("Vous connaissez déjà le sort %s : le livre reste dans l'inventaire.\n", SortBouleDeFeu)
			return
		}
		c.RemoveInventory(ItemLivreBouleDeFeu)
		fmt.Printf("Vous apprenez le sort %s !\n", SortBouleDeFeu)
	default:
		equipement, ok := TrouverEquipement(item)
		if ok {
			c.ChangerEquipement(equipement)
			return
		}
		fmt.Printf("%s n'a pas d'effet utilisable.\n", item)
	}
}

// spellBook ajoute le sort Boule de Feu à la liste des sorts. Un même sort
// ne s'apprend qu'une fois : renvoie false s'il est déjà connu, sans rien
// modifier.
func (c *Character) spellBook() bool {
	if slices.Contains(c.Skill, SortBouleDeFeu) {
		return false
	}
	c.Skill = append(c.Skill, SortBouleDeFeu)
	return true
}

func (c *Character) takePot() {
	// À pleins PV, la potion serait gâchée : elle est refusée et reste
	// dans l'inventaire (même principe que le livre d'un sort déjà connu).
	if c.Inventaire[ItemPotionDeVie] > 0 && c.PVActuel >= c.PVMaxTotal {
		fmt.Println("Vous êtes déjà en pleine santé")
		return
	}
	if c.RemoveInventory(ItemPotionDeVie) {
		avant := c.PVActuel
		c.PVActuel += 50
		if c.PVActuel >= c.PVMaxTotal {
			c.PVActuel = c.PVMaxTotal
		}
		// Le gain annoncé est calculé après le plafond : à 70 / 100, +30.
		fmt.Printf("Vous buvez une %s (+%d PV)\n", ItemPotionDeVie, c.PVActuel-avant)
		fmt.Printf("PV : %d / %d\n", c.PVActuel, c.PVMaxTotal)
	} else {
		fmt.Println("Aucune potion dans l'inventaire")
	}
}

// poisonPot boit une potion de poison : 10 dégâts par seconde pendant 3 s.
// Le poison s'arrête si le personnage meurt.
func (c *Character) poisonPot() {
	if !c.RemoveInventory(ItemPotionDePoison) {
		fmt.Println("Aucune potion de poison dans l'inventaire")
		return
	}
	fmt.Printf("Vous buvez une %s...\n", ItemPotionDePoison)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.PVActuel -= 10
		if c.PVActuel < 0 {
			c.PVActuel = 0 // les PV ne s'affichent jamais en négatif
		}
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", c.Nom, c.PVActuel, c.PVMaxTotal)
		if c.isDead() {
			fmt.Println("Le poison cesse de faire effet.")
			return
		}
	}
}
