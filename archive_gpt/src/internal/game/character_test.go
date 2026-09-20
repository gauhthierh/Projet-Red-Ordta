package game

import "testing"

func TestCharacterCreationFromRequirements(t *testing.T) {
	for _, tc := range []struct {
		class string
		max   int
	}{{"Humain", 100}, {"Elfe", 80}, {"Nain", 120}} {
		d := NewDemo()
		if err := d.CreateCharacter("éLODIE", tc.class); err != nil {
			t.Fatal(err)
		}
		s := d.Snapshot()
		if s.Name != "Élodie" || s.Class != tc.class || s.Level != 1 || s.MaxHP != tc.max || s.HP != tc.max/2 || s.Gold != 100 || s.Potions != 3 || len(s.Skills) != 1 || s.Skills[0] != "Coup de poing" {
			t.Fatalf("bad %s creation: %+v", tc.class, s)
		}
		if countItem(s, "wolf_fur") != 0 {
			t.Fatal("new character retained demo materials")
		}
	}
}

func TestCharacterNameValidation(t *testing.T) {
	d := NewDemo()
	for _, name := range []string{"", "Jean 2", "Jean-Pierre", "Nom avec espace"} {
		if err := d.CreateCharacter(name, "Humain"); err == nil {
			t.Fatalf("accepted invalid name %q", name)
		}
	}
	if d.Snapshot().Name != "Aventurier" {
		t.Fatal("failed creation mutated character")
	}
	if err := d.CreateCharacter("ALICE", "Mage"); err == nil {
		t.Fatal("accepted invalid class")
	}
}

func slotWithID(s Snapshot, id string) int {
	for i, item := range s.Slots {
		if item.ID == id {
			return i
		}
	}
	return -1
}

func TestArmorIsCraftedThenEquipped(t *testing.T) {
	d := NewDemo()
	if err := d.CreateCharacter("ALICE", "Elfe"); err != nil {
		t.Fatal(err)
	}
	d.add("raven_feather", 1)
	d.add("boar_leather", 2)
	d.add("wolf_fur", 3)
	d.add("troll_hide", 1)
	for _, tc := range []struct {
		action       Action
		id           string
		place, bonus int
	}{{CraftHat, "adventurer_hat", 0, 10}, {CraftTunic, "adventurer_tunic", 1, 25}, {CraftBoots, "adventurer_boots", 2, 15}} {
		if err := d.Execute(tc.action); err != nil {
			t.Fatal(err)
		}
		s := d.Snapshot()
		if countItem(s, tc.id) != 1 || s.Equipment[tc.place] != "" {
			t.Fatalf("%s should enter bag first", tc.id)
		}
		if err := d.EquipSlot(slotWithID(s, tc.id)); err != nil {
			t.Fatal(err)
		}
		s = d.Snapshot()
		if countItem(s, tc.id) != 0 || s.Equipment[tc.place] != tc.id {
			t.Fatalf("%s did not move to equipment", tc.id)
		}
	}
	s := d.Snapshot()
	if s.MaxHP != 130 || s.HP != 40 {
		t.Fatalf("armor bonuses incorrect: %d/%d", s.HP, s.MaxHP)
	}
	if err := d.Unequip(1); err != nil {
		t.Fatal(err)
	}
	s = d.Snapshot()
	if s.MaxHP != 105 || s.Equipment[1] != "" || countItem(s, "adventurer_tunic") != 1 {
		t.Fatal("unequip did not restore item and health cap")
	}
}
