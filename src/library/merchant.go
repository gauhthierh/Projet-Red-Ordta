package library

import (
	"fmt"
)

// =================================== //
// === Ce fichier gère le marchand === //
// =================================== //

// Représente un article du marchand //
type Item struct {
	Nom       string
	Prix      int
	NiveauMin int
}

// Affiche les articles, leur prix et leur niveau requis, puis achète l'article choisi //
func (c *Character) Merchant() {
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

// Liste des articles vendus par le marchand //
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

// Renvoie le prix d'un article pour ce joueur //
// La première potion de vie est offerte //
func (c *Character) PrixPour(i Item) int {
	if i.Nom == ItemPotionDeVie && !c.PotionGratuitePrise {
		return 0
	}
	return i.Prix
}

// Achète un article si le joueur a le niveau, l'argent et la place nécessaires //
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
