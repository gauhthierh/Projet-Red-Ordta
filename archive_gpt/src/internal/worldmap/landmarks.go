package worldmap

// StarterMap is the authored village hub. X grows east and Z grows south.
// The wilderness beyond it stays streamed and deterministic.
type Landmark struct {
	ID, Name     string
	X, Z, Radius float32
	HalfX, HalfZ float32 // non-zero for solid buildings and workstations
}

var StarterMap = []Landmark{
	{"arena", "Cercle d'entrainement", 0, 0, 4.4, 0, 0},
	{"market", "Marche", -5, -1, 3.2, 1.25, .65},
	{"forge", "Forge", 6, 1, 3.2, .9, .7},
	{"house_west", "Maison ouest", -6, -6, 0, 1.55, 1.35},
	{"house_north", "Maison nord", 0, -9, 0, 1.55, 1.35},
	{"house_east", "Maison est", 6, -6, 0, 1.55, 1.35},
	{"south_gate", "Porte sud", 0, 12, 0, 0, 0},
	{"west_trail", "Sentier des loups", -15, 4, 0, 0, 0},
	{"east_trail", "Route des collines", 15, 4, 0, 0, 0},
	{"wolf_camp", "Camp des loups", -22, 4, 2.5, 0, 0},
	{"hill_cache", "Cache des collines", 22, 4, 2.5, 0, 0},
	{"old_shrine", "Ancien sanctuaire", 0, 23, 2.5, 0, 0},
}

func LandmarkByID(id string) Landmark {
	for _, p := range StarterMap {
		if p.ID == id {
			return p
		}
	}
	return Landmark{}
}

// OnStarterRoad reserves the three deliberate exits from village clutter.
func OnStarterRoad(x, z float32) bool {
	return (x > -2.1 && x < 2.1 && z > -11 && z < 26) ||
		(z > 2.3 && z < 5.7 && x > -25 && x < 25)
}

func NearSite(x, z float32) string {
	for _, id := range []string{"wolf_camp", "hill_cache", "old_shrine"} {
		p := LandmarkByID(id)
		dx, dz := x-p.X, z-p.Z
		if dx*dx+dz*dz <= p.Radius*p.Radius {
			return id
		}
	}
	return ""
}
