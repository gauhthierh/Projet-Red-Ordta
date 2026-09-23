package monde

import "math"

// Collision : Décrit une empreinte au sol : Hauteur est une dimension sur Y, pas une altitude Z. Rotation est en degrés.
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

// estDansMap : Garde tout le disque du personnage à l'intérieur des limites de la carte.
func estDansMap(x float32, y float32, taillemap float32, margeperso float32) bool {
	demiTaille := taillemap / 2
	limite := demiTaille - margeperso

	return x >= -limite && x <= limite && y >= -limite && y <= limite
}

// personnageToucheObstacleRond : Compare la distance entre centres à la somme des rayons dans le plan du sol.
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

// personnageToucheRectangle : Ramène le joueur dans le repère de chaque rectangle pour tester aussi les obstacles tournés.
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

		// Ramener le point dans les bornes donne le point du rectangle le plus
		// proche du joueur, y compris près d'un coin.
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

// positionAutorisee : Vérifie les limites et les obstacles ; un pont n'autorise que le franchissement de la rivière.
func positionAutorisee(x float32, y float32, rayonPersonnage float32, donnees DonneesMonde) bool {
	if !estDansMap(x, y, donnees.Taille, rayonPersonnage) {
		return false
	}

	if personnageToucheObstacleRond(x, y, rayonPersonnage, donnees.Obstacles) {
		return false
	}

	surPont := false
	for _, pont := range donnees.Ponts {
		angle := float64(pont.Rotation) * math.Pi / 180
		dx, dy := float64(x-pont.X), float64(y-pont.Y)
		localX := dx*math.Cos(angle) + dy*math.Sin(angle)
		localY := -dx*math.Sin(angle) + dy*math.Cos(angle)
		if math.Abs(localX) <= float64(pont.Largeur/2) && math.Abs(localY) <= float64(pont.Hauteur/2-rayonPersonnage) {
			surPont = true
			break
		}
	}
	for _, obstacle := range donnees.Obstacles {
		if obstacle.Categorie == "river" && surPont {
			continue
		}
		if personnageToucheRectangle(x, y, rayonPersonnage, []Collision{obstacle}) {
			return false
		}
	}
	return true
}

// Plusieurs petits pas évitent de sauter par-dessus un tronc lors d'une
// image lente. On conserve le glissement le long des obstacles.
func deplacerAvecCollisions(x, y, dx, dy, rayon float32, donnees DonneesMonde) (float32, float32) {
	pas := int(math.Ceil(math.Hypot(float64(dx), float64(dy)) / .15))
	if pas < 1 {
		return x, y
	}
	dx /= float32(pas)
	dy /= float32(pas)
	for i := 0; i < pas; i++ {
		if positionAutorisee(x+dx, y, rayon, donnees) {
			x += dx
		}
		if positionAutorisee(x, y+dy, rayon, donnees) {
			y += dy
		}
	}
	return x, y
}
