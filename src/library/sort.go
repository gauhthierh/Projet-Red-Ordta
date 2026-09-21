package library

/* Ce fichier gère l'utilisation des sorts pendant les combats.
   Il définit les dégâts des sorts et permet au personnage de choisir
   un sort parmi ceux qu'il a appris*/

import "fmt"

const (
	degatsCoupDePoing = 8
	degatsBouleDeFeu  = 18
)

/*La fonction ChoixSort affiche les sorts connus par le personnage et lui permet
  d'en utiliser un contre le monstre.
  Elle renvoie true si un sort est utilisé et false si le joueur revient.*/

func (c *Character) ChoixSort(m *Monster) bool {
	if len(c.Skill) == 0 {
		fmt.Println("Vous ne connaissez aucun sort.")
		return false
	}

	for {
		fmt.Println("\n=== SORTS ===")

		for indice, sort := range c.Skill {
			fmt.Printf("%d. %s\n", indice+1, sort)
		}

		fmt.Println("0. Retour")

		choix, ok := ReadChoice("Votre choix : ")

		if !ok || choix < 0 || choix > len(c.Skill) {
			fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d.\n", len(c.Skill))
			continue
		}

		if choix == 0 {
			return false
		}

		sortChoisi := c.Skill[choix-1]
		degats := 0

		switch sortChoisi {
		case SortCoupDePoing:
			degats = degatsCoupDePoing
		case SortBouleDeFeu:
			degats = degatsBouleDeFeu
		default:
			fmt.Println("Ce sort ne peut pas être utilisé.")
			continue
		}

		m.PVActuel -= degats

		if m.PVActuel < 0 {
			m.PVActuel = 0
		}

		fmt.Printf("%s utilise %s sur %s et lui inflige %d dégâts !\n", c.Nom, sortChoisi, m.Nom, degats)

		fmt.Printf("PV de %s : %d / %d\n", m.Nom, m.PVActuel, m.PVMax)
		return true
	}
}
