package main

import (
	"slices"
	"testing"
)

// nouveauPersonnage crée un personnage de test sans passer par la saisie
// clavier : les tests appellent initCharacter directement.
func nouveauPersonnage(pvMax, pvActuel int) Character {
	return initCharacter("Test", "Elfe", 1, pvMax, pvActuel)
}

// T2, T10, T12, T13 : initCharacter donne au personnage tout ce que les
// tâches demandent au départ. characterCreation passe par initCharacter,
// donc elle en hérite.
func TestInitCharacterValeursDeDepart(t *testing.T) {
	c := initCharacter("Test", "Nain", 1, 120, 60)

	if c.Niveau != 1 || c.PVMax != 120 || c.PVActuel != 60 {
		t.Errorf("niveau/PV : obtenu %d, %d/%d", c.Niveau, c.PVActuel, c.PVMax)
	}
	if c.Inventaire[itemPotionDeVie] != 3 {
		t.Errorf("3 potions de vie attendues, obtenu %d", c.Inventaire[itemPotionDeVie])
	}
	if !slices.Equal(c.Skill, []string{sortCoupDePoing}) {
		t.Errorf("seul sort attendu : Coup de poing, obtenu %v", c.Skill)
	}
	if c.CapaciteInventaire != 10 {
		t.Errorf("capacité de départ attendue 10, obtenu %d", c.CapaciteInventaire)
	}
	if c.Argent != 100 {
		t.Errorf("argent de départ attendu 100, obtenu %d", c.Argent)
	}
}

// T8 : un personnage vivant n'est pas touché. Contre-témoin des deux tests
// suivants : sans lui, un isDead qui ressusciterait tout le monde passerait.
func TestIsDeadVivant(t *testing.T) {
	c := nouveauPersonnage(80, 1)

	if c.isDead() {
		t.Error("à 1 PV, le personnage est vivant")
	}
	if c.PVActuel != 1 {
		t.Errorf("les PV d'un vivant ne doivent pas changer, obtenu %d", c.PVActuel)
	}
}

// T8 : à exactement 0 PV, le personnage meurt et ressuscite à 50 %.
func TestIsDeadAZero(t *testing.T) {
	c := nouveauPersonnage(80, 0)

	if !c.isDead() {
		t.Error("à 0 PV, le personnage est mort")
	}
	if c.PVActuel != 40 {
		t.Errorf("résurrection à 50 %% de 80 attendue (40), obtenu %d", c.PVActuel)
	}
}

// T8 : sous 0 PV aussi (les dégâts peuvent dépasser les PV restants).
// Un test « PV == 0 » laisserait passer ce cas.
func TestIsDeadNegatif(t *testing.T) {
	c := initCharacter("Test", "Nain", 1, 120, -7)

	if !c.isDead() {
		t.Error("à -7 PV, le personnage est mort")
	}
	if c.PVActuel != 60 {
		t.Errorf("résurrection à 50 %% de 120 attendue (60), obtenu %d", c.PVActuel)
	}
}
