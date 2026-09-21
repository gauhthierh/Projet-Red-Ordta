package monde

import "math"

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

		// On place temporairement le personnage dans le repère local du
		// rectangle. Le reste du calcul reste alors identique à celui d'un
		// rectangle droit, même pour un mur incliné.
		angle := obstacle.Rotation * math.Pi / 180
		cosinus := float32(math.Cos(float64(angle)))
		sinus := float32(math.Sin(float64(angle)))

		differenceCentreX := xPersonnage - obstacle.X
		differenceCentreY := yPersonnage - obstacle.Y
		positionLocaleX := differenceCentreX*cosinus + differenceCentreY*sinus
		positionLocaleY := -differenceCentreX*sinus + differenceCentreY*cosinus

		demiLargeur := obstacle.Largeur / 2
		demiHauteur := obstacle.Hauteur / 2

		minX := -demiLargeur
		maxX := demiLargeur
		minY := -demiHauteur
		maxY := demiHauteur

		plusProcheX := positionLocaleX
		plusProcheY := positionLocaleY

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

		differenceX := positionLocaleX - plusProcheX
		differenceY := positionLocaleY - plusProcheY

		distanceCarree := differenceX*differenceX +
			differenceY*differenceY

		if distanceCarree < rayonPersonnage*rayonPersonnage {
			return true
		}
	}

	return false
}

func positionAutorisee(x float32, y float32, rayonPersonnage float32, donnees DonneesMonde) bool {
	if !estDansMap(x, y, donnees.Taille, rayonPersonnage) {
		return false
	}

	if personnageToucheObstacleRond(x, y, rayonPersonnage, donnees.Obstacles) {
		return false
	}

	return !personnageToucheRectangle(x, y, rayonPersonnage, donnees.Obstacles)
}
