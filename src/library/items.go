package library

import (
	"fmt"
	"slices"
	"time"
)

/* Ce fichier définit les objets et les sorts puis applique leurs effets lorsqu'ils sont utilisés par le personnage. */

/* Ces constantes définissent les noms des objets et des emplacements utilisés dans le jeu. */
const (
	ItemPotionDeVie            = "Potion de vie"
	ItemPotionDePoison         = "Potion de poison"
	ItemLivreBouleDeFeu        = "Livre de Sort : Boule de Feu"
	ItemFourrureDeLoup         = "Fourrure de loup"
	ItemPeauDeTroll            = "Peau de Troll"
	ItemCuirDeSanglier         = "Cuir de sanglier"
	ItemPlumeDeCorbeau         = "Plume de corbeau"
	ItemChapeauAventurier      = "Chapeau de l'aventurier"
	ItemTuniqueAventurier      = "Tunique de l'aventurier"
	ItemBottesAventurier       = "Bottes de l'aventurier"
	EmplacementTete            = "tête"
	EmplacementTorse           = "torse"
	EmplacementPied            = "pied"
	ItemAugmentationInventaire = "Augmentation d'inventaire"
	ItemPotionDeMana           = "Potion de Mana"
)

/* La méthode useItem applique l'effet correspondant à l'objet sélectionné dans l'inventaire. */
func (c *Character) useItem(item string) {
	switch item {
	case ItemPotionDeVie:
		c.takePot()
	case ItemPotionDePoison:
		c.poisonPot()
	case ItemPotionDeMana:
		c.takePotMana()
	case ItemLivreBouleDeFeu:
		if !c.spellBook() {
			fmt.Printf("Vous connaissez déjà le sort %s : le livre reste dans l'inventaire.\n", SortBouleDeFeu)
			return
		}
		c.RemoveInventory(ItemLivreBouleDeFeu)
		fmt.Printf("Vous apprenez le sort %s !\n", SortBouleDeFeu)
	case ItemAugmentationInventaire:
		if !c.UpgradeInventorySlot() {
			fmt.Println("Vous avez déjà atteint la limite maximale d'améliorations d'inventaire")
			return
		}
		c.RemoveInventory(ItemAugmentationInventaire)
		fmt.Printf("Capacité d'inventaire : %d (+%d). Augmentations restantes : %d\n",
			c.CapaciteInventaire,
			bonusAugmentationInventaire,
			maxAugmentationsInventaire-c.AugmentationInventaireUtilisee)
	default:
		equipement, ok := TrouverEquipement(item)
		if ok {
			c.ChangerEquipement(equipement)
			return
		}
		fmt.Printf("%s n'a pas d'effet utilisable.\n", item)
	}
}

/* La méthode spellBook apprend le sort Boule de Feu et refuse de l'ajouter lorsqu'il est déjà connu. */
func (c *Character) spellBook() bool {
	if slices.Contains(c.Skill, SortBouleDeFeu) {
		return false
	}
	c.Skill = append(c.Skill, SortBouleDeFeu)
	return true
}

/* La méthode takePot consomme une potion de vie pour soigner le personnage sans dépasser ses points de vie maximum. */
func (c *Character) takePot() {
	if c.Inventaire[ItemPotionDeVie] > 0 && c.PVActuel >= c.PVMaxTotal {
		fmt.Println("Vous êtes déjà en pleine santé")
		return
	}
	if c.RemoveInventory(ItemPotionDeVie) {
		avant := c.PVActuel
		c.PVActuel += 50
		if c.PVActuel >= c.PVMaxTotal {
			c.PVActuel = c.PVMaxTotal
		}
		fmt.Printf("Vous buvez une %s (+%d PV)\n", ItemPotionDeVie, c.PVActuel-avant)
		fmt.Printf("PV : %d / %d\n", c.PVActuel, c.PVMaxTotal)
	} else {
		fmt.Println("Aucune potion de vie dans l'inventaire")
	}
}

func (c *Character) takePotMana() {
	if c.Inventaire[ItemPotionDeMana] > 0 && c.ManaActuel >= c.ManaMax {
		fmt.Println("Vous avez déjà votre mana au maximum")
		return
	}
	if c.RemoveInventory(ItemPotionDeMana) {
		avant := c.ManaActuel
		c.ManaActuel += 40
		if c.ManaActuel >= c.ManaMax {
			c.ManaActuel = c.ManaMax
		}
		fmt.Printf("Vous buvez une %s (+%d mana)\n", ItemPotionDeMana, c.ManaActuel-avant)
		fmt.Printf("Mana : %d / %d\n", c.ManaActuel, c.ManaMax)
	} else {
		fmt.Println("Aucune potion de mana dans l'inventaire")
	}
}

/* La méthode poisonPot consomme une potion de poison qui inflige dix dégâts par seconde pendant trois secondes ou jusqu'à la mort. */
func (c *Character) poisonPot() {
	if !c.RemoveInventory(ItemPotionDePoison) {
		fmt.Println("Aucune potion de poison dans l'inventaire")
		return
	}
	fmt.Printf("Vous buvez une %s...\n", ItemPotionDePoison)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.PVActuel -= 10
		if c.PVActuel < 0 {
			c.PVActuel = 0
		}
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", c.Nom, c.PVActuel, c.PVMaxTotal)
		if c.isDead() {
			fmt.Println("Le poison cesse de faire effet.")
			return
		}
	}
}
