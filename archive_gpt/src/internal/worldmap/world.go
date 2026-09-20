// Package worldmap owns deterministic world layout independently of the renderer.
package worldmap

import (
	"math"
	"math/rand"
)

const SectorSize = 12

type SectorID struct{ X, Z int }
type FeatureKind uint8

const (
	Tree FeatureKind = iota
	Rock
	Shrub
)

type Biome uint8

const (
	Meadow Biome = iota
	Forest
	Highland
)

type Feature struct {
	Kind  FeatureKind
	X, Z  float32
	Scale float32
}

type World struct {
	seed  int64
	cache map[SectorID][]Feature
}

func New(seed int64) *World { return &World{seed: seed, cache: make(map[SectorID][]Feature)} }

// BiomeAt groups multiple sectors into one visual region.
func (w *World) BiomeAt(id SectorID) Biome {
	if id.X >= -1 && id.X <= 1 && id.Z >= -1 && id.Z <= 1 {
		return Meadow
	}
	x, z := floorDiv(id.X, 3), floorDiv(id.Z, 3)
	v := uint64(w.seed) ^ uint64(int64(x)*73856093) ^ uint64(int64(z)*19349663)
	return Biome((v ^ (v >> 7)) % 3)
}

func floorDiv(a, b int) int {
	q := a / b
	if a%b < 0 {
		q--
	}
	return q
}

func SectorAt(x, z float32) SectorID {
	return SectorID{int(math.Floor(float64(x) / SectorSize)), int(math.Floor(float64(z) / SectorSize))}
}

func (w *World) Features(id SectorID) []Feature {
	if features, ok := w.cache[id]; ok {
		return features
	}
	rng := rand.New(rand.NewSource(w.seed ^ (int64(id.X) * 73856093) ^ (int64(id.Z) * 19349663)))
	biome := w.BiomeAt(id)
	count := 12
	if biome == Forest {
		count = 16
	}
	if biome == Highland {
		count = 10
	}
	features := make([]Feature, 0, count)
	for i := 0; i < count; i++ {
		x := float32(id.X*SectorSize) + 1 + float32(rng.Float64())*10
		z := float32(id.Z*SectorSize) + 1 + float32(rng.Float64())*10
		// The central village and its exits remain passable.
		if (math.Abs(float64(x)) < 11 && math.Abs(float64(z)) < 14) || OnStarterRoad(x, z) || NearSite(x, z) != "" {
			continue
		}
		roll := rng.Intn(10)
		kind := Tree
		switch biome {
		case Meadow:
			if roll >= 4 && roll < 7 {
				kind = Rock
			}
			if roll >= 7 {
				kind = Shrub
			}
		case Forest:
			if roll >= 7 && roll < 9 {
				kind = Rock
			}
			if roll >= 9 {
				kind = Shrub
			}
		case Highland:
			if roll >= 2 && roll < 9 {
				kind = Rock
			}
			if roll >= 9 {
				kind = Shrub
			}
		}
		features = append(features, Feature{Kind: kind, X: x, Z: z, Scale: .75 + float32(rng.Float64())*.65})
	}
	w.cache[id] = features
	return features
}

// Forget releases cached layout data when a sector leaves the visible set.
// Regeneration remains byte-for-byte deterministic from the seed and ID.
func (w *World) Forget(id SectorID) { delete(w.cache, id) }

// Walkable uses the same generated layout as the renderer. Radius is the
// character's horizontal footprint. Sliding is handled by the caller.
func (w *World) Walkable(x, z, radius float32) bool {
	// Buildings and stalls in the starter village.
	for _, b := range StarterMap {
		if b.HalfX == 0 {
			continue
		}
		bz := b.Z
		if b.ID == "market" || b.ID == "forge" {
			bz = 0
		} // visual stall/furnace centers
		if x > b.X-b.HalfX-radius && x < b.X+b.HalfX+radius && z > bz-b.HalfZ-radius && z < bz+b.HalfZ+radius {
			return false
		}
	}
	id := SectorAt(x, z)
	for dz := -1; dz <= 1; dz++ {
		for dx := -1; dx <= 1; dx++ {
			for _, f := range w.Features(SectorID{id.X + dx, id.Z + dz}) {
				if f.Kind == Shrub {
					continue
				}
				r := float32(.36) * f.Scale
				if f.Kind == Rock {
					r = .46 * f.Scale
				}
				vx, vz := x-f.X, z-f.Z
				if vx*vx+vz*vz < (radius+r)*(radius+r) {
					return false
				}
			}
		}
	}
	return true
}
