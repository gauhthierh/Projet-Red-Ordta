package personnages

import "math"

func (p *Personnage) Animer(delta float32, bouge bool) {
	if !bouge {
		p.jambeDroite.SetRotationX(0)
		p.jambeGauche.SetRotationX(0)
		return
	}

	p.phasemarche += delta * 5
	angle := float32(math.Sin(float64(p.phasemarche))) * 0.5
	p.jambeDroite.SetRotationX(angle)
	p.jambeGauche.SetRotationX(-angle)
}
