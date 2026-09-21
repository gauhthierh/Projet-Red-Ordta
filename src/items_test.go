package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// capturerSortie exécute f et renvoie le texte qu'elle a affiché. Certaines
// règles portent sur l'affichage et pas seulement sur les valeurs : on
// redirige donc temporairement la sortie standard vers un tube.
func capturerSortie(t *testing.T, f func()) string {
	t.Helper()
	ancienne := os.Stdout
	lecture, ecriture, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = ecriture
	f()
	ecriture.Close()
	os.Stdout = ancienne
	sortie, _ := io.ReadAll(lecture)
	return string(sortie)
}

// T5 : une potion rend 50 PV et est retirée de l'inventaire.
func TestTakePotSoin(t *testing.T) {
	c := nouveauPersonnage(100, 40)

	c.takePot()

	if c.PVActuel != 90 {
		t.Errorf("40 + 50 : attendu 90, obtenu %d", c.PVActuel)
	}
	if c.Inventaire[itemPotionDeVie] != 2 {
		t.Errorf("une potion doit être consommée : attendu 2, obtenu %d", c.Inventaire[itemPotionDeVie])
	}
}

// T5 : les PV ne dépassent jamais les PV max.
func TestTakePotPlafond(t *testing.T) {
	c := nouveauPersonnage(100, 70)

	c.takePot()

	if c.PVActuel != 100 {
		t.Errorf("70 + 50 plafonné : attendu 100, obtenu %d", c.PVActuel)
	}
}

// Décision de l'équipe : à pleins PV, la potion est refusée et n'est pas
// consommée, pour ne pas la gâcher.
func TestTakePotPleineSante(t *testing.T) {
	c := nouveauPersonnage(80, 80)

	c.takePot()

	if c.Inventaire[itemPotionDeVie] != 3 {
		t.Errorf("la potion ne doit pas être consommée : attendu 3, obtenu %d", c.Inventaire[itemPotionDeVie])
	}
}

// T5 : sans potion, rien ne change (message « Aucune potion »).
func TestTakePotSansPotion(t *testing.T) {
	c := nouveauPersonnage(80, 40)
	c.Inventaire = map[string]int{}

	c.takePot()

	if c.PVActuel != 40 {
		t.Errorf("sans potion, les PV ne doivent pas changer : obtenu %d", c.PVActuel)
	}
}

// T10 : spellBook apprend Boule de Feu une seule fois.
func TestSpellBookUneSeuleFois(t *testing.T) {
	c := nouveauPersonnage(80, 40)

	if !c.spellBook() {
		t.Fatal("le premier apprentissage doit réussir")
	}
	if c.spellBook() {
		t.Error("le second apprentissage doit être refusé")
	}
	if len(c.Skill) != 2 {
		t.Errorf("attendu 2 sorts (Coup de poing, Boule de Feu), obtenu %v", c.Skill)
	}
}

// T10 : le livre est consommé quand le sort est appris, et reste dans
// l'inventaire si le sort est déjà connu (décision de l'équipe).
func TestUseItemLivre(t *testing.T) {
	c := nouveauPersonnage(80, 40)
	c.Inventaire[itemLivreBouleDeFeu] = 2

	c.useItem(itemLivreBouleDeFeu) // apprend le sort, consomme un livre
	c.useItem(itemLivreBouleDeFeu) // déjà connu : le livre reste

	if c.Inventaire[itemLivreBouleDeFeu] != 1 {
		t.Errorf("attendu 1 livre restant, obtenu %d", c.Inventaire[itemLivreBouleDeFeu])
	}
}

// Un objet sans effet l'annonce et n'est pas consommé.
func TestUseItemObjetSansEffet(t *testing.T) {
	c := nouveauPersonnage(80, 40)
	c.Inventaire["Caillou"] = 1

	c.useItem("Caillou")

	if c.Inventaire["Caillou"] != 1 {
		t.Error("un objet sans effet ne doit pas être consommé")
	}
}

// sautSiCourt saute les tests du poison avec go test -short : poisonPot
// attend réellement 1 seconde par tick (time.Sleep), comme le demande T9.
func sautSiCourt(t *testing.T) {
	if testing.Short() {
		t.Skip("test long (attentes réelles de 1 s) : sauté avec -short")
	}
}

// T9 : 3 ticks de 10 dégâts, et la potion de poison est consommée.
func TestPoisonPot(t *testing.T) {
	sautSiCourt(t)
	c := nouveauPersonnage(80, 80)
	c.Inventaire[itemPotionDePoison] = 1

	c.poisonPot()

	if c.PVActuel != 50 {
		t.Errorf("80 - 3 x 10 : attendu 50, obtenu %d", c.PVActuel)
	}
	if c.Inventaire[itemPotionDePoison] != 0 {
		t.Error("la potion de poison doit être consommée")
	}
}

// T8 + T9 : le poison s'arrête à la mort. À 20 PV : 10, puis 0, mort,
// résurrection à 40. Si le poison continuait, le 3e tick donnerait 30.
func TestPoisonPotArretALaMort(t *testing.T) {
	sautSiCourt(t)
	c := nouveauPersonnage(80, 20)
	c.Inventaire[itemPotionDePoison] = 1

	c.poisonPot()

	if c.PVActuel != 40 {
		t.Errorf("attendu 40 (résurrection, poison arrêté), obtenu %d", c.PVActuel)
	}
}

// T9 : les PV ne s'affichent jamais en négatif. À 5 PV, le premier tick
// ferait -5 : ils sont ramenés à 0, puis le personnage meurt et ressuscite.
// On vérifie l'affichage : la valeur finale seule ne suffit pas, car
// isDead (PV <= 0) ressusciterait aussi le personnage à -5, et le test
// passerait sans la remise à 0.
func TestPoisonPotSousZero(t *testing.T) {
	sautSiCourt(t)
	c := nouveauPersonnage(100, 5)
	c.Inventaire[itemPotionDePoison] = 1

	sortie := capturerSortie(t, c.poisonPot)

	if !strings.Contains(sortie, "PV : 0 / 100") || strings.Contains(sortie, "-5") {
		t.Errorf("attendu « PV : 0 / 100 » et aucun PV négatif, affiché :\n%s", sortie)
	}
	if c.PVActuel != 50 {
		t.Errorf("attendu 50 (résurrection à 50 %% de 100), obtenu %d", c.PVActuel)
	}
}
