package personnages

import "github.com/g3n/engine/core"

// Le personnage et la scène utilisent tous deux Z comme axe vertical.
// Les membres sont des nœuds séparés pour pouvoir les animer.
type Personnage struct {
	noeud *core.Node

	jambeGauche *core.Node
	jambeDroite *core.Node
	brasGauche  *core.Node
	brasDroit   *core.Node

	phasemarche float32
	chapeau     *core.Node
	tunique     *core.Node
	bottes      []*core.Node
}

func NouveauAncienModele() *Personnage {
	modele := lireModele()
	p := &Personnage{noeud: core.NewNode()}

	ajouterPieces(p.noeud, modele.Corps)

	p.jambeGauche = ajouterMembre(p.noeud, -.24, .8, modele.JambeGauche)
	p.jambeDroite = ajouterMembre(p.noeud, .24, .8, modele.JambeDroite)
	p.brasGauche = ajouterMembre(p.noeud, -.24*1.65, 1.42, modele.BrasGauche)
	p.brasDroit = ajouterMembre(p.noeud, .24*1.65, 1.42, modele.BrasDroit)

	return p
}

// Le pivot se trouve à la hanche ou à l'épaule. Les pièces attachées
// au pivot le suivent quand son angle change.
func ajouterMembre(parent *core.Node, x, hauteur float32, pieces []piece) *core.Node {
	pivot := core.NewNode()
	pivot.SetPosition(x, 0, hauteur)
	parent.Add(pivot)
	ajouterPieces(pivot, pieces)
	return pivot
}

// Noeud renvoie la racine à ajouter à la scène et à déplacer sur le sol XY.
func (p *Personnage) Noeud() *core.Node {
	return p.noeud
}
