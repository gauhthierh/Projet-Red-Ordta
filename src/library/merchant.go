package library

import (
	"fmt"
)

// merchant affiche les objets disponibles chez le marchand.
func (c *Character) merchant() {
	for {
		fmt.Println("\n=== MARCHAND ===")
		fmt.Printf("Bourse : %d Po\n", c.Argent)
		for i, article := range Boutique {
			prix := c.PrixPour(article)
			if prix == 0 {
				fmt.Printf("%d. %s - Gratuit\n", i+1, article.Nom)
			} else {
				fmt.Printf("%d. %s - %d Po\n", i+1, article.Nom, prix)
			}
		}
		fmt.Println("0. Retour")

		choix, ok := ReadChoice("Votre choix : ")

		if !ok || choix < 0 || choix > len(Boutique) {
			fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d !\n", len(Boutique))
			continue
		}
		if choix == 0 {
			return
		}
		c.AchatMarchand(Boutique[choix-1])
	}
}

// buy ajoute l'objet acheté à l'inventaire, ou explique pourquoi c'est
// impossible : « Vous avez acheté » ne s'affiche que si l'ajout a eu lieu.
func (c *Character) buy(item string) {
	if !c.AddInventory(item) {
		fmt.Printf("Inventaire plein (%d / %d) : impossible d'ajouter %s.\n",
			c.TotalInventaire(), c.CapaciteInventaire, item)
		return
	}
	fmt.Println("Vous avez acheté :", item)
}

type Item struct {
	Nom  string
	Prix int
}

var Boutique = []Item{
	{Nom: ItemPotionDeVie, Prix: 3},
	{Nom: ItemPotionDePoison, Prix: 6},
	{Nom: ItemLivreBouleDeFeu, Prix: 25},
	{Nom: ItemFourrureDeLoup, Prix: 4},
	{Nom: ItemPeauDeTroll, Prix: 7},
	{Nom: ItemCuirDeSanglier, Prix: 3},
	{Nom: ItemPlumeDeCorbeau, Prix: 1},
	{Nom: ItemAugmentationInventaire, Prix: 30},
}

func (c *Character) PrixPour(i Item) int {
	if i.Nom == ItemPotionDeVie && !c.PotionGratuitePrise {
		return 0
	}
	return i.Prix
}

func (c *Character) AchatMarchand(i Item) {
	achat := c.PrixPour(i)
	if c.Argent < achat {
		fmt.Printf("Argent insuffisant, il manque %d Po\n", achat-c.Argent)
		return
	}
	if !c.AddInventory(i.Nom) {
		fmt.Println("Plus aucune place disponible dans l'inventaire")
		return
	}
	c.Argent -= achat
	if i.Nom == ItemPotionDeVie {
		c.PotionGratuitePrise = true
	}
	fmt.Printf("Vous avez acheté %s, il vous reste %d Po\n", i.Nom, c.Argent)
}
