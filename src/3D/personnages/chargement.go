package personnages

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

// Ce chargeur utilise des maillages et textures fixes, embarqués dans le programme.
// Les tenues construites dans tenues.go suivent un autre chemin de création.
//
//go:embed assets/*
var fichiers embed.FS

// piece : Contient les données d'un maillage : positions, normales, UV et indices de triangles.
type piece struct {
	Sommets  []float32  `json:"vertices"`
	Normales []float32  `json:"normals"`
	UV       []float32  `json:"uv"`
	Indices  []uint32   `json:"indices"`
	Couleur  [3]float32 `json:"color"`
	Texture  string     `json:"texture"`
}

// modeleFichier : Range les maillages de l'ancien personnage par membre articulé.
type modeleFichier struct {
	Corps       []piece `json:"body"`
	JambeGauche []piece `json:"left_leg"`
	JambeDroite []piece `json:"right_leg"`
	BrasGauche  []piece `json:"left_arm"`
	BrasDroit   []piece `json:"right_arm"`
}

// lireModele : Décode le modèle embarqué ; une erreur signifie que l'asset livré avec le programme est invalide.
func lireModele() modeleFichier {
	contenu, err := fichiers.ReadFile("assets/personnage.json")
	if err != nil {
		panic(fmt.Errorf("lire le modèle du personnage : %w", err))
	}
	var modele modeleFichier
	if err := json.Unmarshal(contenu, &modele); err != nil {
		panic(fmt.Errorf("décoder le modèle du personnage : %w", err))
	}
	return modele
}

// ajouterPieces attache les maillages fixes à un nœud. Les coordonnées du
// fichier sont directement en Z vertical, comme celles de la scène.
func ajouterPieces(parent *core.Node, pieces []piece) {
	for _, partie := range pieces {
		forme := geometry.NewGeometry()
		forme.SetIndices(math32.ArrayU32(partie.Indices))
		forme.AddVBO(gls.NewVBO(math32.ArrayF32(partie.Sommets)).AddAttrib(gls.VertexPosition))
		forme.AddVBO(gls.NewVBO(math32.ArrayF32(partie.Normales)).AddAttrib(gls.VertexNormal))
		forme.AddVBO(gls.NewVBO(math32.ArrayF32(partie.UV)).AddAttrib(gls.VertexTexcoord))
		parent.Add(graphic.NewMesh(forme, matiere(partie.Couleur, partie.Texture)))
	}
}
