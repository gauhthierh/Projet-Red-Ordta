package library

import "fmt"

// =================================================== //
// === Ce fichier gère l'expérience et les niveaux === //
// =================================================== //

// Ajoute de l'expérience au joueur et le fait monter de niveau si besoin //
// Chaque niveau augmente ses statistiques selon sa classe //
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
