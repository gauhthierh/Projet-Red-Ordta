package library

/* Le fichier gère le déroulement d'un combat d'entrainement.
   Il organise les tours du personnage et du monstre, vérifie les pv
   après chaque attaque et met fin au combat à la mort de l'un des deux*/
import "fmt"

/* La fonction trainingFight prépare et lance un combat d'entrainement*/

func (c *Character) trainingFight() {
	adversaire := InitGoblin()
	tour := 1
	fmt.Printf("Début du combat contre %s \n", adversaire.Nom)

	for c.PVActuel > 0 && adversaire.PVActuel > 0 {
		fmt.Printf("\n=== TOUR %d ===\n", tour)
		c.CharacterTurn(&adversaire)

		if adversaire.PVActuel <= 0 {
			fmt.Printf("%s est vaincu\n", adversaire.Nom)
			fmt.Printf("Vous avez gagné l'entrainement, Bien joué !")
			return
		}
		adversaire.GoblinPattern(c, tour)

		if c.PVActuel <= 0 {
			c.isDead()
			fmt.Println("Le combat d'entrainement est terminé")
			return
		}
		tour++
	}
}
