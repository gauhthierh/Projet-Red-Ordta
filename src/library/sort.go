package library

/* Ce fichier gère l'utilisation des sorts pendant les combats.
   Il définit les dégâts des sorts et permet au personnage de choisir
   un sort parmi ceux qu'il a appris*/

import (
	"fmt"
	"slices"
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
			_, coutManaSort, connu := InfosSort(sort)
			if !connu {
				fmt.Printf("%d. %s (indisponible)\n", indice+1, sort)
			} else {
				fmt.Printf("%d. %s (%d mana)\n", indice+1, sort, coutManaSort)
			}
		}
		fmt.Println("0. Retour")
		choix := ReadChoiceEntre("Votre choix : ", len(c.Skill))
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
			fmt.Printf("Impossible de lancer %s : il vous manque %d de mana.\n", sortChoisi, coutMana-c.ManaActuel)
			continue
		}
		c.ManaActuel -= coutMana
		m.SubirDegats(degats)
		fmt.Printf("%s utilise %s sur %s et lui inflige %d dégâts !\n", c.Nom, sortChoisi, m.Nom, degats)
		fmt.Printf("Pv de %s : %d / %d\n", m.Nom, m.PvActuel, m.PvMax)
		fmt.Printf("Mana de %s : %d / %d\n", c.Nom, c.ManaActuel, c.ManaMax)
		return true
	}
}

func InfosSort(sort string) (int, int, bool) {
	switch sort {
	case SortCoupDePoing:
		return DegatsCoupDePoing, CoutManaCoupDePoing, true
	case SortGrosseBouleDeFeu:
		return DegatsGrosseBouleDeFeu, CoutManaGrosseBouleDeFeu, true
	case SortLameDuDestin:
		return DegatsLameDuDestin, CoutManaLameDuDestin, true
	case SortEclatsDuGardien:
		return DegatsEclatsDuGardien, CoutManaEclatsDuGardien, true
	case SortFlecheDeLumiere:
		return DegatsFlecheDeLumiere, CoutManaFlecheDeLumiere, true
	case SortJugementDesGeants:
		return DegatsJugementDesGeants, CoutManaJugementDesGeants, true
	default:
		return 0, 0, false
	}
}

func SortDuLivre(objet string) (string, bool) {
	switch objet {
	case ItemLivreGrosseBouleDeFeu:
		return SortGrosseBouleDeFeu, true
	case ItemLivreLameDuDestin:
		return SortLameDuDestin, true
	case ItemLivreEclatsDuGardien:
		return SortEclatsDuGardien, true
	case ItemLivreFlecheDeLumiere:
		return SortFlecheDeLumiere, true
	case ItemLivreJugementDesGeants:
		return SortJugementDesGeants, true
	default:
		return "", false
	}
}

/* La méthode spellBook apprend le sort Boule de Feu et refuse de l'ajouter lorsqu'il est déjà connu. */
func (c *Character) SpellBook(sort string) bool {
	if slices.Contains(c.Skill, sort) {
		return false
	}
	c.Skill = append(c.Skill, sort)
	return true
}
