package library

/* Ce fichier gère l'utilisation des sorts pendant les combats.
   Il définit les dégâts des sorts et permet au personnage de choisir
   un sort parmi ceux qu'il a appris*/

import "fmt"

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
			_, coutManaSort, connu := InfosSort(sort)
			if !connu {
				fmt.Printf("%d. %s (indisponible)\n", indice+1, sort)
			} else {
				fmt.Printf("%d. %s (%d mana)\n", indice+1, sort, coutManaSort)
			}
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
		degats, coutMana, connu := InfosSort(sortChoisi)
		if !connu {
			fmt.Println("Ce sort ne peut pas être utilisé.")
			continue
		}
		if c.ManaActuel < coutMana {
			fmt.Printf("Impossible de lancer %s il vous manque %d de mana\n", sortChoisi, coutMana-c.ManaActuel)
			continue
		}
		m.PVActuel -= degats
		c.ManaActuel -= coutMana
		if m.PVActuel < 0 {
			m.PVActuel = 0
		}
		fmt.Printf("%s utilise %s sur %s et lui inflige %d dégâts !\n", c.Nom, sortChoisi, m.Nom, degats)
		fmt.Printf("PV de %s : %d / %d\n", m.Nom, m.PVActuel, m.PVMax)
		fmt.Printf("Mana de %s : %d / %d\n", c.Nom, c.ManaActuel, c.ManaMax)
		return true
	}
}

func InfosSort(sort string) (int, int, bool) {
	switch sort {
	case SortCoupDePoing:
		return degatsCoupDePoing, coutManaCoupDePoing, true
	case SortBouleDeFeu:
		return degatsBouleDeFeu, coutManaBouleDeFeu, true
	case SortLameDuDestin:
		return degatsLameDuDestin, coutManaLameDuDestin, true
	case SortEclatDuGardien:
		return degatsEclatDuGardien, coutManaEclatDuGardien, true
	case SortFlecheDeLumiere:
		return degatsFlecheDeLumiere, coutManaFlecheDeLumiere, true
	case SortJugementDesGeants:
		return degatsJugementDesGeants, coutManaJugementDesGeants, true
	default:
		return 0, 0, false
	}
}
