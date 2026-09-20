package game

import "testing"

func countItem(s Snapshot, id string) int {
	for _, v := range s.Inventory {
		if v.ID == id {
			return v.Count
		}
	}
	return 0
}

func TestPlayableLoop(t *testing.T) {
	g := New()
	g.Player.Position = Vector2{X: -5, Z: 2}
	before := g.Backend.Snapshot()
	g.Dispatch(BuyPotion)
	after := g.Backend.Snapshot()
	if after.Gold != before.Gold-3 || countItem(after, "health_potion") != 4 {
		t.Fatal("purchase not reflected in inventory", after)
	}
	g.Player.Position = Vector2{X: 6, Z: 2}
	g.Dispatch(CraftHat)
	after = g.Backend.Snapshot()
	if countItem(after, "wolf_fur") != 2 || countItem(after, "raven_feather") != 0 || countItem(after, "boar_leather") != 0 || countItem(after, "adventurer_hat") != 1 || after.MaxHP != 100 {
		t.Fatal("craft did not place armor in bag", after)
	}
	for i, item := range after.Slots {
		if item.ID == "adventurer_hat" {
			if err := g.EquipSlot(i); err != nil {
				t.Fatal(err)
			}
			break
		}
	}
	after = g.Backend.Snapshot()
	if countItem(after, "adventurer_hat") != 0 || after.Equipment[0] != "adventurer_hat" || after.MaxHP != 110 {
		t.Fatal("armor did not equip", after)
	}
	g.Player.Position = Vector2{Z: 3.8}
	g.Dispatch(Train)
	for i := 0; i < 6 && g.Backend.Snapshot().Fighting; i++ {
		g.Dispatch(Attack)
	}
	after = g.Backend.Snapshot()
	if after.Fighting || after.Victories != 1 || after.Gold != before.Gold-3-5+8 {
		t.Fatal("combat reward incorrect", after)
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	d := NewDemo()
	s := d.Snapshot()
	s.Inventory[0].Count = 999
	if d.Snapshot().Potions == 999 {
		t.Fatal("shared mutable inventory")
	}
}

func TestForgeRequiresIngredients(t *testing.T) {
	g := New()
	g.Player.Position = Vector2{X: 6, Z: 2}
	before := g.Backend.Snapshot()
	g.Dispatch(CraftTunic)
	after := g.Backend.Snapshot()
	if after.Gold != before.Gold || countItem(after, "adventurer_tunic") != 0 {
		t.Fatal("crafted without materials")
	}
}

func TestExplorationRewardsOnce(t *testing.T) {
	g := New()
	g.Player.Position = Vector2{X: -22, Z: 4}
	before := g.Backend.Snapshot()
	g.Dispatch(ExploreWest)
	after := g.Backend.Snapshot()
	if countItem(after, "wolf_fur") != countItem(before, "wolf_fur")+2 || len(after.Discovered) != 1 {
		t.Fatal("exploration reward missing", after)
	}
	g.Dispatch(ExploreWest)
	if countItem(g.Backend.Snapshot(), "wolf_fur") != countItem(after, "wolf_fur") {
		t.Fatal("site rewarded twice")
	}
	g.Player.Position = Vector2{}
	g.Dispatch(ExploreEast)
	if len(g.Backend.Snapshot().Discovered) != 1 {
		t.Fatal("explored distant site")
	}
}
