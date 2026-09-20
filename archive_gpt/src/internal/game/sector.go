package game

import "projet-red/internal/worldmap"

// SectorSize is a world-space grid for future region streaming and saving.
// It does not load or unload G3N objects by itself.
const SectorSize = worldmap.SectorSize

type SectorID = worldmap.SectorID

func SectorAt(position Vector2) SectorID {
	return worldmap.SectorAt(position.X, position.Z)
}
