package game3d

import "fmt"

func newHeroCharacter(class string, equipment []string) *character {
	a := newCharacter(false)
	switch class {
	case "Elfe":
		a.root.SetScale(1, 1.08, 1)
		for _, side := range []float32{-1, 1} {
			ear := cone(a.root, side*.37, 1.88, 0, .12, .32, skin)
			ear.SetRotationZ(-side * .52)
		}
	case "Nain":
		a.root.SetScale(1.13, .83, 1.1)
		box(a.root, 0, 1.59, .23, .30, .30, .18, wood)
	}
	for _, id := range equipment {
		switch id {
		case "adventurer_hat":
			box(a.root, 0, 2.11, 0, .73, .11, .63, wood)
			cone(a.root, 0, 2.34, 0, .35, .46, roof)
		case "adventurer_tunic":
			box(a.root, 0, 1.13, .22, .67, .77, .14, steel)
			box(a.root, 0, .83, .31, .54, .10, .08, gold)
		case "adventurer_boots":
			for i := range a.legs {
				box(a.legs[i], 0, -.58, .07, .32, .34, .42, wood)
			}
		}
	}
	return a
}

func (w *world) updateHeroAppearance() {
	s := w.game.Backend.Snapshot()
	stamp := fmt.Sprintf("%s/%v", s.Class, s.Equipment)
	if stamp == w.heroVisualStamp {
		return
	}
	w.heroVisualStamp = stamp
	if w.player != nil {
		w.scene.Remove(w.player)
		w.player.DisposeChildren(true)
	}
	w.hero = newHeroCharacter(s.Class, s.Equipment)
	w.player, w.sword = w.hero.root, w.hero.blade
	w.scene.Add(w.player)
}
