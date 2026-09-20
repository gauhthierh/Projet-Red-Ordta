package worldmap

import "testing"

func TestStarterLandmarks(t *testing.T) {
	for _, id := range []string{"market", "forge", "arena", "house_west", "south_gate", "west_trail", "east_trail", "wolf_camp", "hill_cache", "old_shrine"} {
		if LandmarkByID(id).ID != id {
			t.Fatalf("missing landmark %s", id)
		}
	}
	w := New(2026)
	for _, p := range [][2]float32{{0, 12}, {-15, 4}, {15, 4}, {-22, 4}, {22, 4}, {0, 23}} {
		if !w.Walkable(p[0], p[1], .24) {
			t.Fatalf("authored road blocked at %v", p)
		}
	}
}
