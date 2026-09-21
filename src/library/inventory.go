package library

import (
	"fmt"
	"sort"
)

func (c Character) AccessInventory() {
	vide := true
	// Même ordre que le menu d'inventaire : sortedItems ne garde que les
	// objets possédés et les trie par nom.
	for _, nom := range c.SortedItems() {
		fmt.Printf("- %s : %d\n", nom, c.Inventaire[nom])
		vide = false
	}
	if vide {
		fmt.Println("(inventaire vide)")
	}
}

// addInventory ajoute un exemplaire d'un objet dans l'inventaire.
// Renvoie false, sans rien ajouter, si l'inventaire est plein.
func (c *Character) AddInventory(item string) bool {
	if !c.VerifPlaceInventaire() {
		return false
	}
	c.Inventaire[item]++
	return true
}

func (c Character) TotalInventaire() int {
	var total int
	for _, quantite := range c.Inventaire {
		total += quantite
	}
	return total
}

// isInventoryFull indique si l'inventaire a atteint sa capacité. Le nombre
// d'objets est la somme des quantités (3 potions comptent pour 3), et non
// le nombre de sortes d'objets.
func (c Character) VerifPlaceInventaire() bool {
	return c.TotalInventaire() < c.CapaciteInventaire
}

// removeInventory retire un exemplaire d'un objet de l'inventaire.
// La clé est supprimée quand la quantité tombe à 0. Renvoie false si
// l'objet n'était pas dans l'inventaire.
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

// sortedItems renvoie les noms des objets possédés, triés par ordre
// alphabétique. Le parcours d'une map étant aléatoire en Go, ce tri
// garantit qu'un même numéro désigne toujours le même objet.
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

// inventoryMenu affiche l'inventaire numéroté et utilise l'objet choisi.
func (c *Character) InventoryMenu() {
	for {
		fmt.Println("\n=== INVENTAIRE ===")
		// La capacité est lue dans le personnage : l'affichage suivra
		// son augmentation (T18).
		fmt.Printf("Inventaire : %d / %d\n", c.TotalInventaire(), c.CapaciteInventaire)
		items := c.SortedItems()
		if len(items) == 0 {
			fmt.Println("(inventaire vide)")
		}
		for i, item := range items {
			fmt.Printf("%d. %s : %d\n", i+1, item, c.Inventaire[item])
		}
		fmt.Println("0. Retour")

		choice, ok := ReadChoice("Votre choix : ")

		if !ok || choice < 0 || choice > len(items) {
			if len(items) == 0 {
				fmt.Println("Choix invalide. Entrez 0 pour revenir au menu principal.")
			} else {
				fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d.\n", len(items))
			}
			continue
		}

		if choice == 0 {
			return
		}

		c.useItem(items[choice-1])
	}
}

func (c *Character) UpgradeInventorySlot() bool {
	if c.AugmentationInventaireUtilisee >= maxAugmentationsInventaire {
		return false
	}
	c.AugmentationInventaireUtilisee++
	c.CapaciteInventaire += bonusAugmentationInventaire
	return true
}
