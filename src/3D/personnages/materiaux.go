package personnages

import (
	"image"
	imagecolor "image/color"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Palette et textures procédurales reprises du personnage de l'ancien jeu.
// Elles sont produites en Go ; aucun fichier de texture n'est nécessaire.
var (
	bois   = couleur(.28, .15, .08)
	acier  = couleur(.58, .66, .72)
	or     = couleur(.92, .62, .12)
	tissu  = couleur(.12, .29, .40)
	peau   = couleur(.83, .58, .38)
	sombre = couleur(.045, .035, .025)

	textures = make(map[*math32.Color]*texture.Texture2D)
)

func couleur(rouge, vert, bleu float32) *math32.Color {
	return &math32.Color{R: rouge, G: vert, B: bleu}
}

func matiereSurface(teinte *math32.Color) *material.Standard {
	matiere := material.NewStandard(teinte)
	if typeSurface := typeTexture(teinte); typeSurface != "" {
		imageTexture := textures[teinte]
		if imageTexture == nil {
			imageTexture = texture.NewTexture2DFromRGBA(imageSurface(typeSurface))
			imageTexture.SetMagFilter(gls.LINEAR)
			textures[teinte] = imageTexture
		}
		matiere.AddTexture(imageTexture.Incref())
	}
	if teinte == acier || teinte == or {
		matiere.SetShininess(70)
	} else {
		matiere.SetShininess(6)
		matiere.SetSpecularColor(couleur(.10, .10, .10))
	}
	return matiere
}

func typeTexture(teinte *math32.Color) string {
	switch teinte {
	case bois:
		return "bois"
	case acier, or:
		return "metal"
	case tissu:
		return "tissu"
	case peau:
		return "peau"
	default:
		return ""
	}
}

func imageSurface(typeSurface string) *image.RGBA {
	const cote = 64
	img := image.NewRGBA(image.Rect(0, 0, cote, cote))
	for y := 0; y < cote; y++ {
		for x := 0; x < cote; x++ {
			bruit := int(bruit2(x, y)%27) - 13
			valeur := 220 + bruit
			switch typeSurface {
			case "bois":
				valeur = 217 + bruit/2 + ((x/5+y/23)%3)*9
				if x%17 == 0 || (x+y/3)%29 == 0 {
					valeur -= 36
				}
			case "tissu":
				valeur = 222 + bruit/3
				if x%4 == 0 || y%4 == 0 {
					valeur -= 15
				}
			case "metal":
				valeur = 220 + bruit/4
				if y%8 == 0 {
					valeur += 12
				}
			case "peau":
				valeur = 231 + bruit/4
			}
			if valeur < 90 {
				valeur = 90
			}
			if valeur > 255 {
				valeur = 255
			}
			composante := uint8(valeur)
			img.SetRGBA(x, y, imagecolor.RGBA{R: composante, G: composante, B: composante, A: 255})
		}
	}
	return img
}

func bruit2(x, y int) uint32 {
	valeur := uint32(x)*374761393 + uint32(y)*668265263 + 2026
	valeur = (valeur ^ (valeur >> 13)) * 1274126177
	return valeur ^ (valeur >> 16)
}
