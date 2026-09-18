package main

import (
	"fmt"
	"sort"
)

func (c Character) accessInventory() {
	fmt.Println("Inventaire :")
	vide := true
	// Même ordre que le menu d'inventaire : sortedItems ne garde que les
	// objets possédés et les trie par nom.
	for _, nom := range c.sortedItems() {
		fmt.Printf("- %s : %d\n", nom, c.Inventaire[nom])
		vide = false
	}
	if vide {
		fmt.Println("(inventaire vide)")
	}
}

// addInventory ajoute un exemplaire d'un objet dans l'inventaire.
func (c *Character) addInventory(item string) {
	c.Inventaire[item]++
}

// removeInventory retire un exemplaire d'un objet de l'inventaire.
// La clé est supprimée quand la quantité tombe à 0. Renvoie false si
// l'objet n'était pas dans l'inventaire.
func (c *Character) removeInventory(item string) bool {
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
func (c Character) sortedItems() []string {
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
func (c *Character) inventoryMenu() {
	for {
		fmt.Println("\n=== INVENTAIRE ===")
		items := c.sortedItems()
		if len(items) == 0 {
			fmt.Println("(inventaire vide)")
		}
		for i, item := range items {
			fmt.Printf("%d. %s : %d\n", i+1, item, c.Inventaire[item])
		}
		fmt.Println("0. Retour")

		choice, ok := readChoice("Votre choix : ")

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
