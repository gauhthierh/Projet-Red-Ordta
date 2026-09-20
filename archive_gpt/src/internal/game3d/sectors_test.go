package game3d

import (
	"testing"

	"github.com/g3n/engine/core"
	"projet-red/internal/game"
	"projet-red/internal/worldmap"
)

func TestSectorBudget(t *testing.T) {
	center := worldmap.SectorID{X: -4, Z: 7}
	views := wantedSectors(center)
	if len(views) != 25 {
		t.Fatalf("loaded %d sectors, want 25", len(views))
	}
	near := 0
	for _, detailed := range views {
		if detailed {
			near++
		}
	}
	if near != 9 {
		t.Fatalf("%d detailed sectors, want 9", near)
	}
	if !views[center] || views[worldmap.SectorID{X: -2, Z: 9}] {
		t.Fatal("bad LOD assignment")
	}
}

func TestSectorLifecycle(t *testing.T) {
	w := &world{scene: core.NewNode(), game: game.New()}
	w.syncSectors()
	if len(w.sectors) != 25 {
		t.Fatalf("initial sector count: %d", len(w.sectors))
	}
	old := worldmap.SectorID{X: -2, Z: 0}
	if _, ok := w.sectors[old]; !ok {
		t.Fatal("initial sector missing")
	}
	w.game.Player.Position.X = worldmap.SectorSize * 4
	w.syncSectors()
	if len(w.sectors) != 25 {
		t.Fatalf("streaming leaked sectors: %d", len(w.sectors))
	}
	if _, ok := w.sectors[old]; ok {
		t.Fatal("distant sector was not unloaded")
	}
	if _, ok := w.sectors[worldmap.SectorID{X: 6, Z: 0}]; !ok {
		t.Fatal("new sector not loaded")
	}
}
