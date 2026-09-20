package library

import (
	"testing"
)

// T7 : un achat ajoute l'objet à l'inventaire.
func TestBuyAjouteLObjet(t *testing.T) {
	c := NouveauPersonnage(80, 40)

	c.buy(ItemPotionDePoison)

	if c.Inventaire[ItemPotionDePoison] != 1 {
		t.Errorf("attendu 1 potion de poison, obtenu %d", c.Inventaire[ItemPotionDePoison])
	}
}

// T12 : un achat avec l'inventaire plein n'ajoute rien (message « Inventaire
// plein » au lieu de « Vous avez acheté »).
func TestBuyInventairePlein(t *testing.T) {
	c := NouveauPersonnage(80, 40)
	c.Inventaire[ItemPotionDeVie] = 10

	c.buy(ItemPotionDePoison)

	if c.Inventaire[ItemPotionDePoison] != 0 {
		t.Error("l'objet ne doit pas être ajouté quand l'inventaire est plein")
	}
}
