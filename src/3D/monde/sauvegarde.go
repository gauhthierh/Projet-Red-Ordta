package monde

import (
	"encoding/json"
	"fmt"
	"github.com/g3n/engine/math32"
	"ordta/library"
	"os"
	"path/filepath"
)

const cheminSauvegarde = "sauvegardes/partie_3d.json"

type Partie3D struct {
	Version         int
	Personnage      library.Character
	Position        math32.Vector3
	AngleHorizontal float32
	AngleVertical   float32
}

func ChargerPartie3D(chemin string) (*Partie3D, error) {
	contenu, err := os.ReadFile(chemin)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var partie Partie3D
	if err = json.Unmarshal(contenu, &partie); err != nil {
		return nil, fmt.Errorf("sauvegarde illisible : %w", err)
	}
	if partie.Version != 1 || partie.Personnage.Inventaire == nil || partie.Personnage.ExperienceMax <= 0 || partie.Personnage.PVMaxBase <= 0 {
		return nil, fmt.Errorf("sauvegarde 3D invalide ou incompatible")
	}
	partie.Personnage.MettreAJourPvMax()
	return &partie, nil
}

// Écriture temporaire puis remplacement : une écriture interrompue ne tronque pas la sauvegarde.
func SauvegarderPartie3D(chemin string, partie Partie3D) error {
	partie.Version = 1
	contenu, err := json.MarshalIndent(partie, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(chemin), 0755); err != nil {
		return err
	}
	temporaire, err := os.CreateTemp(filepath.Dir(chemin), "partie-*.tmp")
	if err != nil {
		return err
	}
	nom := temporaire.Name()
	defer os.Remove(nom)
	if _, err = temporaire.Write(contenu); err != nil {
		temporaire.Close()
		return err
	}
	if err = temporaire.Sync(); err != nil {
		temporaire.Close()
		return err
	}
	if err = temporaire.Close(); err != nil {
		return err
	}
	return os.Rename(nom, chemin)
}
