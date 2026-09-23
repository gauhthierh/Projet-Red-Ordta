package main

import "math"

// Bruit déterministe filtré : chaque génération produit les mêmes effets.
// Les sons de lame et d'impact ne sont plus des gammes de notes aiguës.
func (m *morceau) percussion(debut, duree, volume float64, typeSon string) {
	seed := uint32(747796405) + uint32(debut*rate)
	grave := 0.0
	for i := 0; i < int(duree*rate); i++ {
		t := float64(i) / rate
		seed ^= seed << 13
		seed ^= seed >> 17
		seed ^= seed << 5
		bruit := float64(seed)/2147483648 - 1
		grave += .08 * (bruit - grave)
		aigu := bruit - grave
		env := math.Min(1, t/.003) * math.Min(1, (duree-t)/.03)
		v := 0.0
		switch typeSon {
		case "lame":
			// Souffle gonflant puis décroissant, sans sifflement de synthétiseur.
			env *= math.Pow(math.Sin(math.Pi*t/duree), 2)
			v = grave*2.5 + aigu*.25
		case "metal":
			for h, f := range []float64{680, 1093, 1741, 2587, 3719} {
				v += math.Sin(2*math.Pi*f*t) * math.Exp(-t*(9+float64(h)*4)) / float64(h+2)
			}
			v += aigu * .3 * math.Exp(-90*t)
		case "tambour":
			v = (math.Sin(2*math.Pi*(58*t+1.8*(1-math.Exp(-35*t))))+.3*math.Sin(2*math.Pi*113*t))*math.Exp(-12*t) + grave*.8*math.Exp(-35*t)
		case "caisse":
			v = (aigu*.55 + math.Sin(2*math.Pi*170*t)*.35) * math.Exp(-28*t)
		case "cymbale":
			v = aigu * .35 * math.Exp(-4*t)
		}
		m.ajouter(int(debut*rate)+i, v*env*volume)
	}
}
