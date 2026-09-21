package library

import "fmt"

type Equipment struct {
	Tete  Stuff
	Torse Stuff
	Pied  Stuff
}

func (s Stuff) DescriptionEquipement() string {
	if s.Nom == "" {
		return "Aucun équipement"
	}
	return fmt.Sprintf("%s, Bonus : + %d Pv", s.Nom, s.BonusPv)
}

func TrouverEquipement(nom string) (Stuff, bool) {
	for _, s := range Armurerie {
		if s.Nom == nom {
			return s, true
		}
	}
	return Stuff{}, false
}

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
	fmt.Printf("Vos nouveaux Pv : %d / %d\n", c.PVActuel, c.PVMaxTotal)
}
