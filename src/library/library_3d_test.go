package library

import "testing"

func TestInventaire3DUneCaseParExemplaire(t *testing.T) {
	p := Character{Inventaire: map[string]int{ItemPotionDeVie: 4, ItemPotionDeMana: 2}, CapaciteInventaire: 10}
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
		if comptes[ItemPotionDeVie] != vie || comptes[ItemPotionDeMana] != mana {
			t.Fatalf("répartition incorrecte : %v", comptes)
		}
	}
	verifier(4, 2)
	p.RemoveInventory(ItemPotionDeVie)
	verifier(3, 2)
	p.Inventaire[ItemPotionDeVie] = 8
	verifier(8, 2)
	if p.VerifPlaceInventaire() {
		t.Fatal("10 exemplaires doivent remplir les 10 places")
	}
}
