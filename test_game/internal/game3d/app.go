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
)

type healthBar struct {
	root     *core.Node
	fill     *graphic.Mesh
	material *material.Standard
}

type world struct {
	app          *app.Application
	scene        *core.Node
	camera       *camera.Camera
	game         *game.Game
	player       *core.Node
	goblin       *core.Node
	sword        *graphic.Mesh
	playerBar    healthBar
	goblinBar    healthBar
	status       *gui.Label
	message      *gui.Label
	hud          *gui.Panel
	previousText string
	attackHeld   bool
	potionHeld   bool
	restartHeld  bool
	attackAnim   float32
	goblinAnim   float32
	previousHP   int
	elapsed      float32
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
		desktop.SetTitle("Projet RED - Arene 3D")
		desktop.SetSwapInterval(1)
	}

	w := &world{
		app:   a,
		scene: core.NewNode(),
		game:  game.New(),
	}
	w.previousHP = w.game.Player.HP

	gui.Manager().Set(w.scene)
	w.createCamera()
	w.createArena()
	w.createLights()
	w.player, w.sword = createPlayer()
	w.goblin = createGoblin()
	w.scene.Add(w.player)
	w.scene.Add(w.goblin)
	w.playerBar = createHealthBar(color(0.16, 0.85, 0.36))
	w.goblinBar = createHealthBar(color(0.9, 0.2, 0.16))
	w.scene.Add(w.playerBar.root)
	w.scene.Add(w.goblinBar.root)
	w.createHUD()
	w.syncScene()

	a.Gls().ClearColor(0.035, 0.055, 0.09, 1)
	return w
}

func (w *world) createCamera() {
	w.camera = camera.New(1)
	w.camera.SetPosition(0, 9.5, 12.5)
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
	}
	w.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)
}

func (w *world) createArena() {
	groundMaterial := material.NewStandard(color(0.12, 0.18, 0.22))
	ground := graphic.NewMesh(geometry.NewPlane(16, 16), groundMaterial)
	ground.SetRotationX(-math32.Pi / 2)
	w.scene.Add(ground)

	wallMaterial := material.NewStandard(color(0.28, 0.34, 0.4))
	for _, wall := range []struct {
		x, z, width, depth float32
	}{
		{0, -8, 16.5, 0.35},
		{0, 8, 16.5, 0.35},
		{-8, 0, 0.35, 16.5},
		{8, 0, 0.35, 16.5},
	} {
		mesh := graphic.NewMesh(geometry.NewBox(wall.width, 0.7, wall.depth), wallMaterial)
		mesh.SetPosition(wall.x, 0.35, wall.z)
		w.scene.Add(mesh)
	}

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

	controls := gui.NewLabel("ZQSD/WASD ou fleches : bouger   ESPACE : attaquer   E : potion   R : rejouer   ECHAP : quitter")
	controls.SetPosition(14, 71)
	controls.SetFontSize(13)
	controls.SetColor(color(0.7, 0.76, 0.84))
	w.hud.Add(controls)
}

func (w *world) update(delta float32) {
	keys := w.app.KeyState()
	if keys.Pressed(window.KeyEscape) {
		w.app.Exit()
		return
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

	attackNow := keys.Pressed(window.KeySpace)
	potionNow := keys.Pressed(window.KeyE)
	restartNow := keys.Pressed(window.KeyR)
	input := game.Input{
		MoveX:     moveX,
		MoveZ:     moveZ,
		Attack:    attackNow && !w.attackHeld,
		UsePotion: potionNow && !w.potionHeld,
		Restart:   restartNow && !w.restartHeld,
	}
	w.attackHeld = attackNow
	w.potionHeld = potionNow
	w.restartHeld = restartNow

	w.game.Update(delta, input)
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
}

func (w *world) updateAnimations(delta float32) {
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
	player := w.game.Player
	goblin := w.game.Goblin
	w.player.SetPosition(player.Position.X, 0, player.Position.Z)
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

	status := fmt.Sprintf("Aventurier  %d/%d PV   |   Gobelin  %d/%d PV   |   Potions  %d", player.HP, player.MaxHP, goblin.HP, goblin.MaxHP, w.game.Potions)
	combined := status + "\n" + w.game.Message
	if combined != w.previousText {
		w.status.SetText(status)
		w.message.SetText(w.game.Message)
		w.previousText = combined
	}
}

func createPlayer() (*core.Node, *graphic.Mesh) {
	root := core.NewNode()
	bodyMaterial := material.NewStandard(color(0.12, 0.42, 0.9))
	skinMaterial := material.NewStandard(color(0.93, 0.68, 0.48))
	metalMaterial := material.NewStandard(color(0.72, 0.8, 0.9))

	body := graphic.NewMesh(geometry.NewCylinder(0.46, 1.2, 16, 1, true, true), bodyMaterial)
	body.SetPosition(0, 0.65, 0)
	root.Add(body)
	head := graphic.NewMesh(geometry.NewSphere(0.38, 16, 12), skinMaterial)
	head.SetPosition(0, 1.55, 0)
	root.Add(head)

	sword := graphic.NewMesh(geometry.NewBox(0.12, 1.25, 0.12), metalMaterial)
	sword.SetPosition(0.7, 1.0, -0.05)
	sword.SetRotationZ(-0.55)
	root.Add(sword)
	return root, sword
}

func createGoblin() *core.Node {
	root := core.NewNode()
	green := material.NewStandard(color(0.25, 0.67, 0.2))
	darkGreen := material.NewStandard(color(0.12, 0.36, 0.1))

	body := graphic.NewMesh(geometry.NewCylinder(0.48, 1.0, 12, 1, true, true), darkGreen)
	body.SetPosition(0, 0.55, 0)
	root.Add(body)
	head := graphic.NewMesh(geometry.NewSphere(0.46, 14, 10), green)
	head.SetPosition(0, 1.35, 0)
	root.Add(head)
	for _, x := range []float32{-0.55, 0.55} {
		ear := graphic.NewMesh(geometry.NewCone(0.2, 0.65, 8, 1, true), green)
		ear.SetPosition(x, 1.42, 0)
		ear.SetRotationZ(-x * 1.2)
		root.Add(ear)
	}
	return root
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
