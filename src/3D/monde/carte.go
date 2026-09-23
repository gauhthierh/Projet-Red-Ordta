package monde

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Plusieurs matériaux peuvent partager la même image avec une teinte différente.
var texturesCarte = map[string]string{
	"leaf":        "foliage_detailed.png",
	"leaf_light":  "foliage_detailed.png",
	"grass":       "grass_detailed.png",
	"dark_grass":  "grass_detailed.png",
	"field":       "grass_detailed.png",
	"marsh":       "marsh_detailed.png",
	"crop":        "crop_detailed.png",
	"path":        "path_detailed.png",
	"sand":        "path_detailed.png",
	"hay":         "path_detailed.png",
	"stone":       "stone_detailed.png",
	"light_stone": "stone_detailed.png",
	"rock":        "stone_detailed.png",
	"wall":        "stone_detailed.png",
	"snow":        "stone_detailed.png",
	"wood":        "wood_detailed.png",
	"water":       "water_detailed.png",
	"magic":       "water_detailed.png",
	"roof_red":    "roof_detailed.png",
	"roof_blue":   "roof_detailed.png",
	"metal":       "metal_detailed.png",
	"cloth_red":   "cloth_detailed.png",
	"cloth_blue":  "cloth_detailed.png",
}

var teintesCarte = map[string]math32.Color{
	"leaf":        {R: .75, G: .85, B: .7},
	"leaf_light":  {R: 1, G: 1, B: .85},
	"grass":       {R: .62, G: .72, B: .54},
	"dark_grass":  {R: .48, G: .68, B: .45},
	"field":       {R: .72, G: .62, B: .34},
	"marsh":       {R: 1, G: 1, B: 1},
	"crop":        {R: 1, G: 1, B: 1},
	"path":        {R: 1, G: 1, B: 1},
	"sand":        {R: 1, G: .86, B: .58},
	"hay":         {R: 1, G: .78, B: .32},
	"stone":       {R: .82, G: .86, B: .92},
	"light_stone": {R: 1, G: 1, B: .96},
	"rock":        {R: .58, G: .62, B: .68},
	"wall":        {R: 1, G: .88, B: .72},
	"snow":        {R: 1, G: 1, B: 1},
	"wood":        {R: 1, G: 1, B: 1},
	"water":       {R: 1, G: 1, B: 1},
	"magic":       {R: .44, G: 1, B: 1},
	"roof_red":    {R: .78, G: .24, B: .16},
	"roof_blue":   {R: .22, G: .48, B: .84},
	"metal":       {R: .72, G: .78, B: .86},
	"cloth_red":   {R: .78, G: .16, B: .13},
	"cloth_blue":  {R: .12, G: .34, B: .78},
}

// ChargerMap3D : Assemble les géométries OBJ avec leurs matériaux et leurs textures répétées.
func ChargerMap3D(cheminOBJ, cheminMTL string) (*core.Node, error) {
	decodeur, err := obj.Decode(cheminOBJ, cheminMTL)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire la map : %w", err)
	}

	// G3N bride les textures OBJ au bord. Elles sont rechargées plus bas en
	// mode répétition, avec la même méthode pour tous les matériaux.
	for _, description := range decodeur.Materials {
		description.MapKd = ""
	}

	map3d := core.NewNode()
	texturesChargees := make(map[string]*texture.Texture2D)
	dossierTextures := filepath.Join(filepath.Dir(cheminMTL), "textures")

	for index := range decodeur.Objects {
		objet := &decodeur.Objects[index]
		nomMateriau, err := lireMateriauDansNom(objet.Name)
		if err != nil {
			return nil, err
		}

		geometrie, err := decodeur.NewGeometry(objet)
		if err != nil {
			return nil, fmt.Errorf("géométrie %s invalide : %w", objet.Name, err)
		}

		matiere, err := creerMateriauCarte(nomMateriau, decodeur, dossierTextures, texturesChargees)
		if err != nil {
			return nil, err
		}

		map3d.Add(graphic.NewMesh(geometrie, matiere))
	}

	return map3d, nil
}

// lireMateriauDansNom : Retrouve le matériau selon la convention de nommage utilisée par les objets de cette carte.
func lireMateriauDansNom(nomObjet string) (string, error) {
	const debut = "__material_"
	const fin = "__part_"

	positionDebut := strings.LastIndex(nomObjet, debut)
	positionFin := strings.LastIndex(nomObjet, fin)
	if positionDebut == -1 || positionFin == -1 || positionFin <= positionDebut {
		return "", fmt.Errorf("le matériau est absent du nom OBJ %q", nomObjet)
	}

	positionDebut += len(debut)
	return nomObjet[positionDebut:positionFin], nil
}

// creerMateriauCarte : Associe une couleur et une texture au matériau en réutilisant les images déjà chargées.
func creerMateriauCarte(nom string, decodeur *obj.Decoder, dossierTextures string, texturesChargees map[string]*texture.Texture2D) (*material.Standard, error) {
	description, existe := decodeur.Materials[nom]
	if !existe {
		return nil, fmt.Errorf("le matériau %q est absent du fichier MTL", nom)
	}

	couleur := description.Diffuse
	if teinte, existe := teintesCarte[nom]; existe {
		couleur = teinte
	}

	matiere := material.NewStandard(&couleur)
	matiere.SetSpecularColor(&description.Specular)
	matiere.SetShininess(description.Shininess)

	nomTexture, existe := texturesCarte[nom]
	if !existe {
		return matiere, nil
	}

	imageTexture, existe := texturesChargees[nomTexture]
	if !existe {
		var err error
		imageTexture, err = chargerTextureRepetee(filepath.Join(dossierTextures, nomTexture))
		if err != nil {
			return nil, err
		}
		texturesChargees[nomTexture] = imageTexture
	}

	matiere.AddTexture(imageTexture.Incref())
	return matiere, nil
}

// chargerTextureRepetee : Autorise la répétition de l'image et filtre les textures lointaines pour limiter le scintillement.
func chargerTextureRepetee(chemin string) (*texture.Texture2D, error) {
	imageTexture, err := texture.NewTexture2DFromImage(chemin)
	if err != nil {
		return nil, fmt.Errorf("impossible de charger la texture %s : %w", chemin, err)
	}

	imageTexture.SetWrapS(gls.REPEAT)
	imageTexture.SetWrapT(gls.REPEAT)
	imageTexture.SetMagFilter(gls.NEAREST)
	// Les mipmaps évitent le scintillement des motifs fins vus de loin.
	imageTexture.SetMinFilter(gls.LINEAR_MIPMAP_LINEAR)
	return imageTexture, nil
}
