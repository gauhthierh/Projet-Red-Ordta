package library

import "fmt"

// ==================================== //
// === Ce fichier gère les monstres === //
// ==================================== //

// Représente un monstre //
type Monster struct {
	Nom              string
	PvMax            int
	PvActuel         int
	Attaque          int
	Initiative       int
	ExperienceDonnee int
	OrMin            int
	OrMax            int
	Niveau           int
}

// Crée le gobelin d'entraînement avec ses statistiques de départ //
func InitGoblin(niveau int) Monster {
	bonus := niveau - 1
	pv := PvGobelin + GainPvGobelin*bonus
	return Monster{
		Nom:              "Gobelin d'entraînement",
		Niveau:           niveau,
		PvMax:            pv,
		PvActuel:         pv,
		Attaque:          AttaqueGobelin + GainAttaqueGobelin*bonus,
		ExperienceDonnee: ExperienceGobelin + GainExperienceGobelin*bonus,
		OrMin:            OrMinGobelin + GainOrMinGobelin*bonus,
		OrMax:            OrMaxGobelin + GainOrMaxGobelin*bonus,
	}
}

// Fait attaquer le gobelin selon son schéma de combat //
// Tous les 3 tours, il inflige le double de dégâts //
func (m Monster) GoblinPattern(c *Character, tour int) {
	degats := m.Attaque
	if tour%3 == 0 {
		fmt.Printf("%s charge son attaque !!!\n", m.Nom)
		degats *= 2
	}
	c.SubirDegats(degats)
	fmt.Printf("%s inflige %d dégâts à %s\n", m.Nom, degats, c.Nom)
	fmt.Printf("Pv de %s : %d / %d\n", c.Nom, c.PvActuel, c.PvMaxTotal)
}

// Retire des Pv au monstre sans descendre sous 0 //
func (m *Monster) SubirDegats(degats int) {
	m.PvActuel -= degats
	if m.PvActuel < 0 {
		m.PvActuel = 0
	}
}
