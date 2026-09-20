package game

import "testing"

func TestTurnBased(t *testing.T) {
	g := New()
	g.Dispatch(Train)
	hp := g.Player.HP
	for i := 0; i < 100; i++ {
		g.Update(.1, Input{})
	}
	if g.Player.HP != hp {
		t.Fatal("damage without player action")
	}
	g.Dispatch(Attack)
	if g.Player.HP != hp-5 || g.Goblin.HP != 35 {
		t.Fatal("incorrect turn")
	}
}
func TestReturnPreservesPlayer(t *testing.T) {
	g := New()
	g.Player.Position = Vector2{X: -5, Z: 2}
	g.Dispatch(Buy)
	g.Player.Position = Vector2{X: 0, Z: 3.8}
	g.Dispatch(Train)
	g.Dispatch(Return)
	if g.Potions != 4 || g.Backend.Snapshot().Fighting {
		t.Fatal("return reset inventory")
	}
}

func TestVendorRequiresProximity(t *testing.T) {
	g := New()
	g.Dispatch(Buy)
	if g.Potions != 3 || g.Backend.Snapshot().Gold != 100 {
		t.Fatal("bought from outside the village stall")
	}
}
