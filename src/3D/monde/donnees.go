package monde

import (
	"encoding/json"
	"fmt"
	"os"
)

type Zone struct {
	Identifiant string  `json:"id"`
	Nom         string  `json:"name"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
	Rayon       float32 `json:"radius"`
	Description string  `json:"description"`
}

type DonneesMonde struct {
	Nom       string      `json:"name"`
	Taille    float32     `json:"world_size"`
	Spawn     [3]float32  `json:"spawn"`
	Zones     []Zone      `json:"landmarks"`
	Obstacles []Collision `json:"collisions"`
	Ponts     []Collision `json:"bridges"`
	// north_axis
	// up_axis
	// image_width
	// image_to_world
}

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
