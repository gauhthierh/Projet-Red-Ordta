package game3d

import (
	"fmt"

	"github.com/g3n/engine/gui"
	"projet-red/internal/game"
)

func (w *world) drawCharacterCreation(card *gui.Panel) {
	menuLabel(card, "Choisissez votre nom et votre classe", 28, 99, 17, color(.89, .91, .91))
	menuLabel(card, "Nom (lettres uniquement)", 28, 138, 14, color(.78, .84, .86))
	edit := gui.NewEdit(250, "Nom du personnage")
	edit.MaxLength = 20
	edit.SetPosition(210, 133)
	edit.SetFontSize(16)
	edit.SetText(w.nameDraft)
	card.Add(edit)
	menuLabel(card, "CLASSE", 28, 188, 14, color(.91, .79, .54))
	classes := []struct {
		name, id string
		hp       int
	}{{"Humain", "human", 100}, {"Elfe", "elf", 80}, {"Nain", "dwarf", 120}}
	for i, choice := range classes {
		x := float32(28 + i*201)
		panel := gui.NewPanel(183, 132)
		panel.SetPosition(x, 213)
		panel.SetColor4(rgba(.10, .14, .18, 1))
		panel.SetBorders(2, 2, 2, 2)
		border := color(.36, .44, .48)
		if choice.name == w.classDraft {
			border = color(.98, .75, .28)
		}
		panel.SetBordersColor(border)
		card.Add(panel)
		w.icon(panel, choice.id, 9, 18, 68)
		menuLabel(panel, fmt.Sprintf("%d PV max", choice.hp), 81, 40, 13, color(.82, .87, .89))
		class := choice.name
		w.menuButton(panel, class, 9, 88, 165, 35, "", func() { w.nameDraft = edit.Text(); w.classDraft = class; w.showMenu(menuCreation) })
	}
	w.menuButton(card, "Commencer l'aventure", 205, 360, 240, 42, "", func() {
		w.nameDraft = edit.Text()
		if err := w.game.CreateCharacter(w.nameDraft, w.classDraft); err != nil {
			w.menuDirty = true
			return
		}
		w.characterCreated = true
		w.selectedSlot = -1
		w.updateHeroAppearance()
		w.showMenu(menuNone)
	})
	menuLabel(card, "Depart : niveau 1, 50% des PV, 100 or, 3 potions, Coup de poing.", 28, 404, 12, color(.71, .79, .82))
}

var _ game.CharacterCreator = (*game.Demo)(nil)
