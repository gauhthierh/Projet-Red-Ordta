package library

import "fmt"

func (c *Character) forgeron() {
	for {
		fmt.Println("\n=== FORGERON ===")
		fmt.Printf("Bourse : %d Po\n", c.Argent)
		for i, equipement := range Armurerie {
			fmt.Printf("%d. %s - %d Po\n", i+1, equipement.Nom, equipement.Prix)
		}
		fmt.Println("0. Retour")

		choix, ok := ReadChoice("Votre choix : ")

		if !ok || choix < 0 || choix > len(Armurerie) {
			fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d !\n", len(Armurerie))
			continue
		}
		if choix == 0 {
			return
		}
		c.AchatForgeron(Armurerie[choix-1])
	}
}

type Stuff struct {
	Nom         string
	Prix        int
	Materiaux   map[string]int
	BonusPv     int
	Emplacement string
}

var Armurerie = []Stuff{
	{Nom: ItemChapeauAventurier, Prix: 5, Materiaux: map[string]int{ItemPlumeDeCorbeau: 1, ItemCuirDeSanglier: 1}, BonusPv: 10, Emplacement: EmplacementTete},
	{Nom: ItemTuniqueAventurier, Prix: 5, Materiaux: map[string]int{ItemFourrureDeLoup: 2, ItemPeauDeTroll: 1}, BonusPv: 25, Emplacement: EmplacementTorse},
	{Nom: ItemBottesAventurier, Prix: 5, Materiaux: map[string]int{ItemFourrureDeLoup: 1, ItemCuirDeSanglier: 1}, BonusPv: 15, Emplacement: EmplacementPied},
}

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
	// ne peut pas échouer : on vient de libérer au moins deux places
	for materiel, quantite := range s.Materiaux {
		for i := 1; i <= quantite; i++ {
			c.RemoveInventory(materiel)
		}
	}
	c.AddInventory(s.Nom)
	c.Argent -= achat
	fmt.Printf("Vous avez fabriqué %s, il vous reste %d Po\n", s.Nom, c.Argent)
}

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
