package monde

type Collision struct {
	Identifiant string  `json:"id"`
	Categorie   string  `json:"category"`
	Forme       string  `json:"shape"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
	Largeur     float32 `json:"width"`
	Hauteur     float32 `json:"height"`
	Rayon       float32 `json:"radius"`
	Rotation    float32 `json:"rotation_degrees"`
}

func estDansMap(x float32, y float32, taillemap float32, margeperso float32) bool {
	demiTaille := taillemap / 2
	limite := demiTaille - margeperso

	return x >= -limite && x <= limite && y >= -limite && y <= limite
}

func personnageToucheObstacleRond(xPersonnage float32, yPersonnage float32, rayonPersonnage float32, obstacles []Collision) bool {
	for _, obstacle := range obstacles {
		if obstacle.Forme != "circle" {
			continue
		}

		differenceX := xPersonnage - obstacle.X
		differenceY := yPersonnage - obstacle.Y

		distanceCarree := differenceX*differenceX +
			differenceY*differenceY

		rayonTotal := rayonPersonnage + obstacle.Rayon
		rayonTotalCarre := rayonTotal * rayonTotal

		if distanceCarree < rayonTotalCarre {
			return true
		}
	}

	return false
}

func personnageToucheRectangle(xPersonnage float32, yPersonnage float32, rayonPersonnage float32, obstacles []Collision) bool {
	for _, obstacle := range obstacles {
		if obstacle.Forme != "rectangle" {
			continue
		}

		// Les rectangles inclinés seront traités plus tard.
		if obstacle.Rotation != 0 {
			continue
		}

		demiLargeur := obstacle.Largeur / 2
		demiHauteur := obstacle.Hauteur / 2

		minX := obstacle.X - demiLargeur
		maxX := obstacle.X + demiLargeur
		minY := obstacle.Y - demiHauteur
		maxY := obstacle.Y + demiHauteur

		plusProcheX := xPersonnage
		plusProcheY := yPersonnage

		if plusProcheX < minX {
			plusProcheX = minX
		} else if plusProcheX > maxX {
			plusProcheX = maxX
		}

		if plusProcheY < minY {
			plusProcheY = minY
		} else if plusProcheY > maxY {
			plusProcheY = maxY
		}

		differenceX := xPersonnage - plusProcheX
		differenceY := yPersonnage - plusProcheY

		distanceCarree := differenceX*differenceX +
			differenceY*differenceY

		if distanceCarree < rayonPersonnage*rayonPersonnage {
			return true
		}
	}

	return false
}
