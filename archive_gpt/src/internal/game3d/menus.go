package game3d

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"projet-red/internal/game"
	"projet-red/internal/worldmap"
)

type menuKind uint8

const (
	menuNone menuKind = iota
	menuHome
	menuCreation
	menuInventory
	menuMerchant
	menuForge
	menuCharacter
	menuCombat
	menuMap
)

func rgba(r, g, b, a float32) *math32.Color4 { return &math32.Color4{R: r, G: g, B: b, A: a} }

func (w *world) showMenu(kind menuKind) {
	w.menu = kind
	w.menuDirty = true
}

func (w *world) refreshMenu() {
	if !w.menuDirty {
		return
	}
	w.menuDirty = false
	if w.menuRoot != nil {
		w.scene.Remove(w.menuRoot)
		w.menuRoot.DisposeChildren(true)
		w.menuRoot = nil
	}
	if w.menu == menuNone {
		return
	}
	width, height := w.app.GetSize()
	if width < 700 || height < 520 {
		return
	}
	root := gui.NewPanel(float32(width), float32(height))
	root.SetColor4(rgba(.01, .02, .04, .70))
	w.scene.Add(root)
	w.menuRoot = root
	cardHeight := float32(450)
	if w.menu == menuInventory {
		cardHeight = 500
	}
	card := gui.NewPanel(650, cardHeight)
	card.SetPosition(float32(width)/2-325, float32(height)/2-cardHeight/2)
	card.SetColor4(rgba(.055, .075, .105, .98))
	card.SetBorders(2, 2, 2, 2)
	card.SetBordersColor(color(.65, .51, .29))
	root.Add(card)
	s := w.game.Backend.Snapshot()
	title := map[menuKind]string{menuHome: "PROJET RED", menuCreation: "CREER UN PERSONNAGE", menuInventory: "INVENTAIRE", menuMerchant: "LE MARCHAND", menuForge: "LA FORGE", menuCharacter: "PERSONNAGE", menuCombat: "COMBAT", menuMap: "CARTE DU MONDE"}[w.menu]
	menuLabel(card, title, 28, 22, 25, color(.98, .86, .60))
	if w.menu == menuCreation {
		menuLabel(card, "Une nouvelle aventure commence ici.", 28, 61, 15, color(.73, .82, .91))
	} else {
		menuLabel(card, fmt.Sprintf("%d / %d PV     %d or     %d potions", s.HP, s.MaxHP, s.Gold, s.Potions), 28, 61, 15, color(.73, .82, .91))
	}
	messageY := float32(406)
	if w.menu == menuInventory {
		messageY = 458
	} else if w.menu == menuCreation {
		messageY = 427
	}
	menuLabel(card, w.game.Message, 28, messageY, 13, color(.98, .80, .45))
	if w.menu != menuHome {
		w.menuButton(card, "Retour", 510, 18, 112, 35, "", func() {
			if w.menu == menuCreation || (w.menu == menuMap && w.previousMenu == menuHome) {
				w.showMenu(menuHome)
			} else if s.Fighting {
				w.showMenu(menuCombat)
			} else {
				w.showMenu(menuNone)
			}
		})
	}
	switch w.menu {
	case menuHome:
		w.icon(card, "main_menu", 28, 105, 110)
		menuLabel(card, "Un village, trois routes et une arene.", 170, 112, 17, color(.88, .91, .93))
		menuLabel(card, "Explorez, achetez, forgez, puis combattez.", 170, 142, 16, color(.7, .79, .86))
		if w.characterCreated {
			w.menuButton(card, "Continuer", 170, 200, 200, 45, "", func() { w.showMenu(menuNone) })
			w.menuButton(card, "Nouveau personnage", 390, 200, 200, 45, "", func() { w.showMenu(menuCreation) })
		} else {
			w.menuButton(card, "Creer personnage", 170, 200, 200, 45, "", func() { w.showMenu(menuCreation) })
			w.menuButton(card, "Voir la carte", 390, 200, 200, 45, "", func() { w.previousMenu = menuHome; w.showMenu(menuMap) })
		}
		w.menuButton(card, "Quitter", 170, 265, 420, 42, "", func() { w.app.Exit() })
	case menuCreation:
		w.drawCharacterCreation(card)
	case menuInventory:
		w.drawInventory(card, s)
	case menuMerchant:
		menuLabel(card, "Cliquez pour acheter. Prix en pieces d'or.", 28, 98, 16, color(.8, .86, .9))
		for i, offer := range []struct {
			id     string
			price  int
			action game.Action
		}{{"health_potion", 3, game.BuyPotion}, {"wolf_fur", 4, game.BuyWolfFur}, {"boar_leather", 3, game.BuyBoarLeather}, {"raven_feather", 1, game.BuyRavenFeather}, {"troll_hide", 7, game.BuyTrollHide}} {
			y := float32(121 + i*54)
			w.itemRow(card, offer.id, 0, 30, y)
			a := offer.action
			w.menuButton(card, fmt.Sprintf("Acheter - %d or", offer.price), 410, y+3, 205, 39, "gold_coin", func() { w.menuAction(a, menuMerchant) })
		}
	case menuForge:
		menuLabel(card, "Chaque creation coute 5 or. Equipez-la depuis le sac.", 28, 97, 15, color(.8, .86, .9))
		for i, recipe := range []struct {
			id, materials string
			action        game.Action
		}{{"adventurer_hat", "1 plume + 1 cuir", game.CraftHat}, {"adventurer_tunic", "2 fourrures + 1 peau de troll", game.CraftTunic}, {"adventurer_boots", "1 fourrure + 1 cuir", game.CraftBoots}} {
			y := float32(138 + i*83)
			w.itemRow(card, recipe.id, 0, 28, y)
			menuLabel(card, recipe.materials, 89, y+34, 13, color(.66, .76, .82))
			a := recipe.action
			w.menuButton(card, "Fabriquer", 433, y+8, 180, 40, "crafting", func() { w.menuAction(a, menuForge) })
		}
	case menuCharacter:
		w.icon(card, strings.ToLower(s.Class), 32, 119, 142)
		menuLabel(card, fmt.Sprintf("%s - %s", s.Name, s.Class), 190, 133, 21, color(.93, .9, .78))
		menuLabel(card, fmt.Sprintf("Vie : %d / %d", s.HP, s.MaxHP), 190, 174, 17, color(.8, .86, .91))
		menuLabel(card, fmt.Sprintf("Niveau %d  |  Victoires : %d", s.Level, s.Victories), 190, 204, 17, color(.8, .86, .91))
		menuLabel(card, "Sort : Coup de poing", 190, 234, 15, color(.74, .84, .88))
		menuLabel(card, "Equipement :", 30, 288, 16, color(.95, .82, .56))
		armorCount := 0
		for _, id := range s.Equipment {
			if id != "" {
				armorCount++
			}
		}
		if armorCount == 0 {
			menuLabel(card, "Aucun", 170, 288, 16, color(.7, .78, .82))
		}
		for i, id := range s.Equipment {
			if i > 2 {
				break
			}
			if id != "" {
				w.icon(card, id, 170+float32(i*70), 316, 48)
			}
		}
		w.menuButton(card, "Ouvrir l'inventaire", 395, 320, 222, 41, "", func() { w.showMenu(menuInventory) })
	case menuCombat:
		w.icon(card, "training_goblin", 32, 115, 120)
		menuLabel(card, fmt.Sprintf("Gobelin : %d PV", s.EnemyHP), 178, 131, 20, color(.98, .63, .53))
		armorCount := 0
		for _, id := range s.Equipment {
			if id != "" {
				armorCount++
			}
		}
		menuLabel(card, fmt.Sprintf("Tour %d - Attaque : %d degats", s.Turn, 5+armorCount*2), 178, 168, 16, color(.8, .86, .91))
		w.menuButton(card, "Attaquer", 30, 275, 180, 54, "basic_attack", func() { w.menuAction(game.Attack, menuCombat) })
		w.menuButton(card, fmt.Sprintf("Potion x%d", s.Potions), 228, 275, 180, 54, "health_potion", func() { w.menuAction(game.Potion, menuCombat) })
		w.menuButton(card, "Fuir", 426, 275, 180, 54, "defeat", func() { w.menuAction(game.Return, menuNone) })
	case menuMap:
		w.mapPanel(card)
	}
}

func (w *world) menuAction(action game.Action, back menuKind) {
	w.game.Dispatch(action)
	s := w.game.Backend.Snapshot()
	if (action == game.Attack || action == game.Potion) && !s.Fighting && back == menuCombat {
		back = menuNone
	}
	w.showMenu(back)
}

func menuLabel(parent *gui.Panel, text string, x, y, size float32, c *math32.Color) {
	l := gui.NewLabel(text)
	l.SetFontSize(float64(size))
	l.SetColor(c)
	l.SetPosition(x, y)
	parent.Add(l)
}

func iconPath(id string) string {
	file := id + ".png"
	wd, _ := os.Getwd()
	exe, _ := os.Executable()
	for _, start := range []string{filepath.Dir(exe), wd} {
		dir := start
		for i := 0; i < 6; i++ {
			path := filepath.Join(dir, "assets", "ui", "icons", "png", file)
			if _, err := os.Stat(path); err == nil {
				return path
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ""
}
func (w *world) icon(parent *gui.Panel, id string, x, y, size float32) {
	p := iconPath(id)
	if p == "" {
		return
	}
	im, err := gui.NewImage(p)
	if err != nil {
		return
	}
	im.SetSize(size, size)
	im.SetPosition(x, y)
	parent.Add(im)
}
func (w *world) itemRow(parent *gui.Panel, id string, count int, x, y float32) {
	cell := gui.NewPanel(284, 57)
	cell.SetPosition(x, y)
	cell.SetColor4(rgba(.11, .15, .19, .95))
	cell.SetBorders(1, 1, 1, 1)
	cell.SetBordersColor(color(.25, .34, .41))
	parent.Add(cell)
	w.icon(cell, id, 5, 4, 48)
	label := game.ItemName(id)
	if count > 0 {
		label = fmt.Sprintf("%s  x%d", label, count)
	}
	menuLabel(cell, label, 60, 17, 14, color(.93, .91, .84))
}
func (w *world) menuButton(parent *gui.Panel, title string, x, y, width, height float32, iconID string, click func()) {
	b := gui.NewButton(title)
	b.SetPosition(x, y)
	b.SetSize(width, height)
	b.Label.SetFontSize(16)
	_ = iconID
	b.Subscribe(gui.OnClick, func(_ string, _ interface{}) { click() })
	parent.Add(b)
}

func (w *world) mapPanel(parent *gui.Panel) {
	mapBG := gui.NewPanel(355, 284)
	mapBG.SetPosition(28, 105)
	mapBG.SetColor4(rgba(.11, .18, .16, 1))
	mapBG.SetBorders(1, 1, 1, 1)
	mapBG.SetBordersColor(color(.52, .62, .44))
	parent.Add(mapBG)
	for _, road := range []struct{ x, y, w, h float32 }{{171, 75, 8, 185}, {50, 146, 250, 8}} {
		p := gui.NewPanel(road.w, road.h)
		p.SetPosition(road.x, road.y)
		p.SetColor4(rgba(.55, .53, .43, 1))
		mapBG.Add(p)
	}
	// The overlay uses the same coordinates as collision and interaction logic.
	for _, p := range worldmap.StarterMap {
		x := float32(175) + p.X*5
		y := float32(130) + p.Z*5
		if x < 2 || x > 345 || y < 2 || y > 270 {
			continue
		}
		marker := gui.NewPanel(8, 8)
		marker.SetPosition(x, y)
		marker.SetColor4(rgba(.96, .77, .33, 1))
		mapBG.Add(marker)
		switch p.ID {
		case "market":
			menuLabel(mapBG, "Marche", x-42, y-17, 11, color(.96, .87, .68))
		case "forge":
			menuLabel(mapBG, "Forge", x+9, y-17, 11, color(.96, .87, .68))
		case "arena":
			menuLabel(mapBG, "Arene", x+10, y-8, 11, color(.96, .87, .68))
		case "wolf_camp":
			menuLabel(mapBG, "Loups", x-13, y+9, 11, color(.96, .87, .68))
		case "hill_cache":
			menuLabel(mapBG, "Cache", x-13, y+9, 11, color(.96, .87, .68))
		case "old_shrine":
			menuLabel(mapBG, "Sanctuaire", x+9, y-8, 11, color(.96, .87, .68))
		}
	}
	x := float32(175) + w.game.Player.Position.X*5
	y := float32(130) + w.game.Player.Position.Z*5
	player := gui.NewPanel(11, 11)
	player.SetPosition(x, y)
	player.SetColor4(rgba(.23, .68, 1, 1))
	mapBG.Add(player)
	menuLabel(parent, "O  Village et arene", 409, 126, 15, color(.9, .89, .75))
	menuLabel(parent, "<  Sentier des loups", 409, 167, 14, color(.78, .84, .83))
	menuLabel(parent, ">  Route des collines", 409, 201, 14, color(.78, .84, .83))
	menuLabel(parent, "v  Porte sud", 409, 235, 14, color(.78, .84, .83))
	menuLabel(parent, "Bleu : votre position", 409, 319, 13, color(.4, .75, 1))
	menuLabel(parent, "G : fouiller une destination", 409, 344, 13, color(.9, .76, .48))
}
