package library

import (
	"fmt"
	"strings"
)

type Character struct {
	Nom        string
	Classe     string
	Niveau     int
	PVMaxBase  int
	PVActuel   int
	PVMaxTotal int
	Inventaire map[string]int
	Skill      []string

	// Nombre maximal d'objets dans l'inventaire (T12). C'est un champ et non
	// une constante pour pouvoir l'augmenter plus tard (T18).
	CapaciteInventaire int
	// Pièces d'or du joueur (T13).
	Argent              int
	PotionGratuitePrise bool
	Equipement          Equipment
}

// Valeurs de départ du personnage.
const (
	capaciteInventaireDepart = 10
	argentDepart             = 100
)

func InitCharacter(nom string, classe string, niveau int, pvmax int, pvactuel int) Character {
	inventaire := map[string]int{
		ItemPotionDeVie: 3,
	}

	personnage := Character{
		Nom:                nom,
		Classe:             classe,
		Niveau:             niveau,
		PVMaxBase:          pvmax,
		PVActuel:           pvactuel,
		Inventaire:         inventaire,
		Skill:              []string{SortCoupDePoing},
		CapaciteInventaire: capaciteInventaireDepart,
		Argent:             argentDepart,
	}
	personnage.MettreAJourPvMax()
	return personnage
}

func (c Character) displayInfo() {
	fmt.Println("\n=== PERSONNAGE ===")
	fmt.Printf("Nom : %s\n", c.Nom)
	fmt.Printf("Classe : %s\n", c.Classe)
	fmt.Printf("Niveau : %d\n", c.Niveau)
	fmt.Printf("Pv : %d / %d\n", c.PVActuel, c.PVMaxTotal)
	fmt.Printf("Sorts : %s\n", strings.Join(c.Skill, ", "))
	fmt.Printf("Argent : %d pièces d'or\n", c.Argent)
	fmt.Println("\n=== INVENTAIRE ===")
	c.AccessInventory()
	fmt.Println("\n=== EQUIPEMENT ===")
	fmt.Printf("-- Tête --\n %s\n", c.Equipement.Tete.DescriptionEquipement())
	fmt.Printf("-- Torse --\n %s\n", c.Equipement.Torse.DescriptionEquipement())
	fmt.Printf("-- Pied --\n %s\n", c.Equipement.Pied.DescriptionEquipement())
}

// isDead vérifie si le personnage est mort (PV à 0 ou moins : les dégâts
// peuvent faire passer sous 0). Si oui, il ressuscite avec 50 % de ses PV
// max. Renvoie true s'il est mort, pour que l'appelant puisse s'arrêter.
func (c *Character) isDead() bool {
	if c.PVActuel > 0 {
		return false
	}
	fmt.Printf("%s est mort !\n", c.Nom)
	c.PVActuel = c.PVMaxTotal / 2
	fmt.Printf("%s ressuscite avec %d / %d PV.\n", c.Nom, c.PVActuel, c.PVMaxTotal)
	return true
}

func (c *Character) CalculPvAvecBonus() int {
	return c.Equipement.Tete.BonusPv + c.Equipement.Torse.BonusPv + c.Equipement.Pied.BonusPv + c.PVMaxBase
}

func (c *Character) MettreAJourPvMax() {
	c.PVMaxTotal = c.CalculPvAvecBonus()
	if c.PVActuel > c.PVMaxTotal {
		c.PVActuel = c.PVMaxTotal
	}
}
