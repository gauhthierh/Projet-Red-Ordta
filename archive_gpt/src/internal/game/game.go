// Package game defines the presentation contract independently of G3N or a CLI.
package game

import (
	"fmt"
	"math"
	"projet-red/internal/worldmap"
)

type State int

const (
	Playing State = iota
	Won
	Lost
)

type Vector2 struct{ X, Z float32 }
type Fighter struct {
	Name                 string
	HP, MaxHP, Damage    int
	Position, LastFacing Vector2
}
type Input struct {
	MoveX, MoveZ               float32
	Attack, UsePotion, Restart bool
}
type Snapshot struct {
	HP, MaxHP, Potions, Gold, EnemyHP, Turn int
	Name, Class                             string
	Level                                   int
	Fighting                                bool
	Message                                 string
	Inventory                               []ItemStack
	Slots                                   []ItemStack
	Equipment                               []string
	Skills                                  []string
	Victories                               int
	Discovered                              []string
}
type ItemStack struct {
	ID    string
	Count int
}
type Action string

const (
	Attack          Action = "attack"
	Potion          Action = "potion"
	Buy             Action = "buy"
	Craft           Action = "craft"
	Train           Action = "train"
	Return          Action = "return"
	BuyPotion       Action = "buy_potion"
	BuyWolfFur      Action = "buy_wolf_fur"
	BuyBoarLeather  Action = "buy_boar_leather"
	BuyRavenFeather Action = "buy_raven_feather"
	BuyTrollHide    Action = "buy_troll_hide"
	CraftHat        Action = "craft_hat"
	CraftTunic      Action = "craft_tunic"
	CraftBoots      Action = "craft_boots"
	ExploreWest     Action = "explore_west"
	ExploreEast     Action = "explore_east"
	ExploreSouth    Action = "explore_south"
)

// Backend is the only connection required by future business functions.
// Snapshot must return a value, never shared mutable presentation state.
type Backend interface {
	Snapshot() Snapshot
	Execute(Action) error
}

// SlotMover is optional so the later business backend can adopt the grid
// without coupling its rules to G3N.
type SlotMover interface{ MoveSlot(from, to int) error }
type CharacterCreator interface {
	CreateCharacter(name, class string) error
}
type EquipmentManager interface {
	EquipSlot(slot int) error
	Unequip(slot int) error
}

// Game owns spatial presentation state, not combat rules.
type Game struct {
	Player, Goblin Fighter
	Potions        int
	State          State
	Message        string
	Backend        Backend
	Map            *worldmap.World
}

func New() *Game { return NewWithBackend(NewDemo()) }
func NewWithBackend(b Backend) *Game {
	g := &Game{Backend: b, Map: worldmap.New(2026), Player: Fighter{Name: "Aventurier", Position: Vector2{Z: 3.8}, LastFacing: Vector2{Z: -1}}, Goblin: Fighter{Name: "Gobelin", Position: Vector2{Z: -3.2}, MaxHP: 40}}
	g.sync()
	return g
}
func (g *Game) sync() {
	s := g.Backend.Snapshot()
	g.Player.Name = s.Name
	g.Player.HP = s.HP
	g.Player.MaxHP = s.MaxHP
	g.Goblin.HP = s.EnemyHP
	g.Goblin.MaxHP = 40 + s.Victories*5
	g.Potions = s.Potions
	g.State = Playing
	g.Message = s.Message
	if s.Fighting {
		g.Message = "Combat tour par tour : " + g.Message
	}
}
func (g *Game) Dispatch(a Action) {
	if !g.CanInteract(a) {
		g.Message = fmt.Sprintf("Approchez-vous du lieu pour utiliser %s.", a)
		return
	}
	if err := g.Backend.Execute(a); err != nil {
		g.Message = err.Error()
		return
	}
	g.sync()
}

func (g *Game) MoveSlot(from, to int) {
	mover, ok := g.Backend.(SlotMover)
	if !ok {
		g.Message = "Le deplacement d'objets n'est pas disponible."
		return
	}
	if err := mover.MoveSlot(from, to); err != nil {
		g.Message = err.Error()
		return
	}
	g.sync()
}
func (g *Game) CreateCharacter(name, class string) error {
	creator, ok := g.Backend.(CharacterCreator)
	if !ok {
		return fmt.Errorf("Creation du personnage indisponible")
	}
	if err := creator.CreateCharacter(name, class); err != nil {
		g.Message = err.Error()
		return err
	}
	g.Player.Position = Vector2{Z: 3.8}
	g.sync()
	return nil
}
func (g *Game) EquipSlot(slot int) error {
	manager, ok := g.Backend.(EquipmentManager)
	if !ok {
		return fmt.Errorf("Equipement indisponible")
	}
	if err := manager.EquipSlot(slot); err != nil {
		g.Message = err.Error()
		return err
	}
	g.sync()
	return nil
}
func (g *Game) Unequip(slot int) error {
	manager, ok := g.Backend.(EquipmentManager)
	if !ok {
		return fmt.Errorf("Equipement indisponible")
	}
	if err := manager.Unequip(slot); err != nil {
		g.Message = err.Error()
		return err
	}
	g.sync()
	return nil
}

// CanInteract applies 3D spatial rules before forwarding an action to the
// temporary (or future) business backend. Inventory actions remain global.
func (g *Game) CanInteract(a Action) bool {
	if g.Backend.Snapshot().Fighting {
		return true
	}
	var x, z, r float32
	switch a {
	case Buy, BuyPotion, BuyWolfFur, BuyBoarLeather, BuyRavenFeather, BuyTrollHide:
		p := worldmap.LandmarkByID("market")
		x, z, r = p.X, p.Z, p.Radius
	case Craft, CraftHat, CraftTunic, CraftBoots:
		p := worldmap.LandmarkByID("forge")
		x, z, r = p.X, p.Z, p.Radius
	case Train:
		p := worldmap.LandmarkByID("arena")
		x, z, r = p.X, p.Z, p.Radius
	case ExploreWest, ExploreEast, ExploreSouth:
		id := map[Action]string{ExploreWest: "wolf_camp", ExploreEast: "hill_cache", ExploreSouth: "old_shrine"}[a]
		p := worldmap.LandmarkByID(id)
		x, z, r = p.X, p.Z, p.Radius
	default:
		return true
	}
	dx, dz := g.Player.Position.X-x, g.Player.Position.Z-z
	return dx*dx+dz*dz <= r*r
}
func (g *Game) Update(dt float32, in Input) {
	if dt < 0 {
		dt = 0
	}
	if dt > 0.1 {
		dt = 0.1
	}
	s := g.Backend.Snapshot()
	if !s.Fighting {
		l := float32(math.Hypot(float64(in.MoveX), float64(in.MoveZ)))
		if l > 0 {
			d := Vector2{in.MoveX / l, in.MoveZ / l}
			g.Player.LastFacing = d
			nextX := g.Player.Position.X + d.X*4.4*dt
			nextZ := g.Player.Position.Z + d.Z*4.4*dt
			if g.Map.Walkable(nextX, g.Player.Position.Z, .24) {
				g.Player.Position.X = nextX
			}
			if g.Map.Walkable(g.Player.Position.X, nextZ, .24) {
				g.Player.Position.Z = nextZ
			}
		}
	}
	if in.Attack {
		g.Dispatch(Attack)
	}
	if in.UsePotion {
		g.Dispatch(Potion)
	}
	if in.Restart {
		if g.Backend.Snapshot().Fighting {
			g.Dispatch(Return)
		}
	}
}
