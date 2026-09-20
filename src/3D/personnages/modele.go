package personnages

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

// Personnage conserve la racine mobile et les articulations utiles aux animations.
// Chaque membre est un nœud : tourner ce nœud fait tourner toutes ses pièces.
type Personnage struct {
	noeud  *core.Node // Déplacement dans la scène (Z vertical).
	modele *core.Node // Dessin du personnage (Y vertical dans son repère local).

	jambeGauche *core.Node
	jambeDroite *core.Node
	brasGauche  *core.Node
	brasDroit   *core.Node

	phasemarche float32
}

// Nouveau assemble le personnage sans lancer d'animation.
func Nouveau() *Personnage {
	p := &Personnage{}
	p.noeud = core.NewNode()
	p.modele = core.NewNode()

	// L'ancien modèle a Y pour hauteur. Une rotation de 90° autour de X
	// place cette hauteur sur Z, sans changer les axes de la scène.
	p.modele.SetRotationX(math32.Pi / 2)
	p.noeud.Add(p.modele)

	p.ajouterCorps()
	p.jambeGauche = ajouterJambe(p.modele, -.24)
	p.jambeDroite = ajouterJambe(p.modele, .24)
	p.brasGauche = ajouterBras(p.modele, -.24*1.65)
	p.brasDroit = ajouterBras(p.modele, .24*1.65)
	p.ajouterTete()
	p.ajouterCheveuxCapeEtSac()
	p.ajouterEquipement()

	return p
}

// Noeud renvoie la racine à ajouter au monde et à déplacer sur le sol XY.
func (p *Personnage) Noeud() *core.Node {
	return p.noeud
}

func (p *Personnage) ajouterCorps() {
	// Les appels à boite indiquent : parent, position X/Y/Z, dimensions, matière.
	boite(p.modele, 0, 1.10, 0, .62, .73, .34, tissu)   // Torse.
	boite(p.modele, 0, 1.25, .19, .46, .40, .08, acier) // Plastron.
	cylindre(p.modele, 0, 1.55, 0, .12, .19, peau)      // Cou.
	boite(p.modele, 0, .83, 0, .68, .12, .40, bois)     // Ceinture.
	boite(p.modele, 0, .83, .23, .13, .13, .07, or)     // Boucle.
}

// Le pivot est à la hanche ; les boîtes sont placées sous ce pivot.
func ajouterJambe(parent *core.Node, positionX float32) *core.Node {
	jambe := core.NewNode()
	jambe.SetPosition(positionX, .8, 0)
	parent.Add(jambe)

	boite(jambe, 0, -.24, 0, .22, .48, .24, tissu)  // Pantalon.
	boite(jambe, 0, -.60, .07, .26, .28, .40, bois) // Botte.
	boite(jambe, 0, -.49, .20, .12, .07, .04, or)   // Décoration.
	return jambe
}

// Le pivot est à l'épaule ; tourner ce nœud déplacera tout le bras.
func ajouterBras(parent *core.Node, positionX float32) *core.Node {
	bras := core.NewNode()
	bras.SetPosition(positionX, 1.42, 0)
	parent.Add(bras)

	sphere(bras, 0, 0, 0, .22, acier)             // Épaulette.
	boite(bras, 0, -.24, 0, .19, .42, .21, tissu) // Manche.
	boite(bras, 0, -.46, 0, .19, .15, .23, bois)  // Poignet.
	sphere(bras, 0, -.58, .04, .115, peau)        // Main.
	return bras
}

func (p *Personnage) ajouterTete() {
	tete := sphere(p.modele, 0, 1.82, 0, .29, peau)
	tete.SetScale(1, 1.12, .9)

	ajouterOeil(p.modele, -.105)
	ajouterOeil(p.modele, .105)
	boite(p.modele, 0, 1.80, .265, .075, .10, .07, peau) // Nez.
	boite(p.modele, 0, 1.70, .24, .10, .025, .025, bois) // Bouche.
}

func ajouterOeil(parent *core.Node, positionX float32) {
	boite(parent, positionX, 1.86, .243, .08, .065, .028, couleur(.94, .89, .70)) // Blanc.
	boite(parent, positionX, 1.86, .265, .035, .045, .018, sombre)                // Pupille.
	boite(parent, positionX, 1.925, .24, .12, .035, .03, bois)                    // Sourcil.
	sphere(parent, positionX*2.65, 1.82, 0, .07, peau)                            // Oreille.
}

func (p *Personnage) ajouterCheveuxCapeEtSac() {
	cheveux := sphere(p.modele, 0, 2.02, -.02, .28, bois)
	cheveux.SetScale(1, .55, 1)
	boite(p.modele, -.22, 1.91, -.03, .09, .27, .30, bois)
	boite(p.modele, .22, 1.91, -.03, .09, .27, .30, bois)

	boite(p.modele, 0, 1.12, -.24, .58, .82, .08, couleur(.48, .12, .10)) // Cape.
	boite(p.modele, .13, 1.36, -.31, .36, .42, .16, bois)                 // Sac.
}

func (p *Personnage) ajouterEquipement() {
	// L'épée est attachée au bras droit : elle le suivra quand il tournera.
	boite(p.brasDroit, 0, -.3, .18, .075, .72, .045, acier) // Lame.
	boite(p.brasDroit, 0, -.64, .18, .31, .055, .08, or)    // Garde.
	boite(p.brasDroit, 0, -.74, .18, .075, .18, .07, bois)  // Poignée.

	// Même principe pour le bouclier, attaché au bras gauche.
	bouclier := sphere(p.brasGauche, -.1, -.35, .16, .29, bois)
	bouclier.SetScale(1, 1.3, .22)
	sphere(p.brasGauche, -.1, -.35, .23, .075, acier)
}
