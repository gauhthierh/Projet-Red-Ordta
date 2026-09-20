package game3d

import (
	"image"
	imagecolor "image/color"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// These shared procedural textures are generated in Go once. They add wood
// grain, stone grain, woven cloth and roof tiles without external asset files.
var surfaces = make(map[*math32.Color]*texture.Texture2D)

func surfaceMaterial(tint *math32.Color) *material.Standard {
	m := material.NewStandard(tint)
	if kind := surfaceKind(tint); kind != "" {
		tex := surfaces[tint]
		if tex == nil {
			tex = texture.NewTexture2DFromRGBA(surfaceImage(kind))
			tex.SetMagFilter(gls.LINEAR)
			surfaces[tint] = tex
		}
		m.AddTexture(tex.Incref())
	}
	if tint == steel || tint == gold {
		m.SetShininess(70)
	} else {
		m.SetShininess(6)
		m.SetSpecularColor(color(.10, .10, .10))
	}
	return m
}

func surfaceKind(c *math32.Color) string {
	switch c {
	case wood:
		return "wood"
	case plaster:
		return "plaster"
	case roof:
		return "roof"
	case stone:
		return "stone"
	case steel:
		return "metal"
	case gold:
		return "metal"
	case cloth:
		return "cloth"
	case skin:
		return "skin"
	default:
		return ""
	}
}

func surfaceImage(kind string) *image.RGBA {
	const side = 64
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			n := int(noise2(x, y)%27) - 13
			value := 220 + n
			switch kind {
			case "wood":
				value = 217 + n/2 + ((x/5+y/23)%3)*9
				if x%17 == 0 || (x+y/3)%29 == 0 {
					value -= 36
				}
			case "stone":
				value = 215 + n
				if (x/12+y/15)%4 == 0 {
					value -= 22
				}
				if y%16 == 0 || (x+(y/16)*7)%23 == 0 {
					value -= 34
				}
			case "roof":
				value = 225 + n/2
				if y%10 < 2 {
					value -= 43
				}
				if (x+(y/10%2)*7)%16 < 2 {
					value -= 31
				}
			case "plaster":
				value = 224 + n/2
			case "cloth":
				value = 222 + n/3
				if x%4 == 0 || y%4 == 0 {
					value -= 15
				}
			case "metal":
				value = 220 + n/4
				if y%8 == 0 {
					value += 12
				}
			case "skin":
				value = 231 + n/4
			}
			if value < 90 {
				value = 90
			}
			if value > 255 {
				value = 255
			}
			v := uint8(value)
			img.SetRGBA(x, y, imagecolor.RGBA{v, v, v, 255})
		}
	}
	return img
}
