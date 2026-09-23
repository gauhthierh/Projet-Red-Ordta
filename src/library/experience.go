package library

/* Ce fichier gère le gain d'expérience et les montées de niveau
   du personnage après une victoire*/

import "fmt"

/* La fonction gagnerExperience ajoute des points d'expérience au personnage.
   Lorsque le maximum est atteint, le personnage gagne un niveau.
   L'expérience dépassant le maximum est conservée pour le niveau suivant.*/

func (c *Character) GagnerExperience(gain int) {
	fmt.Printf("%s gagne %d points d'expérience !\n", c.Nom, gain)
	c.ExperienceActuelle += gain
	for c.ExperienceActuelle >= c.ExperienceMax {
		c.ExperienceActuelle -= c.ExperienceMax
		c.Niveau++
		c.ExperienceMax = c.ExperienceMax * AugmentationExperienceNumerateur / AugmentationExperienceDenominateur
		c.PvMaxBase += c.GainPvMax
		c.MettreAJourPvMax()
		c.ManaMax += c.GainManaMax
		c.Attaque += c.GainAttaque

		fmt.Printf("%s passe au niveau %d !\n", c.Nom, c.Niveau)
		fmt.Println("Nouvelles Statistiques :")
		fmt.Printf("Pv : %d / %d\n", c.PvActuel, c.PvMaxTotal)
		fmt.Printf("Mana : %d / %d\n", c.ManaActuel, c.ManaMax)
		fmt.Printf("Bonus Attaque physique : +%d dégâts\n", c.Attaque)
	}
	fmt.Printf("Expérience : %d / %d\n", c.ExperienceActuelle, c.ExperienceMax)
}
