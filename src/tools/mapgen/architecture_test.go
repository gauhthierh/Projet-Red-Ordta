package main

import (
	"math"
	"testing"
)

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
