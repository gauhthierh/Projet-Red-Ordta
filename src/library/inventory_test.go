package library

import (
	"slices"
	"testing"
)

// T7 : removeInventory retire un exemplaire et supprime la clé à 0, pour
// qu'un objet épuisé disparaisse des listes au lieu d'y rester avec « 0 ».
func TestRemoveInventorySupprimeLaCle(t *testing.T) {
	c := NouveauPersonnage(80, 40)
	c.Inventaire = map[string]int{ItemPotionDePoison: 1}

	if !c.RemoveInventory(ItemPotionDePoison) {
		t.Fatal("l'objet était présent, le retrait doit réussir")
	}
	if _, present := c.Inventaire[ItemPotionDePoison]; present {
		t.Error("la clé doit être supprimée quand la quantité tombe à 0")
	}
}

// T7 : retirer un objet absent échoue sans créer de clé à 0 ni de
// quantité négative.
func TestRemoveInventoryObjetAbsent(t *testing.T) {
	c := NouveauPersonnage(80, 40)

	if c.RemoveInventory("Objet absent") {
		t.Error("retirer un objet absent doit renvoyer false")
	}
	if _, present := c.Inventaire["Objet absent"]; present {
		t.Error("retirer un objet absent ne doit pas créer de clé")
	}
}

// Menu d'inventaire : les objets sont toujours dans le même ordre, sinon un
// même numéro désignerait des objets différents. Le parcours d'une map est
// aléatoire en Go : on répète donc l'appel 50 fois.
func TestSortedItemsOrdreStable(t *testing.T) {
	c := NouveauPersonnage(80, 40)
	c.Inventaire = map[string]int{
		ItemPotionDeVie:     1,
		ItemLivreBouleDeFeu: 1,
		ItemPotionDePoison:  1,
		"Zèbre":             1,
		"Arc":               1,
	}
	attendu := []string{"Arc", ItemLivreBouleDeFeu, ItemPotionDePoison, ItemPotionDeVie, "Zèbre"}

	for i := 0; i < 50; i++ {
		if obtenu := c.SortedItems(); !slices.Equal(obtenu, attendu) {
			t.Fatalf("appel %d : %v, attendu %v", i, obtenu, attendu)
		}
	}
}

// T12 : le nombre d'objets est la somme des quantités, pas le nombre de
// sortes. Ici 3 + 2 = 5 objets pour 2 sortes : un calcul par len() donnerait 2.
func TestTotalInventaireSommeDesQuantites(t *testing.T) {
	c := NouveauPersonnage(80, 40)
	c.Inventaire = map[string]int{ItemPotionDeVie: 3, ItemPotionDePoison: 2}

	if total := c.TotalInventaire(); total != 5 {
		t.Errorf("attendu 5 objets, obtenu %d", total)
	}
}

// T12 : au-delà de 10 objets, addInventory refuse et n'ajoute rien.
func TestAddInventoryLimite(t *testing.T) {
	c := NouveauPersonnage(80, 40) // 3 potions au départ

	for i := 0; i < 7; i++ {
		if !c.AddInventory(ItemPotionDePoison) {
			t.Fatalf("ajout %d refusé alors que l'inventaire n'est pas plein", i+1)
		}
	}
	if c.AddInventory(ItemPotionDePoison) {
		t.Error("le 11e objet doit être refusé")
	}
	if total := c.TotalInventaire(); total != 10 {
		t.Errorf("un ajout refusé ne doit rien ajouter : attendu 10, obtenu %d", total)
	}
}

// T12 : une place libérée permet de nouveau d'ajouter un objet.
func TestAddInventoryApresPlaceLiberee(t *testing.T) {
	c := NouveauPersonnage(80, 40)
	// Remplit l'inventaire jusqu'à ce que addInventory refuse.
	for c.AddInventory(ItemPotionDePoison) {
	}

	c.RemoveInventory(ItemPotionDePoison)

	if !c.AddInventory(ItemPotionDePoison) {
		t.Error("après avoir retiré un objet, l'ajout doit être accepté")
	}
}

// T12 : la limite est lue dans le champ CapaciteInventaire et n'est pas
// écrite en dur. Si on augmente le champ (ce que fera la tâche 18),
// l'inventaire accepte davantage d'objets.
func TestCapaciteLueDansLeChamp(t *testing.T) {
	c := NouveauPersonnage(80, 40)
	c.CapaciteInventaire = 12

	// Remplit l'inventaire jusqu'à ce que addInventory refuse.
	for c.AddInventory(ItemPotionDePoison) {
	}

	if total := c.TotalInventaire(); total != 12 {
		t.Errorf("avec une capacité de 12, attendu 12 objets, obtenu %d", total)
	}
}
