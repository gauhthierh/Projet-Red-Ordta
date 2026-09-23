package personnages

import "math"

// Animer : Anime la marche avec delta en secondes ; au repos, les membres reviennent progressivement à zéro.
func (p *Personnage) Animer(delta float32, bouge bool) {
	if !bouge {
		angleActuelbd := p.brasDroit.Rotation().X
		nouvelAnglebd := angleadouci(angleActuelbd, 0, 8, delta)
		p.brasDroit.SetRotationX(nouvelAnglebd)

		angleActuelbg := p.brasGauche.Rotation().X
		nouvelAnglebg := angleadouci(angleActuelbg, 0, 8, delta)
		p.brasGauche.SetRotationX(nouvelAnglebg)

		angleActueljd := p.jambeDroite.Rotation().X
		nouvelAnglejd := angleadouci(angleActueljd, 0, 8, delta)
		p.jambeDroite.SetRotationX(nouvelAnglejd)

		angleActueljg := p.jambeGauche.Rotation().X
		nouvelAnglejg := angleadouci(angleActueljg, 0, 8, delta)
		p.jambeGauche.SetRotationX(nouvelAnglejg)
		return
	}

	// La phase est l'horloge du balancement, pas l'orientation du personnage.
	// Le sinus alterne entre -1 et 1 ; 0.5 fixe l'amplitude en radians.
	p.phasemarche += delta * 5
	angle := float32(math.Sin(float64(p.phasemarche))) * 0.5
	// Les jambes sont opposées, comme les bras ; bras et jambe du même côté
	// partent en sens inverse pour donner une marche naturelle.
	p.jambeDroite.SetRotationX(angle)
	p.jambeGauche.SetRotationX(-angle)
	p.brasDroit.SetRotationX(-angle)
	p.brasGauche.SetRotationX(angle)
}

// angleadouci : Rapproche un angle de sa cible sans la dépasser, même si une image prend plus de temps.
func angleadouci(AngleActuel, AngleCible, Vitesse, Delta float32) float32 {
	facteur := Vitesse * Delta
	if facteur > 1 {
		facteur = 1
	}

	return AngleActuel + (AngleCible-AngleActuel)*facteur
}
