package library

import "fmt"

/* Ce fichier gère le forgeron, les équipements disponibles et leur fabrication à partir de matériaux. */

/* La méthode forgeron affiche l'armurerie et permet au personnage de choisir un équipement à fabriquer. */
func (c *Character) forgeron() {
	for {
		fmt.Println("\n=== FORGERON ===")
		fmt.Printf("Bourse : %d Po\n", c.Argent)
		for i, equipement := range Armurerie {
			fmt.Printf("%d. %s - %d Po\n", i+1, equipement.Nom, equipement.Prix)
			for materiau, quantite := range equipement.Materiaux {
				fmt.Printf("\t%s : %d (Possède : %d)\n", materiau, quantite, c.Inventaire[materiau])
			}
		}
		fmt.Println("0. Retour")
		choix := ReadChoiceEntre("Votre choix : ", len(Armurerie))
		if choix == 0 {
			return
		}
		c.AchatForgeron(Armurerie[choix-1])
	}
}

/* La structure Stuff représente un équipement avec son prix, ses matériaux, son bonus et son emplacement. */
type Stuff struct {
	Nom         string
	Prix        int
	Materiaux   map[string]int
	BonusPv     int
	Emplacement string
}

/* La variable Armurerie contient les équipements pouvant être fabriqués chez le forgeron. */
var Armurerie = []Stuff{
	{Nom: ItemChapeauAventurier, Prix: 5, Materiaux: map[string]int{ItemPlumeDeCorbeau: 1, ItemCuirDeSanglier: 1}, BonusPv: 10, Emplacement: EmplacementTete},
	{Nom: ItemTuniqueAventurier, Prix: 5, Materiaux: map[string]int{ItemFourrureDeLoup: 2, ItemPeauDeTroll: 1}, BonusPv: 25, Emplacement: EmplacementTorse},
	{Nom: ItemBottesAventurier, Prix: 5, Materiaux: map[string]int{ItemFourrureDeLoup: 1, ItemCuirDeSanglier: 1}, BonusPv: 15, Emplacement: EmplacementPied},
}

/* La méthode AchatForgeron vérifie les ressources et l'argent avant de fabriquer un équipement et de l'ajouter à l'inventaire. */
func (c *Character) AchatForgeron(s Stuff) {
	achat := s.Prix
	if c.Argent < achat {
		fmt.Printf("Argent insuffisant, il manque %d Po\n", achat-c.Argent)
		return
	}
	manques := c.MateriauxManquants(s)
	if len(manques) > 0 {
		for materiel, quantite := range manques {
			fmt.Printf("Il vous manque %d %s\n", quantite, materiel)
		}
		return
	}
	for materiel, quantite := range s.Materiaux {
		for i := 1; i <= quantite; i++ {
			c.RemoveInventory(materiel)
		}
	}
	c.AddInventory(s.Nom)
	c.Argent -= achat
	fmt.Printf("Vous avez fabriqué %s, il vous reste %d Po\n", s.Nom, c.Argent)
}

/* La méthode MateriauxManquants renvoie les matériaux et les quantités manquantes pour fabriquer un équipement. */
func (c *Character) MateriauxManquants(s Stuff) map[string]int {
	manques := make(map[string]int)
	for materiau, besoin := range s.Materiaux {
		possede := c.Inventaire[materiau]
		if possede < besoin {
			manques[materiau] = besoin - possede
		}
	}
	return manques
}
