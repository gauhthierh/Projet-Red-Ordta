package monde

import (
	"encoding/json"
	"fmt"
	"os"
)

// Zone : Décrit un lieu du plan ; ce repère ne constitue pas à lui seul un obstacle.
type Zone struct {
	Identifiant string  `json:"id"`
	Nom         string  `json:"name"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
	Rayon       float32 `json:"radius"`
	Description string  `json:"description"`
}

// DonneesMonde : Contient les informations du JSON, dans le même repère XY au sol et Z vertical que les modèles.
type DonneesMonde struct {
	Nom       string      `json:"name"`
	Taille    float32     `json:"world_size"`
	Spawn     [3]float32  `json:"spawn"`
	Zones     []Zone      `json:"landmarks"`
	Obstacles []Collision `json:"collisions"`
	Ponts     []Collision `json:"bridges"`
	// Les métadonnées de projection supplémentaires du JSON ne sont pas lues ici.
	// Le jeu utilise directement le repère XY au sol, avec Z pour la hauteur.
}

// ChargerZones : Lit les repères et collisions du JSON ; le décor visible est chargé depuis l'OBJ.
func ChargerZones(cheminJSON string) (DonneesMonde, error) {
	var plan DonneesMonde

	contenu, err := os.ReadFile(cheminJSON)
	if err != nil {
		return plan, fmt.Errorf("Impossible de lire la zone : %w", err)
	}

	if err := json.Unmarshal(contenu, &plan); err != nil {
		return plan, fmt.Errorf("Zone JSON invalide : %w", err)
	}

	return plan, nil
}
