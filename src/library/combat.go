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

	for c.PVActuel > 0 && adversaire.PVActuel > 0 {
		fmt.Printf("\n=== TOUR %d ===\n", tour)

		if personnageCommence {
			c.CharacterTurn(&adversaire)

			if adversaire.PVActuel <= 0 {
				fmt.Printf("%s est vaincu !\n", adversaire.Nom)
				fmt.Println("Vous avez gagné l'entraînement, bien joué !")
				return
			}

			adversaire.GoblinPattern(c, tour)

			if c.PVActuel <= 0 {
				c.isDead()
				fmt.Println("Le combat d'entraînement est terminé.")
				return
			}
		} else {
			adversaire.GoblinPattern(c, tour)

			if c.PVActuel <= 0 {
				c.isDead()
				fmt.Println("Le combat d'entraînement est terminé.")
				return
			}

			c.CharacterTurn(&adversaire)

			if adversaire.PVActuel <= 0 {
				fmt.Printf("%s est vaincu !\n", adversaire.Nom)
				fmt.Println("Vous avez gagné l'entraînement, bien joué !")
				return
			}
		}

		tour++
	}
}
