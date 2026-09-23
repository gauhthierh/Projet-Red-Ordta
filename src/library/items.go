package library

import (
	"fmt"
	"slices"
	"time"
)

/* Ce fichier définit les objets et les sorts puis applique leurs effets lorsqu'ils sont utilisés par le personnage. */

/* La méthode useItem applique l'effet correspondant à l'objet sélectionné dans l'inventaire. */
func (c *Character) useItem(item string) bool {
	switch item {
	case ItemPotionDeVie:
		return c.takePot()
	case ItemPotionDePoison:
		return c.poisonPot()
	case ItemPotionDeMana:
		return c.takePotMana()
	case ItemAugmentationInventaire:
		if !c.UpgradeInventorySlot() {
			fmt.Println("Vous avez déjà atteint la limite maximale d'améliorations d'inventaire")
			return false
		}
		c.RemoveInventory(ItemAugmentationInventaire)
		fmt.Printf("Capacité d'inventaire : %d (+%d). Augmentations restantes : %d\n",
			c.CapaciteInventaire,
			BonusAugmentationInventaire,
			MaxAugmentationsInventaire-c.AugmentationInventaireUtilisee)
	default:
		sort, estUnLivre := SortDuLivre(item)
		if estUnLivre {
			if !c.spellBook(sort) {
				fmt.Printf("Vous connaissez déjà le sort %s : le livre reste dans l'inventaire.\n", sort)
				return false
			}
			c.RemoveInventory(item)
			fmt.Printf("Vous apprenez le sort %s !\n", sort)
			return true
		}
		attaque, estUnManuel := AttaqueDuManuel(item)
		if estUnManuel {
			if !c.ApprentissageAttaque(attaque) {
				fmt.Printf("Vous connaissez déjà l'attaque %s : le manuel reste dans l'inventaire.\n", attaque)
				return false
			}
			c.RemoveInventory(item)
			fmt.Printf("Vous apprenez l'attaque %s !\n", attaque)
			return true
		}
		equipement, ok := TrouverEquipement(item)
		if ok {
			c.ChangerEquipement(equipement)
			return true
		}
		fmt.Printf("%s n'a pas d'effet utilisable.\n", item)
		return false
	}
	return true
}

/* La méthode spellBook apprend le sort Boule de Feu et refuse de l'ajouter lorsqu'il est déjà connu. */
func (c *Character) spellBook(sort string) bool {
	if slices.Contains(c.Skill, sort) {
		return false
	}
	c.Skill = append(c.Skill, sort)
	return true
}

/* La méthode takePot consomme une potion de vie pour soigner le personnage sans dépasser ses points de vie maximum. */
func (c *Character) takePot() bool {
	if c.Inventaire[ItemPotionDeVie] > 0 && c.PvActuel >= c.PvMaxTotal {
		fmt.Println("Vous êtes déjà en pleine santé")
		return false
	}
	if c.RemoveInventory(ItemPotionDeVie) {
		avant := c.PvActuel
		c.PvActuel += 50
		if c.PvActuel >= c.PvMaxTotal {
			c.PvActuel = c.PvMaxTotal
		}
		fmt.Printf("Vous buvez une %s (+%d Pv)\n", ItemPotionDeVie, c.PvActuel-avant)
		fmt.Printf("Pv : %d / %d\n", c.PvActuel, c.PvMaxTotal)
		return true
	} else {
		fmt.Println("Aucune potion de vie dans l'inventaire")
		return false
	}
}

func (c *Character) takePotMana() bool {
	if c.Inventaire[ItemPotionDeMana] > 0 && c.ManaActuel >= c.ManaMax {
		fmt.Println("Vous avez déjà votre mana au maximum")
		return false
	}
	if c.RemoveInventory(ItemPotionDeMana) {
		avant := c.ManaActuel
		c.ManaActuel += 40
		if c.ManaActuel >= c.ManaMax {
			c.ManaActuel = c.ManaMax
		}
		fmt.Printf("Vous buvez une %s (+%d mana)\n", ItemPotionDeMana, c.ManaActuel-avant)
		fmt.Printf("Mana : %d / %d\n", c.ManaActuel, c.ManaMax)
		return true
	} else {
		fmt.Println("Aucune potion de mana dans l'inventaire")
		return false
	}
}

/* La méthode poisonPot consomme une potion de poison qui inflige dix dégâts par seconde pendant trois secondes ou jusqu'à la mort. */
func (c *Character) poisonPot() bool {
	if !c.RemoveInventory(ItemPotionDePoison) {
		fmt.Println("Aucune potion de poison dans l'inventaire")
		return false
	}
	fmt.Printf("Vous buvez une %s...\n", ItemPotionDePoison)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.SubirDegats(10)
		fmt.Printf("%s a été empoisonné ! Pv : %d / %d\n", c.Nom, c.PvActuel, c.PvMaxTotal)
		if c.PvActuel <= 0 {
			fmt.Println("Le poison cesse de faire effet.")
			return true
		}
	}
	return true
}
