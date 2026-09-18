// Package game contains the engine-independent gameplay rules.
package game

import "math"

const (
	arenaLimit       = float32(7)
	playerSpeed      = float32(4.4)
	goblinSpeed      = float32(1.35)
	playerRange      = float32(2.15)
	goblinRange      = float32(1.65)
	goblinAttackWait = float32(1.1)
)

// State describes the current round state.
type State int

const (
	Playing State = iota
	Won
	Lost
)

// Vector2 is a position on the arena floor.
type Vector2 struct {
	X float32
	Z float32
}

// Fighter stores the small set of statistics needed by the demo.
type Fighter struct {
	Name       string
	HP         int
	MaxHP      int
	Damage     int
	Position   Vector2
	LastFacing Vector2
}

// Input contains actions collected during one frame.
type Input struct {
	MoveX     float32
	MoveZ     float32
	Attack    bool
	UsePotion bool
	Restart   bool
}

// Game owns the round state and does not depend on the 3D renderer.
type Game struct {
	Player         Fighter
	Goblin         Fighter
	Potions        int
	State          State
	Message        string
	goblinCooldown float32
}

// New creates a fresh arena round.
func New() *Game {
	return &Game{
		Player: Fighter{
			Name:       "Aventurier",
			HP:         100,
			MaxHP:      100,
			Damage:     24,
			Position:   Vector2{X: 0, Z: 3.8},
			LastFacing: Vector2{X: 0, Z: -1},
		},
		Goblin: Fighter{
			Name:       "Gobelin d'entrainement",
			HP:         72,
			MaxHP:      72,
			Damage:     11,
			Position:   Vector2{X: 0, Z: -3.2},
			LastFacing: Vector2{X: 0, Z: 1},
		},
		Potions: 3,
		State:   Playing,
		Message: "Approchez-vous du gobelin et attaquez !",
	}
}

// Update advances the simulation by one frame.
func (g *Game) Update(delta float32, input Input) {
	if g.State != Playing {
		if input.Restart {
			*g = *New()
		}
		return
	}

	if delta > 0.1 {
		delta = 0.1
	}
	if g.goblinCooldown > 0 {
		g.goblinCooldown -= delta
	}

	g.movePlayer(delta, input.MoveX, input.MoveZ)

	if input.UsePotion {
		g.usePotion()
	}

	if input.Attack {
		g.playerAttack()
		if g.State != Playing {
			return
		}
	}

	g.updateGoblin(delta)
}

func (g *Game) movePlayer(delta, x, z float32) {
	length := length(x, z)
	if length > 0 {
		x /= length
		z /= length
		g.Player.LastFacing = Vector2{X: x, Z: z}
		g.Player.Position.X = clamp(g.Player.Position.X+x*playerSpeed*delta, -arenaLimit, arenaLimit)
		g.Player.Position.Z = clamp(g.Player.Position.Z+z*playerSpeed*delta, -arenaLimit, arenaLimit)
	}
}

func (g *Game) usePotion() {
	if g.Potions == 0 {
		g.Message = "Vous n'avez plus de potion."
		return
	}
	if g.Player.HP == g.Player.MaxHP {
		g.Message = "Vos points de vie sont deja au maximum."
		return
	}
	g.Potions--
	g.Player.HP += 35
	if g.Player.HP > g.Player.MaxHP {
		g.Player.HP = g.Player.MaxHP
	}
	g.Message = "Potion utilisee : +35 PV."
}

func (g *Game) playerAttack() {
	if g.distance() > playerRange {
		g.Message = "Le gobelin est trop loin."
		return
	}
	g.Goblin.HP -= g.Player.Damage
	if g.Goblin.HP <= 0 {
		g.Goblin.HP = 0
		g.State = Won
		g.Message = "Victoire ! Appuyez sur R pour rejouer."
		return
	}
	g.Message = "Coup d'epee : 24 degats."
}

func (g *Game) updateGoblin(delta float32) {
	dx := g.Player.Position.X - g.Goblin.Position.X
	dz := g.Player.Position.Z - g.Goblin.Position.Z
	distance := length(dx, dz)
	if distance > 0 {
		g.Goblin.LastFacing = Vector2{X: dx / distance, Z: dz / distance}
	}

	if distance > goblinRange {
		step := goblinSpeed * delta
		if step > distance-goblinRange {
			step = distance - goblinRange
		}
		g.Goblin.Position.X += g.Goblin.LastFacing.X * step
		g.Goblin.Position.Z += g.Goblin.LastFacing.Z * step
		return
	}

	if g.goblinCooldown > 0 {
		return
	}
	g.goblinCooldown = goblinAttackWait
	g.Player.HP -= g.Goblin.Damage
	if g.Player.HP <= 0 {
		g.Player.HP = 0
		g.State = Lost
		g.Message = "Vous etes vaincu. Appuyez sur R pour recommencer."
		return
	}
	g.Message = "Le gobelin inflige 11 degats."
}

func (g *Game) distance() float32 {
	return length(
		g.Player.Position.X-g.Goblin.Position.X,
		g.Player.Position.Z-g.Goblin.Position.Z,
	)
}

func length(x, z float32) float32 {
	return float32(math.Sqrt(float64(x*x + z*z)))
}

func clamp(value, minimum, maximum float32) float32 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
