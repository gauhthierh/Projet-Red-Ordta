package main

// Positions communes au décor et aux collisions. Le marché reste libre de maisons.
var villageHouses = [][4]float64{
	{-49, 37, 17, 13}, {-40, 60, 14, 11}, {40, 60, 14, 11},
	{50, 34, 17, 13}, {-54, 12, 15, 12},
	{-50, -30, 18, 14}, {-27, -51, 15, 12}, {27, -51, 15, 12},
	{50, -29, 18, 14}, {-8, -56, 13, 10},
}

type zoneRock struct {
	name                 string
	x, y, radius, height float64
}

var pendingRocks []zoneRock
var generatedRocks []collision

// Attendre que tous les chemins soient connus avant de poser les rochers.
// Le même choix s'applique au maillage et à sa collision.
func (w *objWriter) zoneRocks() {
	for _, r := range pendingRocks {
		if onRoad(r.x, r.y, r.radius*1.1+4.5) {
			continue
		}
		w.object(r.name)
		w.material("rock")
		w.boulder(r.x, r.y, 0, r.radius, r.height)
		generatedRocks = append(generatedRocks, circleCollision(r.name, "rock", r.x, r.y, r.radius*1.1))
	}
}
