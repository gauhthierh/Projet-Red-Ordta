package tests

import (
	"testing"

	"ordta/library"
)

func TestInventaire3DUneCaseParExemplaire(t *testing.T) {
	p := library.Character{
		Inventaire:         map[string]int{library.ItemPotionDeVie: 4, library.ItemPotionDeMana: 2},
		CapaciteInventaire: 10,
	}
	verifier := func(vie, mana int) {
		t.Helper()
		objets := p.Inventaire3D()
		if len(objets) != vie+mana || len(objets) != p.TotalInventaire() {
			t.Fatalf("%d cases affichées, attendu %d", len(objets), vie+mana)
		}
		comptes := map[string]int{}
		for _, objet := range objets {
			if objet.Quantite != 1 {
				t.Fatalf("quantité par case : %d", objet.Quantite)
			}
			comptes[objet.Nom]++
		}
		if comptes[library.ItemPotionDeVie] != vie || comptes[library.ItemPotionDeMana] != mana {
			t.Fatalf("répartition incorrecte : %v", comptes)
		}
	}
	verifier(4, 2)
	p.RemoveInventory(library.ItemPotionDeVie)
	verifier(3, 2)
	p.Inventaire[library.ItemPotionDeVie] = 8
	verifier(8, 2)
	if p.VerifPlaceInventaire() {
		t.Fatal("10 exemplaires doivent remplir les 10 places")
	}
}
