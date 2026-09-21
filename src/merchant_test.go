package main

import "testing"

// T7 : un achat ajoute l'objet à l'inventaire.
func TestBuyAjouteLObjet(t *testing.T) {
	c := nouveauPersonnage(80, 40)

	c.buy(itemPotionDePoison)

	if c.Inventaire[itemPotionDePoison] != 1 {
		t.Errorf("attendu 1 potion de poison, obtenu %d", c.Inventaire[itemPotionDePoison])
	}
}

// T12 : un achat avec l'inventaire plein n'ajoute rien (message « Inventaire
// plein » au lieu de « Vous avez acheté »).
func TestBuyInventairePlein(t *testing.T) {
	c := nouveauPersonnage(80, 40)
	c.Inventaire[itemPotionDeVie] = 10

	c.buy(itemPotionDePoison)

	if c.Inventaire[itemPotionDePoison] != 0 {
		t.Error("l'objet ne doit pas être ajouté quand l'inventaire est plein")
	}
}
