package library

/*Ce fichier gère le marchand et les achats du personnage.
  Il contient la liste des objets vendus, leur prix et les vérifications
  nécessaires avant un achat : argent disponible et place dans l'inventaire*/

import (
	"fmt"
)

/*La fonction merchant affiche le menu du marchand et permet au joueur de choisir
  entre acheter un objet ou revenir au menu principal*/

func (c *Character) merchant() {
	for {
		fmt.Println("\n=== MARCHAND ===")
		fmt.Printf("Bourse : %d Po\n", c.Argent)
		for i, article := range Boutique {
			prix := c.PrixPour(article)
			if prix == 0 {
				fmt.Printf("%d. %s - Gratuit\n", i+1, article.Nom)
			} else if article.NiveauMin > 1 {
				fmt.Printf("%d. %s - %d Po (Niveau %d requis)\n", i+1, article.Nom, prix, article.NiveauMin)
			} else {
				fmt.Printf("%d. %s - %d Po\n", i+1, article.Nom, prix)
			}
		}
		fmt.Println("0. Retour")
		choix := ReadChoiceEntre("Votre choix : ", len(Boutique))
		if choix == 0 {
			return
		}
		c.AchatMarchand(Boutique[choix-1])
	}
}

/*La fonction buy ajoute gratuitement un objet à l'inventaire
  Si l'inventaire est plein l'ajout est refusé*/

func (c *Character) buy(item string) {
	if !c.AddInventory(item) {
		fmt.Printf("Inventaire plein (%d / %d) : impossible d'ajouter %s.\n",
			c.TotalInventaire(), c.CapaciteInventaire, item)
		return
	}
	fmt.Println("Vous avez acheté :", item)
}

/* Item représente un objet vendu par le marchand avec son nom et son prix*/

type Item struct {
	Nom       string
	Prix      int
	NiveauMin int
}

/* Boutique contient tous les objets disponibles chez le marchand*/

var Boutique = []Item{
	{Nom: ItemPotionDeVie, Prix: 3, NiveauMin: 1},
	{Nom: ItemPotionDePoison, Prix: 6, NiveauMin: 1},
	{Nom: ItemLivreGrosseBouleDeFeu, Prix: 25, NiveauMin: 1},
	{Nom: ItemFourrureDeLoup, Prix: 4, NiveauMin: 1},
	{Nom: ItemPeauDeTroll, Prix: 7, NiveauMin: 1},
	{Nom: ItemCuirDeSanglier, Prix: 3, NiveauMin: 1},
	{Nom: ItemPlumeDeCorbeau, Prix: 1, NiveauMin: 1},
	{Nom: ItemAugmentationInventaire, Prix: 30, NiveauMin: 1},
	{Nom: ItemPotionDeMana, Prix: 5, NiveauMin: 1},
	{Nom: ItemLivreLameDuDestin, Prix: 15, NiveauMin: 1},
	{Nom: ItemLivreEclatsDuGardien, Prix: 30, NiveauMin: 2},
	{Nom: ItemLivreFlecheDeLumiere, Prix: 30, NiveauMin: 3},
	{Nom: ItemLivreJugementDesGeants, Prix: 150, NiveauMin: 6},
	{Nom: ItemLivreAttaquePichenette, Prix: 20, NiveauMin: 1},
	{Nom: ItemLivreAttaqueClaquounette, Prix: 40, NiveauMin: 2},
	{Nom: ItemLivreAttaqueCoupsDePied, Prix: 70, NiveauMin: 3},
	{Nom: ItemLivreAttaqueMorsure, Prix: 120, NiveauMin: 5},
	{Nom: ItemLivreAttaqueUppercut, Prix: 200, NiveauMin: 7},
}

/* La fonction PrixPour détermine le prix d'un article pour le personnage*/

func (c *Character) PrixPour(i Item) int {
	if i.Nom == ItemPotionDeVie && !c.PotionGratuitePrise {
		return 0
	}
	return i.Prix
}

/* La fonction AchatMarchand tente d'acheter un article.
   L'achat est refusé si le personnage manque d'argent ou si son inventaire
   est plein. En cas de réussite, l'ajout à lieu et le prix est retiré*/

func (c *Character) AchatMarchand(i Item) {
	achat := c.PrixPour(i)
	if c.Niveau < i.NiveauMin {
		fmt.Printf("Vous n'avez pas le niveau requis pour l'achat : %d / %d\n", c.Niveau, i.NiveauMin)
		return
	}
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
