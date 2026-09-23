package library

import "fmt"

// ======================================= //
// === Ce fichier gère les équipements === //
// ======================================= //

// Représente les emplacements d'équipement //
type Equipment struct {
	Tete  Stuff
	Torse Stuff
	Pied  Stuff
}

// Représente un équipement //
type Stuff struct {
	Nom         string
	Prix        int
	Materiaux   map[string]int
	BonusPv     int
	Emplacement string
}

// Renvoie le nom de l' équipement et son bonus de Pv //
func (s Stuff) DescriptionEquipement() string {
	if s.Nom == "" {
		return "Aucun équipement"
	}
	return fmt.Sprintf("%s, Bonus : + %d Pv", s.Nom, s.BonusPv)
}

// Cherche un équipement par son nom //
// Renvoie false si l'objet n'est pas un équipement //
func TrouverEquipement(nom string) (Stuff, bool) {
	for _, s := range Armurerie {
		if s.Nom == nom {
			return s, true
		}
	}
	return Stuff{}, false
}

// Équipe le joueur et met à jour ses Pv max //
// Retourne dans l'inventaire l'ancien équipement //
func (c *Character) ChangerEquipement(nouveau Stuff) {
	if nouveau.Emplacement != EmplacementTete &&
		nouveau.Emplacement != EmplacementTorse &&
		nouveau.Emplacement != EmplacementPied {
		fmt.Printf("%s ne peut pas être équipé : emplacement inconnu.\n", nouveau.Nom)
		return
	}
	if !c.RemoveInventory(nouveau.Nom) {
		fmt.Printf("Vous ne possédez pas %s.\n", nouveau.Nom)
		return
	}
	ancien := ""
	switch nouveau.Emplacement {
	case EmplacementTete:
		if c.Equipement.Tete.Nom != "" {
			ancien = c.Equipement.Tete.Nom
			c.AddInventory(ancien)
		}
		c.Equipement.Tete = nouveau
	case EmplacementTorse:
		if c.Equipement.Torse.Nom != "" {
			ancien = c.Equipement.Torse.Nom
			c.AddInventory(ancien)
		}
		c.Equipement.Torse = nouveau
	case EmplacementPied:
		if c.Equipement.Pied.Nom != "" {
			ancien = c.Equipement.Pied.Nom
			c.AddInventory(ancien)
		}
		c.Equipement.Pied = nouveau
	}
	c.MettreAJourPvMax()
	if ancien != "" {
		fmt.Printf("%s remplace %s\n", nouveau.Nom, ancien)
	} else {
		fmt.Printf("Vous équipez %s\n", nouveau.Nom)
	}
	fmt.Printf("Vos nouveaux Pv : %d / %d\n", c.PvActuel, c.PvMaxTotal)
}
