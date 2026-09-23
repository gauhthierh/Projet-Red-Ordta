package monde

import (
	"github.com/g3n/engine/math32"
	"ordta/library"
	"os"
	"path/filepath"
	"testing"
)

// TestSauvegarde3DRetrouveProgression : Vérifie l'écriture, la relecture et le remplacement d'une sauvegarde.
func TestSauvegarde3DRetrouveProgression(t *testing.T) {
	chemin := filepath.Join(t.TempDir(), "partie.json")
	joueur := library.NouveauPersonnage3D("Test", "Humain")
	joueur.Argent = 57
	joueur.Equipement.Tete = library.Armurerie[0]
	partie := Partie3D{Personnage: joueur, Position: math32.Vector3{X: 20, Y: 30, Z: .265}, AngleHorizontal: .4}
	if err := SauvegarderPartie3D(chemin, partie); err != nil {
		t.Fatal(err)
	}
	partie.Personnage.Argent = 42
	if err := SauvegarderPartie3D(chemin, partie); err != nil {
		t.Fatal("remplacement :", err)
	}
	relue, err := ChargerPartie3D(chemin)
	if err != nil {
		t.Fatal(err)
	}
	if relue.Personnage.Argent != 42 || relue.Personnage.Equipement.Tete.Nom != joueur.Equipement.Tete.Nom || relue.Position != partie.Position || relue.Personnage.Inventaire[library.ItemPotionDeVie] != 3 {
		t.Fatal("progression non conservée")
	}
}

// TestSauvegarde3DAbsenteEtCorrompue : Vérifie qu'une absence de sauvegarde ne se confond pas avec un fichier abîmé.
func TestSauvegarde3DAbsenteEtCorrompue(t *testing.T) {
	chemin := filepath.Join(t.TempDir(), "partie.json")
	partie, err := ChargerPartie3D(chemin)
	if err != nil || partie != nil {
		t.Fatal("une nouvelle partie doit rester possible")
	}
	if err := os.WriteFile(chemin, []byte("invalide"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ChargerPartie3D(chemin); err == nil {
		t.Fatal("corruption non signalée")
	}
}

// TestAncienneSauvegardeMigration : Vérifie que les anciens noms sont convertis au chargement.
func TestAncienneSauvegardeMigration(t *testing.T) {
	chemin := filepath.Join(t.TempDir(), "ancienne.json")
	p := library.NouveauPersonnage3D("Ancien", "Nain")
	p.Inventaire = map[string]int{"Potion de Mana": 2, "Peau de Troll": 1, "Livre de sort : Boule de feu": 1}
	p.Skill = []string{"Coup de poing", "Boule de feu"}
	p.AttaquesPhysiques = []string{"Attaque Basique"}
	if err := SauvegarderPartie3D(chemin, Partie3D{Personnage: p}); err != nil {
		t.Fatal(err)
	}
	relue, err := ChargerPartie3D(chemin)
	if err != nil {
		t.Fatal(err)
	}
	if relue.Personnage.Nom != "Ancien" || relue.Personnage.Classe != "Nain" || relue.Personnage.Inventaire[library.ItemPotionDeMana] != 2 || relue.Personnage.Skill[1] != library.SortGrosseBouleDeFeu || relue.Personnage.AttaquesPhysiques[0] != library.AttaqueBasique {
		t.Fatalf("migration incomplète : %+v", relue.Personnage)
	}
}
