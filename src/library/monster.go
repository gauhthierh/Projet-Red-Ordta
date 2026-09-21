package library

import "fmt"

/* Ce fichier définit les monstres, leur initialisation et leur comportement pendant les combats. */

/* La structure Monster contient les statistiques et les récompenses associées à un monstre. */
type Monster struct {
	Nom              string
	PVMax            int
	PVActuel         int
	Attaque          int
	Initiative       int
	ExperienceDonnee int
}

/* La fonction InitGoblin crée un gobelin d'entraînement avec ses statistiques et sa récompense d'expérience. */
func InitGoblin() Monster {
	gobelin := Monster{
		Nom:              "Gobelin d'entraînement",
		PVMax:            40,
		Attaque:          5,
		ExperienceDonnee: 40,
	}
	gobelin.PVActuel = gobelin.PVMax
	return gobelin
}

/* La méthode GoblinPattern fait attaquer le gobelin et double ses dégâts tous les trois tours. */
func (m Monster) GoblinPattern(c *Character, tour int) {
	degats := m.Attaque
	if tour%3 == 0 {
		fmt.Printf("%s charge son attaque !!!\n", m.Nom)
		degats *= 2
	}
	c.PVActuel -= degats
	if c.PVActuel < 0 {
		c.PVActuel = 0
	}
	fmt.Printf("%s inflige à %s %d de dégâts\n", m.Nom, c.Nom, degats)
	fmt.Printf("Pv de %s : %d / %d\n", c.Nom, c.PVActuel, c.PVMaxTotal)
}
