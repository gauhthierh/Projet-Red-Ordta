package personnages

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

var textures = make(map[string]*texture.Texture2D)

func matiere(rgb [3]float32, nomTexture string) *material.Standard {
	couleur := &math32.Color{R: rgb[0], G: rgb[1], B: rgb[2]}
	resultat := material.NewStandard(couleur)

	if nomTexture == "metal" {
		resultat.SetShininess(70)
	} else {
		resultat.SetShininess(6)
		resultat.SetSpecularColor(&math32.Color{R: .1, G: .1, B: .1})
	}

	if nomTexture != "" {
		resultat.AddTexture(chargerTexture(nomTexture).Incref())
	}
	return resultat
}

func chargerTexture(nom string) *texture.Texture2D {
	if existante := textures[nom]; existante != nil {
		return existante
	}

	contenu, err := fichiers.ReadFile("assets/" + nom + ".png")
	if err != nil {
		panic(fmt.Errorf("lire la texture %s : %w", nom, err))
	}
	imageLue, err := png.Decode(bytes.NewReader(contenu))
	if err != nil {
		panic(fmt.Errorf("décoder la texture %s : %w", nom, err))
	}

	bitmap := image.NewRGBA(imageLue.Bounds())
	draw.Draw(bitmap, bitmap.Bounds(), imageLue, imageLue.Bounds().Min, draw.Src)
	tex := texture.NewTexture2DFromRGBA(bitmap)
	tex.SetMagFilter(gls.LINEAR)
	textures[nom] = tex
	return tex
}
