package main

import (
	"fmt"
	"math"
	"testing"
)

func TestRockHasNoOpenSeamsAboveGround(t *testing.T) {
	w := &objWriter{groupByMaterial: true, facesByMaterial: make(map[string][]triangleFace)}
	w.material("rock")
	w.boulder(0, 0, 0, 8, 15)
	edges := map[string]int{}
	bottom := map[string]bool{}
	for _, f := range w.facesByMaterial["rock"] {
		v := []point{f.a, f.b, f.c}
		for i := range v {
			a, b := v[i], v[(i+1)%3]
			ka, kb := fmt.Sprintf("%.5f,%.5f,%.5f", a.x+0.0000001, a.y+0.0000001, a.z), fmt.Sprintf("%.5f,%.5f,%.5f", b.x+0.0000001, b.y+0.0000001, b.z)
			if ka > kb {
				ka, kb = kb, ka
			}
			key := ka + "/" + kb
			edges[key]++
			bottom[key] = a.z == 0 && b.z == 0
		}
	}
	for edge, count := range edges {
		if count != 2 && !bottom[edge] {
			t.Errorf("fente dans le rocher : %s (%d faces)", edge, count)
		}
	}
}

func TestSegmentAndRoofOutwardFaces(t *testing.T) {
	w := &objWriter{groupByMaterial: true, facesByMaterial: make(map[string][]triangleFace)}
	w.material("stone")
	w.segment(0, 0, 10, 0, 2, 0, 1)
	faces := w.facesByMaterial["stone"]
	for i := 0; i < 2; i++ {
		f := faces[i]
		nz := (f.b.x-f.a.x)*(f.c.y-f.a.y) - (f.b.y-f.a.y)*(f.c.x-f.a.x)
		if nz <= 0 {
			t.Fatal("dessus du chemin orienté sous le sol")
		}
	}
	w.material("roof")
	w.gableRoof(0, 0, 8, 10, 12, 4)
	for _, f := range w.facesByMaterial["roof"][:4] {
		nz := (f.b.x-f.a.x)*(f.c.y-f.a.y) - (f.b.y-f.a.y)*(f.c.x-f.a.x)
		if nz <= 0 {
			t.Fatal("toiture orientée vers l'intérieur")
		}
	}
}

func TestGatesRemainWalkable(t *testing.T) {
	towers := gateTowerPositions()
	if len(towers) != 8 {
		t.Fatal("huit tours attendues")
	}
	for i := 0; i < 4; i++ {
		a := float64(i) * math.Pi / 2
		x, y := 82*math.Cos(a), 82*math.Sin(a)
		for _, p := range towers {
			if math.Hypot(x-p.x, y-p.y) < 5.5 {
				t.Fatal("une tour bouche une porte")
			}
		}
	}
}
