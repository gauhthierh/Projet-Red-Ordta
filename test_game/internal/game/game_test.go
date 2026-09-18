package game

import "testing"

func TestPlayerWinsInThreeHitsAtCloseRange(t *testing.T) {
	g := New()
	g.Player.Position = Vector2{X: 0, Z: 0}
	g.Goblin.Position = Vector2{X: 0, Z: 1}

	for range 3 {
		g.Update(0, Input{Attack: true})
	}

	if g.State != Won {
		t.Fatalf("expected a victory, got state %v with %d goblin HP", g.State, g.Goblin.HP)
	}
}

func TestPotionHealsWithoutExceedingMaximum(t *testing.T) {
	g := New()
	g.Player.HP = 90
	g.Update(0, Input{UsePotion: true})

	if g.Player.HP != 100 {
		t.Fatalf("expected 100 HP, got %d", g.Player.HP)
	}
	if g.Potions != 2 {
		t.Fatalf("expected 2 potions, got %d", g.Potions)
	}
}

func TestRestartCreatesFreshRound(t *testing.T) {
	g := New()
	g.State = Lost
	g.Player.HP = 0
	g.Update(0, Input{Restart: true})

	if g.State != Playing || g.Player.HP != g.Player.MaxHP || g.Potions != 3 {
		t.Fatalf("round was not reset: %+v", g)
	}
}
