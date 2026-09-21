package monde

import (
	"fmt"
	"path/filepath"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

func ChargerMap3D(cheminOBJ, cheminMTL string) (*core.Node, error) {
	decodeur, err := obj.Decode(cheminOBJ, cheminMTL)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire la map : %w", err)
	}

	// La texture du MTL utilise le mode CLAMP de G3N. On la désactive ici
	// afin de la remplacer par la même image configurée en mode REPEAT.
	descriptionHerbe, existe := decodeur.Materials["grass"]
	if !existe {
		return nil, fmt.Errorf("le matériau grass est absent du fichier MTL")
	}
	descriptionHerbe.MapKd = ""

	map3d, err := decodeur.NewGroup()
	if err != nil {
		return nil, fmt.Errorf("impossible de construire la map : %w", err)
	}

	cheminHerbe := filepath.Join(filepath.Dir(cheminMTL), "textures", "grass_pixel.png")
	matiereHerbe, err := creerMateriauHerbe(cheminHerbe)
	if err != nil {
		return nil, err
	}

	if err := appliquerMateriauAuPremierObjet(map3d, matiereHerbe); err != nil {
		return nil, err
	}
	return map3d, nil
}

func chargerTextureRepetee(chemin string) (*texture.Texture2D, error) {
	imageTexture, err := texture.NewTexture2DFromImage(chemin)
	if err != nil {
		return nil, fmt.Errorf("impossible de charger la texture : %w", err)
	}

	imageTexture.SetWrapS(gls.REPEAT)
	imageTexture.SetWrapT(gls.REPEAT)
	imageTexture.SetMagFilter(gls.NEAREST)
	imageTexture.SetMinFilter(gls.NEAREST)

	return imageTexture, nil
}

func creerMateriauHerbe(chemin string) (*material.Standard, error) {
	imageTexture, err := chargerTextureRepetee(chemin)
	if err != nil {
		return nil, err
	}

	couleurBlanche := &math32.Color{
		R: 1,
		G: 1,
		B: 1,
	}

	matiere := material.NewStandard(couleurBlanche)
	matiere.AddTexture(imageTexture)

	return matiere, nil
}

func appliquerMateriauAuPremierObjet(groupe *core.Node, matiere *material.Standard) error {
	enfants := groupe.Children()

	if len(enfants) == 0 {
		return fmt.Errorf("le groupe ne contient aucun objet")
	}

	premierMaillage, ok := enfants[0].(*graphic.Mesh)
	if !ok {
		return fmt.Errorf("le premier objet n'est pas un maillage")
	}

	premierMaillage.SetMaterial(matiere)
	return nil
}
