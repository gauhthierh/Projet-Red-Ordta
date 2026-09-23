package library

import (
	"fmt"
	"strings"
)

// ======================================== //
// === Ce fichier définit le personnage === //
// ======================================== //

// Représente les joueur //
type Character struct {
	Nom                            string
	Classe                         string
	Niveau                         int
	PvMaxBase                      int
	PvActuel                       int
	PvMaxTotal                     int
	Inventaire                     map[string]int
	Attaque                        int
	AttaquesPhysiques              []string
	Skill                          []string
	Initiative                     int
	ExperienceActuelle             int
	ExperienceMax                  int
	CapaciteInventaire             int
	Argent                         int
	PotionGratuitePrise            bool
	Equipement                     Equipment
	AugmentationInventaireUtilisee int
	ManaActuel                     int
	ManaMax                        int
	GainPvMax                      int
	GainManaMax                    int
	GainAttaque                    int
}

// Crée un personnage à l'aide de sa structure //
func InitCharacter(nom string, lvl int, cl Classe) Character {
	inventaire := map[string]int{
		ItemPotionDeVie: 3,
	}
	personnage := Character{
		Nom:                nom,
		Classe:             cl.Nom,
		Niveau:             lvl,
		PvMaxBase:          cl.PvMax,
		PvActuel:           cl.PvMax / 2,
		Inventaire:         inventaire,
		Attaque:            cl.Attaque,
		AttaquesPhysiques:  []string{AttaqueBasique},
		Skill:              []string{SortCoupDePoing},
		CapaciteInventaire: CapaciteInventaireDepart,
		Argent:             ArgentDepart,
		ExperienceActuelle: Experienceinitiale,
		ExperienceMax:      Experiencemaximale,
		ManaActuel:         cl.ManaMax,
		ManaMax:            cl.ManaMax,
		GainPvMax:          cl.GainPvMax,
		GainManaMax:        cl.GainManaMax,
		GainAttaque:        cl.GainAttaque,
	}
	personnage.MettreAJourPvMax()
	return personnage
}

// Affiche les infos du joueur //
func (c Character) displayInfo() {
	fmt.Println("\n=== PERSONNAGE ===")
	fmt.Printf("Nom : %s\n", c.Nom)
	fmt.Printf("Classe : %s\n", c.Classe)
	fmt.Printf("Niveau : %d\n", c.Niveau)
	fmt.Printf("Pv : %d / %d\n", c.PvActuel, c.PvMaxTotal)
	fmt.Printf("Mana : %d / %d\n", c.ManaActuel, c.ManaMax)
	fmt.Printf("Bonus d'attaque physique : +%d dégâts\n", c.Attaque)
	fmt.Printf("Attaques physiques : %s\n", strings.Join(c.AttaquesPhysiques, ", "))
	fmt.Printf("Sorts : %s\n", strings.Join(c.Skill, ", "))
	fmt.Printf("Argent : %d pièces d'or\n", c.Argent)
	fmt.Printf("Expérience : %d / %d\n", c.ExperienceActuelle, c.ExperienceMax)
	fmt.Println("\n=== INVENTAIRE ===")
	c.AccessInventory()
	fmt.Println("\n=== ÉQUIPEMENT ===")
	fmt.Printf("-- Tête --\n %s\n", c.Equipement.Tete.DescriptionEquipement())
	fmt.Printf("-- Torse --\n %s\n", c.Equipement.Torse.DescriptionEquipement())
	fmt.Printf("-- Pied --\n %s\n", c.Equipement.Pied.DescriptionEquipement())
}

// Ressuscite le joueur avec 50% de ses Pv en cas de décès //
func (c *Character) IsDead() bool {
	if c.PvActuel > 0 {
		return false
	}
	fmt.Printf("%s est mort !\n", c.Nom)
	c.PvActuel = c.PvMaxTotal / 2
	fmt.Printf("%s ressuscite avec %d / %d Pv.\n", c.Nom, c.PvActuel, c.PvMaxTotal)
	return true
}

// Renvoie les Pv max du joueur avec le bonus d'équipement //
func (c *Character) CalculPvAvecBonus() int {
	return c.Equipement.Tete.BonusPv + c.Equipement.Torse.BonusPv + c.Equipement.Pied.BonusPv + c.PvMaxBase
}

// Recalcule les Pv du joueur après un changement de base ou d'équipement //
func (c *Character) MettreAJourPvMax() {
	c.PvMaxTotal = c.CalculPvAvecBonus()
	if c.PvActuel > c.PvMaxTotal {
		c.PvActuel = c.PvMaxTotal
	}
}

// Retire des Pv au joueur sans descendre en négatif //
func (c *Character) SubirDegats(degats int) {
	c.PvActuel -= degats
	if c.PvActuel < 0 {
		c.PvActuel = 0
	}
}
