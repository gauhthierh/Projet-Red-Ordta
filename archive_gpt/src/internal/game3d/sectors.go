package game3d

import (
	"image"
	imagecolor "image/color"
	"math"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/texture"

	"projet-red/internal/worldmap"
)

type sectorView struct {
	root *core.Node
	near bool
}

const visibleSectorRadius = 2

// syncSectors changes the scene only when the player crosses a sector border.
// Nine near sectors carry detailed props; the outer ring uses cheaper meshes.
func (w *world) syncSectors() {
	center := worldmap.SectorAt(w.game.Player.Position.X, w.game.Player.Position.Z)
	if w.sectors != nil && center == w.currentSector {
		return
	}
	if w.sectors == nil {
		w.sectors = make(map[worldmap.SectorID]sectorView)
	}
	w.currentSector = center
	wanted := wantedSectors(center)
	for id, view := range w.sectors {
		near, keep := wanted[id]
		if keep && view.near == near {
			continue
		}
		w.scene.Remove(view.root)
		view.root.DisposeChildren(true)
		delete(w.sectors, id)
		if !keep {
			w.game.Map.Forget(id)
		}
	}
	for id, near := range wanted {
		if _, ok := w.sectors[id]; ok {
			continue
		}
		view := w.buildSector(id, near)
		w.sectors[id] = view
		w.scene.Add(view.root)
	}
}

func wantedSectors(center worldmap.SectorID) map[worldmap.SectorID]bool {
	wanted := make(map[worldmap.SectorID]bool, 25)
	for dz := -visibleSectorRadius; dz <= visibleSectorRadius; dz++ {
		for dx := -visibleSectorRadius; dx <= visibleSectorRadius; dx++ {
			wanted[worldmap.SectorID{X: center.X + dx, Z: center.Z + dz}] = absInt(dx) <= 1 && absInt(dz) <= 1
		}
	}
	return wanted
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (w *world) buildSector(id worldmap.SectorID, near bool) sectorView {
	root := core.NewNode()
	const size = worldmap.SectorSize
	x, z := float32(id.X*size), float32(id.Z*size)
	mat := material.NewStandard(color(1, 1, 1))
	mat.SetSpecularColor(color(.04, .04, .04))
	mat.AddTexture(texture.NewTexture2DFromRGBA(terrainTexture(id, w.game.Map.BiomeAt(id))))
	ground := graphic.NewMesh(geometry.NewPlane(size, size), mat)
	ground.SetRotationX(-math.Pi / 2)
	ground.SetPosition(x+size/2, -.045, z+size/2)
	root.Add(ground)
	for _, feature := range w.game.Map.Features(id) {
		if near {
			addDetailedFeature(root, feature)
		} else {
			addDistantFeature(root, feature)
		}
	}
	return sectorView{root: root, near: near}
}

func addDetailedFeature(root *core.Node, f worldmap.Feature) {
	s := f.Scale
	switch f.Kind {
	case worldmap.Tree:
		shadow(root, f.X, f.Z, .95*s, .65*s)
		cylinder(root, f.X, 1.0*s, f.Z, .19*s, 2*s, wood)
		cone(root, f.X, 2.15*s, f.Z, 1.05*s, 1.8*s, color(.15, .32, .12))
		cone(root, f.X, 2.80*s, f.Z, .78*s, 1.55*s, color(.19, .40, .14))
		cone(root, f.X, 3.35*s, f.Z, .48*s, 1.10*s, color(.25, .49, .18))
		for i := 0; i < 3; i++ {
			a := float64(i) * 2.094
			sphere(root, f.X+float32(math.Cos(a))*.3*s, .13*s, f.Z+float32(math.Sin(a))*.3*s, .15*s, color(.18, .32, .12))
		}
	case worldmap.Rock:
		shadow(root, f.X, f.Z, .7*s, .55*s)
		m := sphere(root, f.X, .27*s, f.Z, .52*s, color(.42, .43, .41))
		m.SetScale(1, .55, .8)
		m = sphere(root, f.X+.19*s, .38*s, f.Z-.1*s, .28*s, color(.52, .51, .45))
		m.SetScale(1, .55, .75)
	case worldmap.Shrub:
		for _, d := range [][2]float32{{-.2, 0}, {.15, .12}, {.04, -.2}} {
			sphere(root, f.X+d[0]*s, .20*s, f.Z+d[1]*s, .28*s, color(.21, .42, .16))
		}
	}
}

func addDistantFeature(root *core.Node, f worldmap.Feature) {
	s := f.Scale
	switch f.Kind {
	case worldmap.Tree:
		cylinder(root, f.X, .9*s, f.Z, .16*s, 1.8*s, wood)
		cone(root, f.X, 2.2*s, f.Z, .92*s, 2.7*s, color(.18, .37, .14))
	case worldmap.Rock:
		m := sphere(root, f.X, .24*s, f.Z, .48*s, stone)
		m.SetScale(1, .5, .9)
	case worldmap.Shrub:
		sphere(root, f.X, .19*s, f.Z, .32*s, color(.20, .38, .15)).SetScale(1, .55, 1)
	}
}

func terrainTexture(id worldmap.SectorID, biome worldmap.Biome) *image.RGBA {
	const pixels = 128
	img := image.NewRGBA(image.Rect(0, 0, pixels, pixels))
	for py := 0; py < pixels; py++ {
		for px := 0; px < pixels; px++ {
			wx := float32(id.X*worldmap.SectorSize) + float32(px)*worldmap.SectorSize/pixels
			wz := float32(id.Z*worldmap.SectorSize) + float32(py)*worldmap.SectorSize/pixels
			n := noise2(int(math.Floor(float64(wx*7))), int(math.Floor(float64(wz*7))))
			shade := int(n%29) - 14
			var c imagecolor.RGBA
			// A broad worn lane through the settlement, then irregular grass and soil.
			if math.Abs(float64(wx)) < 1.45 && wz > -13 && wz < 8 {
				c = imagecolor.RGBA{uint8(117 + shade), uint8(105 + shade), uint8(82 + shade), 255}
			} else {
				switch biome {
				case worldmap.Forest:
					c = imagecolor.RGBA{uint8(66 + shade), uint8(96 + shade), uint8(46 + shade), 255}
				case worldmap.Highland:
					c = imagecolor.RGBA{uint8(116 + shade), uint8(111 + shade), uint8(86 + shade), 255}
				default:
					c = imagecolor.RGBA{uint8(75 + shade), uint8(114 + shade), uint8(53 + shade), 255}
				}
				if n%17 < 3 {
					c.R += 10
					c.G += 2
					c.B += 2
				}
			}
			img.SetRGBA(px, py, c)
		}
	}
	return img
}

func noise2(x, z int) uint32 {
	v := uint32(x)*374761393 + uint32(z)*668265263 + 2026
	v = (v ^ (v >> 13)) * 1274126177
	return v ^ (v >> 16)
}
