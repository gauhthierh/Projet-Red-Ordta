package library

/* Le fichier gère le déroulement d'un combat d'entrainement.
   Il organise les tours du personnage et du monstre, vérifie les pv
   après chaque attaque et met fin au combat à la mort de l'un des deux*/
import (
	"fmt"
	"math/rand"
)

/* La fonction trainingFight prépare et lance un combat d'entrainement*/

func (c *Character) trainingFight() {
	adversaire := InitGoblin()
	tour := 1
	c.Initiative = rand.Intn(10) + 1
	adversaire.Initiative = rand.Intn(10) + 1
	fmt.Printf("\nDébut du combat contre %s !\n", adversaire.Nom)
	fmt.Printf("Initiative de %s : %d\n", c.Nom, c.Initiative)
	fmt.Printf("Initiative de %s : %d\n", adversaire.Nom, adversaire.Initiative)
	personnageCommence := c.Initiative >= adversaire.Initiative
	if personnageCommence {
		fmt.Printf("%s commence le combat !\n", c.Nom)
	} else {
		fmt.Printf("%s commence le combat !\n", adversaire.Nom)
	}
	for {
		fmt.Printf("\n=== TOUR %d ===\n", tour)
		if personnageCommence {
			c.CharacterTurn(&adversaire)
			if c.FinCombat(&adversaire) {
				return
			}
			adversaire.GoblinPattern(c, tour)
			if c.FinCombat(&adversaire) {
				return
			}
		} else {
			adversaire.GoblinPattern(c, tour)
			if c.FinCombat(&adversaire) {
				return
			}
			c.CharacterTurn(&adversaire)
			if c.FinCombat(&adversaire) {
				return
			}
		}
		tour++
	}
}

func (c *Character) FinCombat(m *Monster) bool {
	if m.PvActuel <= 0 {
		fmt.Printf("%s est vaincu !\n", m.Nom)
		fmt.Println("Vous avez gagné l'entraînement, bien joué !")
		c.GagnerExperience(m.ExperienceDonnee)
		return true
	}
	if c.PvActuel <= 0 {
		c.IsDead()
		fmt.Println("Le combat d'entraînement est terminé.")
		return true
	}
	return false
}
