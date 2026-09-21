package library

import "fmt"

type Monster struct {
	Nom      string
	PVMax    int
	PVActuel int
	Attaque  int
}

func InitGoblin() Monster {
	gobelin := Monster{
		Nom:     "Gobelin d'entraînement",
		PVMax:   40,
		Attaque: 5,
	}
	gobelin.PVActuel = gobelin.PVMax
	return gobelin
}

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
