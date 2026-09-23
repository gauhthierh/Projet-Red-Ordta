package main

import (
	"math"
	"path/filepath"
	"testing"
)

func TestArenaAccessAndCombatFloor(t *testing.T) {
	writeWorld(filepath.Join(t.TempDir(), "world.obj"))
	obstacles := arenaDetailCollisions()
	touches := func(x, y float64, c collision) bool {
		if c.Shape == "circle" {
			return math.Hypot(x-c.X, y-c.Y) < c.Radius+.3
		}
		a := c.RotationDegrees * math.Pi / 180
		dx, dy := x-c.X, y-c.Y
		return math.Abs(dx*math.Cos(a)+dy*math.Sin(a)) < c.Width/2+.3 && math.Abs(-dx*math.Sin(a)+dy*math.Cos(a)) < c.Height/2+.3
	}
	// Vérifie le passage sur la largeur des deux routes d'arrivée.
	for _, route := range mapRoutes {
		for i := 1; i < len(route); i++ {
			a, b := route[i-1], route[i]
			length := math.Hypot(b[0]-a[0], b[1]-a[1])
			for d := 0.0; d <= length; d += .5 {
				x, y := a[0]+(b[0]-a[0])*d/length, a[1]+(b[1]-a[1])*d/length
				if math.Hypot(x-170, y+165) > 52 {
					continue
				}
				for _, offset := range []float64{-3, 0, 3} {
					px, py := x-(b[1]-a[1])*offset/length, y+(b[0]-a[0])*offset/length
					for _, c := range obstacles {
						if touches(px, py, c) {
							t.Fatalf("accès bloqué par %s à %.1f %.1f", c.ID, px, py)
						}
					}
				}
			}
		}
	}
	for _, p := range [][2]float64{{166, -165}, {174, -165}, {175, -162}, {175, -168}, {170, -180}} {
		for _, c := range obstacles {
			if touches(p[0], p[1], c) {
				t.Fatalf("emplacement combat bloqué : %s", c.ID)
			}
		}
	}
}
