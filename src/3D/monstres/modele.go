// Package monstres charge les modèles fixes, sans gérer les règles de combat.
package monstres

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Les noms des membres sont les clés du fichier JSON : tete, bras_droit, etc.
// Chaque instance possède ses propres pivots, donc peut être animée séparément.
type Monstre struct {
	noeud   *core.Node
	Membres map[string]*core.Node
}

type partie struct {
	Nom    string     `json:"nom"`
	Pivot  [3]float32 `json:"pivot"`
	Pieces []piece    `json:"pieces"`
}
type piece struct {
	Forme    string     `json:"forme"`
	Position [3]float32 `json:"position"`
	Taille   [3]float32 `json:"taille"`
	Couleur  [3]float32 `json:"couleur"`
	Texture  string     `json:"texture,omitempty"`
}
type modele struct {
	Nom     string   `json:"nom"`
	Parties []partie `json:"parties"`
}

// Charger utilise un chemin explicite, comme le chargement de la carte.
// Z est vertical ; les monstres regardent vers -Y, comme le personnage.
func Charger(chemin string) (*Monstre, error) {
	contenu, err := os.ReadFile(chemin)
	if err != nil {
		return nil, fmt.Errorf("modèle monstre : %w", err)
	}
	var description modele
	if err := json.Unmarshal(contenu, &description); err != nil {
		return nil, err
	}
	m := &Monstre{noeud: core.NewNode(), Membres: make(map[string]*core.Node)}
	for _, partie := range description.Parties {
		if _, existe := m.Membres[partie.Nom]; existe {
			return nil, fmt.Errorf("membre répété : %s", partie.Nom)
		}
		pivot := core.NewNode()
		pivot.SetPosition(partie.Pivot[0], partie.Pivot[1], partie.Pivot[2])
		m.Membres[partie.Nom] = pivot
		m.noeud.Add(pivot)
		for _, p := range partie.Pieces {
			if p.Taille[0] <= 0 || p.Taille[1] <= 0 || p.Taille[2] <= 0 {
				return nil, fmt.Errorf("taille invalide : %s", partie.Nom)
			}
			var forme *geometry.Geometry
			switch p.Forme {
			case "boite":
				forme = geometry.NewBox(p.Taille[0], p.Taille[1], p.Taille[2])
			case "ellipsoide":
				forme = geometry.NewSphere(.5, 10, 8)
			default:
				return nil, fmt.Errorf("forme inconnue : %s", p.Forme)
			}
			matiere := material.NewStandard(&math32.Color{R: p.Couleur[0], G: p.Couleur[1], B: p.Couleur[2]})
			matiere.SetShininess(8)
			if p.Texture != "" {
				tex, err := texture.NewTexture2DFromImage(filepath.Join(filepath.Dir(chemin), "textures", p.Texture+".png"))
				if err != nil {
					return nil, fmt.Errorf("texture %s : %w", p.Texture, err)
				}
				matiere.AddTexture(tex)
			}
			mesh := graphic.NewMesh(forme, matiere)
			if p.Forme == "ellipsoide" {
				mesh.SetScale(p.Taille[0], p.Taille[1], p.Taille[2])
			}
			mesh.SetPosition(p.Position[0], p.Position[1], p.Position[2])
			pivot.Add(mesh)
		}
	}
	return m, nil
}

func (m *Monstre) Noeud() *core.Node { return m.noeud }
