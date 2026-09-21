package main

import (
	"math"
	"testing"
)

func TestFenceCrossingOpensVisibleAndCollisionGap(t *testing.T) {
	old := mapRoutes
	defer func() { mapRoutes = old }()
	mapRoutes = [][][2]float64{{{-30, 3}, {30, 3}}}
	sections := fenceSections("test", 0, 0, 20, 20, 4, "south")
	collisions := appendFenceCollisions(nil, "test", 0, 0, 20, 20, 4, "south")
	if len(sections) != len(collisions) {
		t.Fatal("dessin et collisions désynchronisés")
	}
	for i, s := range sections {
		for j := 0; j <= 100; j++ {
			f := float64(j) / 100
			x, y := s[0]+(s[2]-s[0])*f, s[1]+(s[3]-s[1])*f
			if onRoad(x, y, 5.19) {
				t.Fatalf("clôture sur le chemin : %v", s)
			}
		}
		if math.Abs(collisions[i].X-(s[0]+s[2])/2) > 1e-6 {
			t.Fatal("collision décalée")
		}
	}
}

func TestBoundaryClosed(t *testing.T) {
	segments := boundarySegments()
	for i, s := range segments {
		next := segments[(i+1)%len(segments)]
		if math.Hypot(s[2]-next[0], s[3]-next[1]) > 1e-8 {
			t.Fatal("trou entre deux falaises")
		}
	}
}
