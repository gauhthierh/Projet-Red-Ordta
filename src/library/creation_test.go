package library

import "testing"

// T11 : formatName accepte uniquement des lettres (accentuées comprises) et
// remet le nom en forme : majuscule initiale, reste en minuscules.
// Écrit sous forme de tableau : chaque ligne est un cas, et en ajouter un
// ne demande qu'une ligne.
func TestFormatName(t *testing.T) {
	cas := []struct {
		saisie    string
		attendu   string
		estValide bool
	}{
		{"jean", "Jean", true},
		{"jEAN", "Jean", true},
		{"a", "A", true},           // une seule lettre : l'ancien code l'ignorait
		{"éLOÏSE", "Éloïse", true}, // accents : il faut traiter des runes, pas des octets
		{"jean pierre", "", false}, // espace refusé
		{"Jean-Pierre", "", false}, // tiret refusé
		{"jean3", "", false},       // chiffre refusé
		{"", "", false},            // vide refusé
	}

	for _, c := range cas {
		obtenu, valide := FormatName(c.saisie)
		if valide != c.estValide || obtenu != c.attendu {
			t.Errorf("formatName(%q) = (%q, %v), attendu (%q, %v)",
				c.saisie, obtenu, valide, c.attendu, c.estValide)
		}
	}
}
