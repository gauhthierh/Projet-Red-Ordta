package game3d

import (
	"fmt"
	"strings"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/window"
	"projet-red/internal/game"
)

// drawInventory is a chest-style RPG inventory, with no on-screen hotbar.
// Crafted armor stays in the bag until the player equips it here.
func (w *world) drawInventory(card *gui.Panel, s game.Snapshot) {
	w.icon(card, strings.ToLower(s.Class), 28, 91, 70)
	menuLabel(card, fmt.Sprintf("%s  -  %s", s.Name, s.Class), 113, 101, 16, color(.96, .87, .68))
	menuLabel(card, fmt.Sprintf("Niveau %d   |   %d/%d PV", s.Level, s.HP, s.MaxHP), 113, 129, 13, color(.75, .83, .87))
	menuLabel(card, "ARMURE", 365, 90, 14, color(.88, .76, .51))
	armor := []struct{ id, label string }{{"adventurer_hat", "Tete"}, {"adventurer_tunic", "Torse"}, {"adventurer_boots", "Pieds"}}
	for i, slot := range armor {
		idx := i
		item := game.ItemStack{}
		if i < len(s.Equipment) && s.Equipment[i] != "" {
			item = game.ItemStack{ID: s.Equipment[i], Count: 1}
		}
		w.paintSlot(card, item, float32(361+i*61), 112, 48, false, func() {
			if w.selectedSlot >= 0 && w.selectedSlot < len(s.Slots) && s.Slots[w.selectedSlot].ID == slot.id {
				if w.game.EquipSlot(w.selectedSlot) == nil {
					w.selectedSlot = -1
				}
			} else if item.ID != "" {
				w.game.Unequip(idx)
			} else {
				w.game.Message = "Selectionnez une piece d'armure dans le sac."
			}
			w.menuDirty = true
		})
		menuLabel(card, slot.label, float32(366+i*61), 161, 10, color(.71, .77, .8))
	}
	menuLabel(card, "SAC D'AVENTURIER", 28, 174, 14, color(.88, .76, .51))
	for row := 0; row < 4; row++ {
		for col := 0; col < 9; col++ {
			idx := row*9 + col
			w.inventorySlot(card, s, idx, float32(28+col*53), float32(198+row*53))
		}
	}
	if w.selectedSlot >= 0 && w.selectedSlot < len(s.Slots) && s.Slots[w.selectedSlot].ID != "" {
		item := s.Slots[w.selectedSlot]
		w.icon(card, item.ID, 534, 207, 72)
		menuLabel(card, fmt.Sprintf("x%d", item.Count), 542, 275, 14, color(.98, .91, .72))
		for i, line := range wrappedItemName(game.ItemName(item.ID), 16) {
			if i > 1 {
				break
			}
			menuLabel(card, line, 514, float32(302+i*19), 12, color(.88, .9, .9))
		}
		switch item.ID {
		case "health_potion":
			w.menuButton(card, "Boire", 513, 346, 110, 32, "", func() { w.menuAction(game.Potion, menuInventory) })
		case "adventurer_hat", "adventurer_tunic", "adventurer_boots":
			w.menuButton(card, "Equiper", 513, 346, 110, 32, "", func() {
				if w.game.EquipSlot(w.selectedSlot) == nil {
					w.selectedSlot = -1
				}
				w.menuDirty = true
			})
		}
	} else {
		menuLabel(card, "Selectionnez", 516, 225, 12, color(.68, .76, .79))
		menuLabel(card, "un objet", 516, 244, 12, color(.68, .76, .79))
	}
	menuLabel(card, "Deux clics : deplacer ou echanger. Cliquez l'armure pour l'equiper/retirer.", 28, 421, 12, color(.73, .81, .84))
}

func wrappedItemName(name string, max int) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(name) {
		if len(line)+len(word)+1 > max && line != "" {
			lines = append(lines, line)
			line = word
		} else if line == "" {
			line = word
		} else {
			line += " " + word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func (w *world) inventorySlot(card *gui.Panel, s game.Snapshot, idx int, x, y float32) {
	item := game.ItemStack{}
	if idx < len(s.Slots) {
		item = s.Slots[idx]
	}
	w.paintSlot(card, item, x, y, 49, idx == w.selectedSlot, func() {
		if w.selectedSlot < 0 {
			if item.ID == "" {
				return
			}
			w.selectedSlot = idx
		} else if w.selectedSlot == idx {
			w.selectedSlot = -1
		} else {
			w.game.MoveSlot(w.selectedSlot, idx)
			w.selectedSlot = idx
		}
		w.menuDirty = true
	})
}

func (w *world) paintSlot(parent *gui.Panel, item game.ItemStack, x, y, size float32, selected bool, clicked func()) *gui.Panel {
	cell := gui.NewPanel(size, size)
	cell.SetPosition(x, y)
	cell.SetColor4(rgba(.10, .13, .17, 1))
	cell.SetBorders(2, 2, 2, 2)
	border := color(.32, .38, .43)
	if selected {
		border = color(1, .77, .27)
	}
	cell.SetBordersColor(border)
	parent.Add(cell)
	if item.ID != "" {
		w.icon(cell, item.ID, 4, 4, size-8)
		if item.Count > 1 {
			menuLabel(cell, fmt.Sprintf("%d", item.Count), size-17, size-19, 13, color(1, 1, 1))
		}
	}
	if clicked != nil {
		cell.Subscribe(gui.OnMouseUp, func(_ string, ev interface{}) {
			if ev.(*window.MouseEvent).Button == window.MouseButtonLeft {
				clicked()
			}
		})
	}
	return cell
}
