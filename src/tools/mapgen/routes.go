package main

import (
	"fmt"
	"math"
)

var mapRoutes [][][2]float64
var bridgePassages []collision
var generatedTrees []collision
var rivers = [][4]float64{{0, 300, 5, 105}, {5, 105, -75, 65}, {-75, 65, -102, 0}, {-102, 0, -92, -80}, {-92, -80, 0, -105}, {0, -105, 105, -88}, {105, -88, 145, -35}, {145, -35, 300, -55}}

func distanceSegment(x, y float64, s [4]float64) float64 {
	dx, dy := s[2]-s[0], s[3]-s[1]
	t := ((x-s[0])*dx + (y-s[1])*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(x-s[0]-t*dx, y-s[1]-t*dy)
}
func overWater(x, y float64) bool {
	for _, s := range rivers {
		if distanceSegment(x, y, s) < 12 {
			return true
		}
	}
	return false
}
func onRoad(x, y, margin float64) bool {
	for _, r := range mapRoutes {
		for i := 1; i < len(r); i++ {
			if distanceSegment(x, y, [4]float64{r[i-1][0], r[i-1][1], r[i][0], r[i][1]}) < margin {
				return true
			}
		}
	}
	return false
}

// Les tabliers remplacent le chemin au-dessus de l'eau. Même axe, même
// largeur, même hauteur de marche ; aucune coordonnée de pont indépendante.
func (w *objWriter) connectedRoads(routes [][][2]float64) {
	bridgePassages = nil
	for ri, r := range routes {
		for i := 1; i < len(r); i++ {
			a, b := r[i-1], r[i]
			length := math.Hypot(b[0]-a[0], b[1]-a[1])
			steps := int(math.Ceil(length))
			for j := 0; j < steps; j++ {
				t0, t1 := float64(j)/float64(steps), float64(j+1)/float64(steps)
				x0, y0 := a[0]+(b[0]-a[0])*t0, a[1]+(b[1]-a[1])*t0
				x1, y1 := a[0]+(b[0]-a[0])*t1, a[1]+(b[1]-a[1])*t1
				name := fmt.Sprintf("road_%d_%d_%d", ri, i, j)
				w.object(name)
				w.material("path")
				if overWater((x0+x1)/2, (y0+y1)/2) {
					w.material("wood")
					// Les petits rectangles se chevauchent pour une collision continue.
					passage := segmentCollision(name, "bridge", x0, y0, x1, y1, 8)
					passage.Width += 1
					bridgePassages = append(bridgePassages, passage)
				}
				w.segment(x0, y0, x1, y1, 8, .17, .03)
			}
		}
		// Jonctions rondes : comblent les angles entre tronçons.
		for i, p := range r {
			w.object(fmt.Sprintf("road_joint_%d_%d", ri, i))
			w.material("path")
			if overWater(p[0], p[1]) {
				w.material("wood")
			}
			w.cylinder(p[0], p[1], .17, 4, .03, 20)
		}
	}
}
