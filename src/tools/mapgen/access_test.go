package main

import (
	"math"
	"path/filepath"
	"testing"
)

func obstacleAt(x, y, margin float64, c collision) bool {
	if c.Shape == "circle" {
		return math.Hypot(x-c.X, y-c.Y) < c.Radius+margin
	}
	a := c.RotationDegrees * math.Pi / 180
	dx, dy := x-c.X, y-c.Y
	return math.Abs(dx*math.Cos(a)+dy*math.Sin(a)) < c.Width/2+margin && math.Abs(-dx*math.Sin(a)+dy*math.Cos(a)) < c.Height/2+margin
}

// Toutes les routes doivent conserver un couloir de marche de six mètres.
func TestAllRoadsAccessible(t *testing.T) {
	writeWorld(filepath.Join(t.TempDir(), "world.obj"))
	obstacles := worldCollisions()
	seen := map[string]bool{}
	for ri, route := range mapRoutes {
		for i := 1; i < len(route); i++ {
			a, b := route[i-1], route[i]
			length := math.Hypot(b[0]-a[0], b[1]-a[1])
			for d := 0.0; d <= length; d += .5 {
				for _, offset := range []float64{-3, 0, 3} {
					x := a[0] + (b[0]-a[0])*d/length - (b[1]-a[1])*offset/length
					y := a[1] + (b[1]-a[1])*d/length + (b[0]-a[0])*offset/length
					bridge := false
					for _, p := range bridgePassages {
						if obstacleAt(x, y, -.35, p) {
							bridge = true
							break
						}
					}
					for _, c := range obstacles {
						if c.Category == "river" && bridge {
							continue
						}
						if !seen[c.ID] && obstacleAt(x, y, .3, c) {
							seen[c.ID] = true
							t.Errorf("route %d bloquée par %s à %.1f %.1f", ri, c.ID, x, y)
						}
					}
				}
			}
		}
	}
}

func TestEveryZoneReachableFromSpawn(t *testing.T) {
	writeWorld(filepath.Join(t.TempDir(), "world.obj"))
	obstacles := worldCollisions()
	// Grille de deux mètres, avec une marge suffisante pour relier les cases.
	type cell struct{ x, y int }
	free := map[cell]bool{}
	for x := -135; x <= 135; x++ {
		for y := -135; y <= 135; y++ {
			px, py := float64(x*2), float64(y*2)
			bridge := false
			for _, p := range bridgePassages {
				if obstacleAt(px, py, -.4, p) {
					bridge = true
					break
				}
			}
			blocked := false
			for _, c := range obstacles {
				if c.Category == "river" && bridge {
					continue
				}
				if obstacleAt(px, py, 1.1, c) {
					blocked = true
					break
				}
			}
			free[cell{x, y}] = !blocked
		}
	}
	start := cell{0, -10}
	if !free[start] {
		t.Fatal("point de départ bloqué")
	}
	queue := []cell{start}
	seen := map[cell]bool{start: true}
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		for _, d := range []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := cell{p.x + d.x, p.y + d.y}
			if free[n] && !seen[n] {
				seen[n] = true
				queue = append(queue, n)
			}
		}
	}
	for name, p := range map[string]cell{
		"village": {0, 19}, "marchand": {17, 4}, "forgeron": {24, 2},
		"guilde": {-75, 65}, "forge": {-97, -12}, "arène": {85, -82},
		"sanctuaire": {0, -98}, "champs": {-85, -82}, "forêt": {-112, 41},
		"corbeaux": {0, 105}, "marais": {90, 85}, "mana": {95, -6},
	} {
		if !seen[p] {
			t.Errorf("zone inaccessible depuis le départ : %s", name)
		}
	}
}
