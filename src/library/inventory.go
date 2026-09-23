package library

import (
	"fmt"
	"sort"
)

/* Ce fichier gère l'affichage, le contenu, la capacité et l'utilisation de l'inventaire du personnage. */

/* La méthode AccessInventory affiche les objets possédés par le personnage dans un ordre stable. */
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

/* La méthode AddInventory ajoute un objet à l'inventaire et renvoie false lorsqu'il est plein. */
func (c *Character) AddInventory(item string) bool {
	if !c.VerifPlaceInventaire() {
		return false
	}
	c.Inventaire[item]++
	return true
}

/* La méthode TotalInventaire calcule le nombre total d'objets présents dans l'inventaire. */
func (c Character) TotalInventaire() int {
	var total int
	for _, quantite := range c.Inventaire {
		total += quantite
	}
	return total
}

/* La méthode VerifPlaceInventaire indique si l'inventaire possède encore une place disponible. */
func (c Character) VerifPlaceInventaire() bool {
	return c.TotalInventaire() < c.CapaciteInventaire
}

/* La méthode RemoveInventory retire un objet et supprime son entrée lorsque sa quantité atteint zéro. */
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

/* La méthode SortedItems renvoie les noms des objets possédés dans l'ordre alphabétique. */
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

/* La méthode InventoryMenu affiche l'inventaire et permet au joueur de choisir un objet à utiliser. */
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
		c.IsDead()
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
