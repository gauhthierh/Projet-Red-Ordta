package library

import "fmt"

func (c *Character) ChoixAttaquePhysique(m *Monster) bool {
	if len(c.AttaquesPhysiques) == 0 {
		fmt.Println("Vous ne connaissez aucune attaque physique.")
		return false
	}
	for {
		fmt.Println("\n=== ATTAQUES PHYSIQUES ===")
		for indice, attaque := range c.AttaquesPhysiques {
			degats, connu := c.CalculerDegatsPhysique(attaque)
			if !connu {
				fmt.Printf("%d. %s (indisponible)\n", indice+1, attaque)
			} else {
				fmt.Printf("%d. %s (%d dégâts)\n", indice+1, attaque, degats)
			}
		}
		fmt.Println("0. Retour")
		choix, ok := ReadChoice("Votre choix : ")
		if !ok || choix < 0 || choix > len(c.AttaquesPhysiques) {
			fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d.\n", len(c.AttaquesPhysiques))
			continue
		}
		if choix == 0 {
			return false
		}
		attaqueChoisi := c.AttaquesPhysiques[choix-1]
		degats, connu := c.CalculerDegatsPhysique(attaqueChoisi)
		if !connu {
			fmt.Println("Cette attaque ne peut pas être utilisée.")
			continue
		}
		m.PVActuel -= degats
		if m.PVActuel < 0 {
			m.PVActuel = 0
		}
		fmt.Printf("%s utilise %s sur %s et lui inflige %d dégâts !\n", c.Nom, attaqueChoisi, m.Nom, degats)
		fmt.Printf("PV de %s : %d / %d\n", m.Nom, m.PVActuel, m.PVMax)
		return true
	}
}

func InfosAttaquePhysique(attaque string) (int, bool) {
	switch attaque {
	case attaqueBasique:
		return degatsAttaqueBasique, true
	default:
		return 0, false
	}
}

func (c *Character) CalculerDegatsPhysique(attaque string) (int, bool) {
	degats, existe := InfosAttaquePhysique(attaque)
	if !existe {
		return 0, false
	} else {
		return degats + c.Attaque, existe
	}
}
