package library

/* Ce fichier gère le gain d'expérience et les montées de niveau
   du personnage après une victoire*/

import "fmt"

/* La fonction gagnerExperience ajoute des points d'expérience au personnage.
   Lorsque le maximum est atteint, le personnage gagne un niveau.
   L'expérience dépassant le maximum est conservée pour le niveau suivant.*/

func (c *Character) gagnerExperience(gain int) {
	fmt.Printf("%s gagne %d points d'expérience !\n", c.Nom, gain)

	c.ExperienceActuelle += gain

	for c.ExperienceActuelle >= c.ExperienceMax {
		c.ExperienceActuelle -= c.ExperienceMax
		c.Niveau++
		c.ExperienceMax += 50

		fmt.Printf("%s passe au niveau %d !\n", c.Nom, c.Niveau)
	}

	fmt.Printf("Expérience : %d / %d\n", c.ExperienceActuelle, c.ExperienceMax)
}
