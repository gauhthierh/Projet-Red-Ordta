// Package game3d adapts the gameplay rules to a G3N scene.
package game3d

import (
	"fmt"
	"math"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"

	"projet-red/internal/game"
	"projet-red/internal/worldmap"
)

type healthBar struct {
	root     *core.Node
	fill     *graphic.Mesh
	material *material.Standard
}

type world struct {
	app              *app.Application
	scene            *core.Node
	camera           *camera.Camera
	game             *game.Game
	player           *core.Node
	goblin           *core.Node
	sword            *graphic.Mesh
	playerBar        healthBar
	goblinBar        healthBar
	status           *gui.Label
	message          *gui.Label
	hud              *gui.Panel
	menuRoot         *gui.Panel
	menu             menuKind
	previousMenu     menuKind
	menuDirty        bool
	selectedSlot     int
	characterCreated bool
	nameDraft        string
	classDraft       string
	previousText     string
	attackAnim       float32
	goblinAnim       float32
	previousHP       int
	elapsed          float32
	held             map[window.Key]bool
	keyEvents        map[window.Key]bool
	hero             *character
	heroVisualStamp  string
	enemy            *character
	lastPosition     game.Vector2
	sectors          map[worldmap.SectorID]sectorView
	currentSector    worldmap.SectorID
	playerShadow     *graphic.Mesh
	goblinShadow     *graphic.Mesh
}

// Run opens the window and starts the playable arena.
func Run() {
	w := newWorld()
	w.app.Run(func(rend *renderer.Renderer, delta time.Duration) {
		w.update(float32(delta.Seconds()))
		w.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		rend.Render(w.scene, w.camera)
	})
}

func newWorld() *world {
	a := app.App()
	if desktop, ok := a.IWindow.(*window.GlfwWindow); ok {
		desktop.SetTitle("Projet RED - Village 3D")
		desktop.SetSwapInterval(1)
	}

	w := &world{
		app:          a,
		scene:        core.NewNode(),
		game:         game.New(),
		held:         make(map[window.Key]bool),
		keyEvents:    make(map[window.Key]bool),
		selectedSlot: -1,
		classDraft:   "Humain",
	}
	a.Subscribe(window.OnKeyDown, func(_ string, event interface{}) {
		w.keyEvents[event.(*window.KeyEvent).Key] = true
	})
	w.previousHP = w.game.Player.HP

	gui.Manager().Set(w.scene)
	w.createCamera()
	w.createTerrain()
	w.createVillage()
	w.syncSectors()
	w.createLights()
	w.hero = newCharacter(false)
	w.enemy = newCharacter(true)
	w.player, w.sword = w.hero.root, w.hero.blade
	w.goblin = w.enemy.root
	w.scene.Add(w.player)
	w.scene.Add(w.goblin)
	w.playerShadow = shadow(w.scene, 0, 0, .85, .45)
	w.goblinShadow = shadow(w.scene, 0, 0, .75, .4)
	w.playerBar = createHealthBar(color(0.16, 0.85, 0.36))
	w.goblinBar = createHealthBar(color(0.9, 0.2, 0.16))
	w.scene.Add(w.playerBar.root)
	w.scene.Add(w.goblinBar.root)
	w.createHUD()
	w.showMenu(menuHome)
	w.refreshMenu()
	w.syncScene()

	a.Gls().ClearColor(0.38, 0.53, 0.65, 1)
	return w
}

func (w *world) createCamera() {
	w.camera = camera.New(1)
	w.camera.SetPosition(0, 5.7, 8.5)
	target := math32.Vector3{X: 0, Y: 0.5, Z: 0}
	up := math32.Vector3{X: 0, Y: 1, Z: 0}
	w.camera.LookAt(&target, &up)
	w.scene.Add(w.camera)

	onResize := func(_ string, _ interface{}) {
		width, height := w.app.GetSize()
		if height == 0 {
			return
		}
		w.app.Gls().Viewport(0, 0, int32(width), int32(height))
		w.camera.SetAspect(float32(width) / float32(height))
		if w.hud != nil {
			w.hud.SetWidth(float32(width - 32))
		}
		if w.menu != menuNone {
			w.menuDirty = true
		}
	}
	w.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)
}

func (w *world) createTerrain() {
	// The training ring is a landmark, not a wall around the village.
	markerMaterial := material.NewStandard(color(0.55, 0.14, 0.12))
	marker := graphic.NewMesh(geometry.NewTorus(2.2, 0.06, 8, 48, math32.Pi*2), markerMaterial)
	marker.SetRotationX(math32.Pi / 2)
	marker.SetPosition(0, 0.03, 0)
	w.scene.Add(marker)
}

func (w *world) createLights() {
	w.scene.Add(light.NewAmbient(color(1, 1, 1), 0.55))
	key := light.NewDirectional(color(1, 0.9, 0.75), 1.25)
	key.SetPosition(6, 10, 7)
	w.scene.Add(key)
	fill := light.NewPoint(color(0.25, 0.45, 1), 2.2)
	fill.SetPosition(-5, 4, 3)
	w.scene.Add(fill)
}

func (w *world) createHUD() {
	w.hud = gui.NewPanel(768, 106)
	w.hud.SetPosition(16, 16)
	w.hud.SetColor4(&math32.Color4{R: 0.03, G: 0.045, B: 0.07, A: 0.9})
	w.hud.SetBorders(1, 1, 1, 1)
	w.hud.SetBordersColor(color(0.35, 0.55, 0.8))
	w.scene.Add(w.hud)

	w.status = gui.NewLabel("")
	w.status.SetPosition(14, 10)
	w.status.SetFontSize(16)
	w.status.SetColor(color(0.9, 0.95, 1))
	w.hud.Add(w.status)

	w.message = gui.NewLabel("")
	w.message.SetPosition(14, 39)
	w.message.SetFontSize(15)
	w.message.SetColor(color(1, 0.78, 0.28))
	w.hud.Add(w.message)

	controls := gui.NewLabel("ZQSD: marcher | M marchand | F forge | T arene | I sac | C perso | B carte | G fouiller")
	controls.SetPosition(14, 71)
	controls.SetFontSize(13)
	controls.SetColor(color(0.7, 0.76, 0.84))
	w.hud.Add(controls)
}

func (w *world) update(delta float32) {
	keys := w.app.KeyState()
	if w.justPressed(window.KeyEscape) {
		if w.menu == menuCreation {
			w.showMenu(menuHome)
		} else if w.menu == menuHome && !w.characterCreated {
			w.showMenu(menuCreation)
		} else if w.menu == menuHome {
			w.showMenu(menuNone)
		} else if w.menu != menuNone {
			w.showMenu(menuNone)
		} else {
			w.showMenu(menuHome)
		}
	}
	if !w.characterCreated && w.menu == menuNone {
		w.showMenu(menuHome)
	}
	if w.menu == menuNone {
		if w.justPressed(window.KeyI) {
			w.showMenu(menuInventory)
		}
		if w.justPressed(window.KeyC) {
			w.showMenu(menuCharacter)
		}
		if w.justPressed(window.KeyB) {
			w.previousMenu = menuNone
			w.showMenu(menuMap)
		}
		if w.justPressed(window.KeyM) {
			if w.game.CanInteract(game.Buy) {
				w.showMenu(menuMerchant)
			} else {
				w.game.Message = "Approchez-vous du marchand."
			}
		}
		if w.justPressed(window.KeyF) {
			if w.game.CanInteract(game.Craft) {
				w.showMenu(menuForge)
			} else {
				w.game.Message = "Approchez-vous de la forge."
			}
		}
		if w.justPressed(window.KeyT) {
			if w.game.CanInteract(game.Train) {
				w.game.Dispatch(game.Train)
				w.showMenu(menuCombat)
			} else {
				w.game.Message = "Approchez-vous de l'arene."
			}
		}
		if w.justPressed(window.KeyG) {
			switch worldmap.NearSite(w.game.Player.Position.X, w.game.Player.Position.Z) {
			case "wolf_camp":
				w.game.Dispatch(game.ExploreWest)
			case "hill_cache":
				w.game.Dispatch(game.ExploreEast)
			case "old_shrine":
				w.game.Dispatch(game.ExploreSouth)
			default:
				w.game.Message = "Aucun lieu a fouiller ici. Consultez la carte (B)."
			}
		}
	}

	moveX, moveZ := float32(0), float32(0)
	if keys.Pressed(window.KeyD) || keys.Pressed(window.KeyRight) {
		moveX++
	}
	if keys.Pressed(window.KeyA) || keys.Pressed(window.KeyQ) || keys.Pressed(window.KeyLeft) {
		moveX--
	}
	if keys.Pressed(window.KeyS) || keys.Pressed(window.KeyDown) {
		moveZ++
	}
	if keys.Pressed(window.KeyW) || keys.Pressed(window.KeyZ) || keys.Pressed(window.KeyUp) {
		moveZ--
	}

	input := game.Input{
		MoveX: moveX, MoveZ: moveZ,
		Attack:    w.justPressed(window.KeySpace),
		UsePotion: w.justPressed(window.KeyE),
		Restart:   w.justPressed(window.KeyR),
	}

	if w.menu != menuNone {
		input = game.Input{}
	} else if w.game.Backend.Snapshot().Fighting && input.Attack {
		w.showMenu(menuCombat)
	}
	w.game.Update(delta, input)
	if w.game.Backend.Snapshot().Fighting && w.menu == menuNone {
		w.showMenu(menuCombat)
	}
	w.syncSectors()
	if input.Attack {
		w.attackAnim = 0.28
	}
	if w.game.Player.HP < w.previousHP {
		w.goblinAnim = 0.25
	}
	if input.Restart {
		w.previousHP = w.game.Player.HP
	}
	w.previousHP = w.game.Player.HP
	w.elapsed += delta
	w.updateAnimations(delta)
	w.syncScene()
	w.refreshMenu()
	clear(w.keyEvents)
}

func (w *world) justPressed(key window.Key) bool {
	pressed := w.app.KeyState().Pressed(key)
	fresh := w.keyEvents[key] || (pressed && !w.held[key])
	w.held[key] = pressed
	return fresh
}

func (w *world) updateAnimations(delta float32) {
	moving := w.game.Player.Position != w.lastPosition
	w.lastPosition = w.game.Player.Position
	w.hero.animate(delta, moving, w.attackAnim)
	w.enemy.animate(delta, false, w.goblinAnim)
	if w.attackAnim > 0 {
		w.attackAnim -= delta
		progress := 1 - w.attackAnim/0.28
		w.sword.SetRotationZ(-0.55 + float32(math.Sin(float64(progress*math32.Pi)))*1.5)
	} else {
		w.sword.SetRotationZ(-0.55)
	}
	if w.goblinAnim > 0 {
		w.goblinAnim -= delta
	}
}

func (w *world) syncScene() {
	w.updateHeroAppearance()
	player := w.game.Player
	goblin := w.game.Goblin
	w.player.SetPosition(player.Position.X, 0, player.Position.Z)
	w.playerShadow.SetPosition(player.Position.X, -.012, player.Position.Z)
	w.goblinShadow.SetPosition(goblin.Position.X, -.012, goblin.Position.Z)
	w.camera.SetPosition(player.Position.X, 5.7, player.Position.Z+8.5)
	w.camera.LookAt(&math32.Vector3{X: player.Position.X, Y: 1, Z: player.Position.Z - 2.2}, &math32.Vector3{Y: 1})
	w.goblin.SetPosition(goblin.Position.X, 0, goblin.Position.Z)
	w.player.SetRotationY(facingAngle(player.LastFacing))
	w.goblin.SetRotationY(facingAngle(goblin.LastFacing))

	playerBob := float32(math.Sin(float64(w.elapsed*7))) * 0.025
	goblinBob := float32(math.Sin(float64(w.elapsed*8+1))) * 0.045
	if w.goblinAnim > 0 {
		goblinBob += float32(math.Sin(float64((1-w.goblinAnim/0.25)*math32.Pi))) * 0.25
	}
	w.player.SetPositionY(playerBob)
	w.goblin.SetPositionY(goblinBob)

	updateHealthBar(&w.playerBar, player.HP, player.MaxHP, player.Position.X, player.Position.Z)
	updateHealthBar(&w.goblinBar, goblin.HP, goblin.MaxHP, goblin.Position.X, goblin.Position.Z)

	snapshot := w.game.Backend.Snapshot()
	status := fmt.Sprintf("%s (%s niv.%d)  %d/%d PV   |   %d or   |   Potions  %d", snapshot.Name, snapshot.Class, snapshot.Level, player.HP, player.MaxHP, snapshot.Gold, w.game.Potions)
	if worldmap.NearSite(player.Position.X, player.Position.Z) != "" {
		status += "   |   G fouiller"
	}
	combined := status + "\n" + w.game.Message
	if combined != w.previousText {
		w.status.SetText(status)
		w.message.SetText(w.game.Message)
		w.previousText = combined
	}
}

func createHealthBar(fillColor *math32.Color) healthBar {
	root := core.NewNode()
	background := graphic.NewMesh(geometry.NewBox(1.9, 0.18, 0.12), material.NewStandard(color(0.025, 0.025, 0.03)))
	root.Add(background)
	fillMaterial := material.NewStandard(fillColor)
	fill := graphic.NewMesh(geometry.NewBox(1.72, 0.1, 0.14), fillMaterial)
	fill.SetPosition(0, 0, -0.02)
	root.Add(fill)
	return healthBar{root: root, fill: fill, material: fillMaterial}
}

func updateHealthBar(bar *healthBar, hp, maxHP int, x, z float32) {
	ratio := float32(hp) / float32(maxHP)
	if ratio < 0.001 {
		ratio = 0.001
	}
	bar.root.SetPosition(x, 2.35, z)
	bar.fill.SetScaleX(ratio)
	bar.fill.SetPositionX(0.86 * (ratio - 1))
	if ratio > 0.55 {
		bar.material.SetColor(color(0.16, 0.85, 0.36))
	} else if ratio > 0.25 {
		bar.material.SetColor(color(0.95, 0.66, 0.12))
	} else {
		bar.material.SetColor(color(0.92, 0.15, 0.12))
	}
}

func facingAngle(direction game.Vector2) float32 {
	return float32(math.Atan2(float64(direction.X), float64(direction.Z)))
}

func color(red, green, blue float32) *math32.Color {
	return &math32.Color{R: red, G: green, B: blue}
}
