package worldmap

import (
	"reflect"
	"testing"
)

func TestGenerationIsStable(t *testing.T) {
	id := SectorID{4, -3}
	a, b := New(2026), New(2026)
	if !reflect.DeepEqual(a.Features(id), b.Features(id)) {
		t.Fatal("world layout changed")
	}
}
func TestBuildingsBlockMovement(t *testing.T) {
	w := New(2026)
	if w.Walkable(-6, -6, .24) {
		t.Fatal("house is traversable")
	}
	if !w.Walkable(0, 3.8, .24) {
		t.Fatal("spawn is blocked")
	}
}

func TestBiomeIsStableAndVillageIsMeadow(t *testing.T) {
	a, b := New(2026), New(2026)
	if a.BiomeAt(SectorID{}) != Meadow {
		t.Fatal("village biome changed")
	}
	for x := -8; x <= 8; x++ {
		for z := -8; z <= 8; z++ {
			id := SectorID{X: x, Z: z}
			if a.BiomeAt(id) != b.BiomeAt(id) {
				t.Fatalf("biome mismatch at %+v", id)
			}
		}
	}
}

func TestVisibleObstacleMatchesCollision(t *testing.T) {
	w := New(2026)
	for _, feature := range w.Features(SectorID{X: 4, Z: 4}) {
		if feature.Kind == Shrub {
			continue
		}
		if w.Walkable(feature.X, feature.Z, .2) {
			t.Fatalf("visible obstacle is traversable: %+v", feature)
		}
		return
	}
	t.Fatal("test sector has no solid feature")
}

func TestForgottenSectorRegeneratesIdentically(t *testing.T) {
	w := New(2026)
	id := SectorID{X: -5, Z: 6}
	before := append([]Feature(nil), w.Features(id)...)
	w.Forget(id)
	if !reflect.DeepEqual(before, w.Features(id)) {
		t.Fatal("sector changed after unload")
	}
}
