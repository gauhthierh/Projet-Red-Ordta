package monde

import (
	"os"
	"testing"
)

// Tester le format ne nécessite ni sortie son ni contexte OpenAL.
func TestFormatAudioLivreEtFichiersInvalides(t *testing.T) {
	data, err := os.ReadFile("../../../assets/audio/menu.wav")
	if err != nil {
		t.Fatal(err)
	}
	if !formatAudioValide(data) {
		t.Fatal("musique livrée refusée")
	}
	for _, invalide := range [][]byte{nil, data[:20], data[:len(data)-2], []byte("pas un WAV")} {
		if formatAudioValide(invalide) {
			t.Fatal("fichier invalide accepté")
		}
	}
}
