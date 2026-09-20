package game

import "testing"

func TestInventorySlotsMoveAndSwap(t *testing.T) {
	d := NewDemo()
	if got := d.Snapshot().Slots; len(got) != 36 || got[0].ID != "health_potion" || got[9].ID != "wolf_fur" {
		t.Fatal("initial slot layout incorrect")
	}
	if err := d.MoveSlot(0, 5); err != nil {
		t.Fatal(err)
	}
	if got := d.Snapshot().Slots; got[0].ID != "" || got[5].ID != "health_potion" || got[5].Count != 3 {
		t.Fatal("move failed")
	}
	if err := d.MoveSlot(5, 9); err != nil {
		t.Fatal(err)
	}
	if got := d.Snapshot().Slots; got[9].ID != "health_potion" || got[5].ID != "wolf_fur" {
		t.Fatal("swap failed")
	}
	if err := d.MoveSlot(0, 36); err == nil {
		t.Fatal("invalid move accepted")
	}
	if d.Snapshot().Potions != 3 {
		t.Fatal("moving slots changed item count")
	}
}

func TestSlotSnapshotCannotChangeBackend(t *testing.T) {
	d := NewDemo()
	s := d.Snapshot()
	s.Slots[0].ID = "poison_potion"
	if d.Snapshot().Slots[0].ID != "health_potion" {
		t.Fatal("slots share backend memory")
	}
}

func TestLastPotionClearsSlot(t *testing.T) {
	d := NewDemo()
	d.state.HP = 1
	d.add("health_potion", -2)
	if err := d.Execute(Potion); err != nil {
		t.Fatal(err)
	}
	if d.Snapshot().Potions != 0 || d.Snapshot().Slots[0].ID != "" {
		t.Fatal("empty stack remained in inventory")
	}
}
