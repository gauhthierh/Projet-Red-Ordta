package personnages

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Chaque pièce a une position fixe. Les vêtements suivent les pivots animés.
func boite(parent *core.Node, x, y, z, largeur, profondeur, hauteur float32, couleur uint32) {
	m := material.NewStandard(&math32.Color{R: float32(couleur>>16&255) / 255, G: float32(couleur>>8&255) / 255, B: float32(couleur&255) / 255})
	mesh := graphic.NewMesh(geometry.NewBox(largeur, profondeur, hauteur), m)
	mesh.SetPosition(x, y, z)
	parent.Add(mesh)
}

// Nouveau : Crée le personnage joueur avec la tenue par défaut.
func Nouveau() *Personnage { return NouveauAvecTenue("joueur") }

// NouveauAvecTenue : Construit le corps et ses pivots, puis ajoute les détails propres au joueur ou au PNJ.
func NouveauAvecTenue(tenue string) *Personnage {
	p := &Personnage{noeud: core.NewNode()}
	chemise := uint32(0xd5c8a8)
	if tenue == "marchand" {
		chemise = 0x426543
	}
	if tenue == "forgeron" {
		chemise = 0x67534a
	}
	boite(p.noeud, 0, 0, 1.18, .63, .34, .65, chemise)
	boite(p.noeud, 0, 0, .87, .65, .37, .09, 0x50301b)
	boite(p.noeud, 0, -.2, .87, .12, .05, .10, 0xcdaa4b)
	boite(p.noeud, 0, 0, 1.55, .17, .19, .18, 0xc28e65)
	boite(p.noeud, 0, 0, 1.78, .43, .39, .43, 0xc28e65)
	boite(p.noeud, 0, .035, 1.99, .46, .4, .12, 0x493021)
	boite(p.noeud, 0, .19, 1.81, .45, .08, .32, 0x493021)
	// Le visage regarde vers -Y.
	boite(p.noeud, -.1, -.201, 1.82, .08, .025, .05, 0xf4eee0)
	boite(p.noeud, .1, -.201, 1.82, .08, .025, .05, 0xf4eee0)
	boite(p.noeud, -.1, -.219, 1.82, .03, .015, .04, 0x29221b)
	boite(p.noeud, .1, -.219, 1.82, .03, .015, .04, 0x29221b)
	boite(p.noeud, 0, -.23, 1.73, .07, .09, .1, 0xb97f57)
	boite(p.noeud, 0, -.202, 1.64, .12, .025, .025, 0x6b3d2c)
	p.jambeGauche = core.NewNode()
	p.jambeGauche.SetPosition(-.19, 0, .8)
	p.noeud.Add(p.jambeGauche)
	p.jambeDroite = core.NewNode()
	p.jambeDroite.SetPosition(.19, 0, .8)
	p.noeud.Add(p.jambeDroite)
	for _, jambe := range []*core.Node{p.jambeGauche, p.jambeDroite} {
		boite(jambe, 0, 0, -.33, .25, .28, .66, 0x4c4b43)
		boite(jambe, 0, -.045, -.71, .27, .39, .16, 0x614630)
		botte := creerBotteDetaillee(jambe)
		p.bottes = append(p.bottes, botte)
	}
	p.brasGauche = core.NewNode()
	p.brasGauche.SetPosition(-.43, 0, 1.43)
	p.noeud.Add(p.brasGauche)
	p.brasDroit = core.NewNode()
	p.brasDroit.SetPosition(.43, 0, 1.43)
	p.noeud.Add(p.brasDroit)
	for _, bras := range []*core.Node{p.brasGauche, p.brasDroit} {
		boite(bras, 0, 0, -.19, .22, .28, .40, chemise)
		boite(bras, 0, 0, -.45, .18, .22, .16, 0xc28e65)
		p.protectionsBras = append(p.protectionsBras, creerProtectionBras(bras))
	}
	p.chapeau = core.NewNode()
	p.noeud.Add(p.chapeau)
	if tenue == "marchand" {
		// Le PNJ garde son chapeau de marchand.
		boite(p.chapeau, 0, 0, 2.04, .62, .55, .07, 0x6c4827)
		boite(p.chapeau, 0, .025, 2.15, .43, .38, .19, 0x805c33)
		boite(p.chapeau, 0, -.19, 2.1, .44, .025, .045, 0xd5ad57)
	} else {
		creerChapeauDetaille(p.chapeau)
	}
	p.tunique = core.NewNode()
	p.noeud.Add(p.tunique)
	creerTuniqueDetaillee(p.tunique)
	p.Equiper(false, false, false)
	if tenue == "marchand" {
		p.chapeau.SetVisible(true)
		boite(p.noeud, 0, -.185, 1.22, .42, .055, .52, 0x305437)
		boite(p.noeud, .26, -.23, .89, .19, .15, .23, 0x976329)
	}
	if tenue == "forgeron" {
		boite(p.noeud, 0, -.195, 1.04, .52, .065, .88, 0x683d26)
		boite(p.noeud, 0, -.21, 1.57, .33, .09, .15, 0x493021)
		boite(p.brasDroit, 0, -.05, -.65, .07, .09, .48, 0x735137)
		boite(p.brasDroit, 0, -.05, -.83, .34, .16, .17, 0x777f82)
	}
	return p
}

// Equiper : Change uniquement la visibilité des pièces 3D ; les statistiques restent gérées dans library.
func (p *Personnage) Equiper(tete, torse, pieds bool) {
	p.chapeau.SetVisible(tete)
	p.tunique.SetVisible(torse)
	for _, protection := range p.protectionsBras {
		protection.SetVisible(torse)
	}
	for _, botte := range p.bottes {
		botte.SetVisible(pieds)
	}
}
