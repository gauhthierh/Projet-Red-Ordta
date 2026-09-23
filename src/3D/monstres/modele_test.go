package monstres

import "testing"

// Charge les vrais fichiers et textures sans ouvrir de fenêtre.
func TestModelesLivres(t *testing.T) {
	for _, nom := range []string{"gobelin", "gobelin_cuirasse", "chaman", "loup", "troll"} {
		t.Run(nom, func(t *testing.T) {
			a, err := Charger("../../../assets/models/monstres/" + nom + ".json")
			if err != nil {
				t.Fatal(err)
			}
			b, err := Charger("../../../assets/models/monstres/" + nom + ".json")
			if err != nil {
				t.Fatal(err)
			}
			if a.Membres["tete"] == nil || a.Membres["corps"] == nil {
				t.Fatal("pivots manquants")
			}
			a.Membres["tete"].SetRotationZ(.5)
			if b.Membres["tete"].Rotation().Z != 0 {
				t.Fatal("les instances partagent leurs articulations")
			}
		})
	}
}
