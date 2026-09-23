package library

import (
	"fmt"
	"sort"
)

// ============================================== //
// === Ce fichier gère l'inventaire du joueur === //
// ============================================== //

// Affiche les objets de l'inventaire et leur quantité //
func (c Character) AccessInventory() {
	vide := true
	for _, nom := range c.SortedItems() {
		fmt.Printf("- %s : %d\n", nom, c.Inventaire[nom])
		vide = false
	}
	if vide {
		fmt.Println("(inventaire vide)")
	}
}

// Ajoute un exemplaire d'un objet à l'inventaire //
// Renvoie false si l'inventaire est plein //
func (c *Character) AddInventory(item string) bool {
	if !c.VerifPlaceInventaire() {
		return false
	}
	c.Inventaire[item]++
	return true
}

// Renvoie le nombre total d'objets, en comptant chaque exemplaire //
func (c Character) TotalInventaire() int {
	var total int
	for _, quantite := range c.Inventaire {
		total += quantite
	}
	return total
}

// Renvoie true s'il reste au moins une place dans l'inventaire //
func (c Character) VerifPlaceInventaire() bool {
	return c.TotalInventaire() < c.CapaciteInventaire
}

// Retire un exemplaire d'un objet de l'inventaire //
// Renvoie false si l'objet n'y était pas //
func (c *Character) RemoveInventory(item string) bool {
	if c.Inventaire[item] <= 0 {
		return false
	}
	c.Inventaire[item]--
	if c.Inventaire[item] == 0 {
		delete(c.Inventaire, item)
	}
	return true
}

// Renvoie les noms des objets possédés //
// triés par ordre alphabétique //
func (c Character) SortedItems() []string {
	items := []string{}
	for item, quantite := range c.Inventaire {
		if quantite > 0 {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	return items
}

// Affiche l'inventaire hors combat et utilise l'objet choisi //
func (c *Character) InventoryMenu() {
	for {
		fmt.Println("\n=== INVENTAIRE ===")
		fmt.Printf("Inventaire : %d / %d\n", c.TotalInventaire(), c.CapaciteInventaire)
		items := c.SortedItems()
		if len(items) == 0 {
			fmt.Println("(inventaire vide)")
		}
		for i, item := range items {
			fmt.Printf("%d. %s : %d\n", i+1, item, c.Inventaire[item])
		}
		fmt.Println("0. Retour")
		choice := ReadChoiceEntre("Votre choix : ", len(items))
		if choice == 0 {
			return
		}
		c.UseItem(items[choice-1])
		if c.IsDead() {
			return
		}
	}
}

/* La méthode UpgradeInventorySlot augmente la capacité de l'inventaire tant que la limite d'améliorations n'est pas atteinte. */
func (c *Character) UpgradeInventorySlot() bool {
	if c.AugmentationInventaireUtilisee >= MaxAugmentationsInventaire {
		return false
	}
	c.AugmentationInventaireUtilisee++
	c.CapaciteInventaire += BonusAugmentationInventaire
	return true
}
