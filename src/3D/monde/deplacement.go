package monde

import "github.com/g3n/engine/math32"

// CalculerDirectionPersonnage transforme une commande locale
// en direction dans le monde.
func CalculerDirectionPersonnage(angleHorizontal float32, commandeAvant float32, commandeDroite float32) (float32, float32) {
	// Direction horizontale regardée par le personnage.
	avantX := math32.Sin(angleHorizontal)
	avantY := math32.Cos(angleHorizontal)

	// Direction située à droite du personnage.
	droiteX := math32.Cos(angleHorizontal)
	droiteY := -math32.Sin(angleHorizontal)

	// Combinaison des commandes.
	directionMondeX := avantX*commandeAvant + droiteX*commandeDroite

	directionMondeY := avantY*commandeAvant + droiteY*commandeDroite

	// Normalisation pour éviter une accélération en diagonale.
	longueur := math32.Sqrt(directionMondeX*directionMondeX + directionMondeY*directionMondeY)

	if longueur > 1 {
		directionMondeX /= longueur
		directionMondeY /= longueur
	}

	return directionMondeX, directionMondeY
}
