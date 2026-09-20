package game

import "testing"

func TestSectorAt(t *testing.T) {
	for _, test := range []struct {
		position Vector2
		want     SectorID
	}{
		{Vector2{0, 0}, SectorID{X: 0, Z: 0}},
		{Vector2{11.9, 12}, SectorID{X: 0, Z: 1}},
		{Vector2{-0.1, -12}, SectorID{X: -1, Z: -1}},
		{Vector2{-12.1, -12.1}, SectorID{X: -2, Z: -2}},
	} {
		if got := SectorAt(test.position); got != test.want {
			t.Errorf("SectorAt(%+v) = %+v, want %+v", test.position, got, test.want)
		}
	}
}

func TestExplorationExtendsBeyondOldArena(t *testing.T) {
	g := New()
	for i := 0; i < 100; i++ {
		g.Update(.1, Input{MoveX: 1})
	}
	if g.Player.Position.X <= 7 {
		t.Fatalf("unexpected exploration boundary: %v", g.Player.Position.X)
	}
}
